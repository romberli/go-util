package linux

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pingcap/errors"
	"github.com/pkg/sftp"
	"github.com/romberli/go-util/constant"
	"golang.org/x/crypto/ssh"

	utilrand "k8s.io/apimachinery/pkg/util/rand"
)

const (
	defaultRandStringLength = 6

	DefaultSplitStr           = constant.CRLFString
	DefaultSuccessReturnValue = constant.ZeroInt
	DefaultFailedReturnValue  = constant.OneInt
	DefaultSSHTimeout         = 10 * time.Second
	DefaultSSHPortNum         = 22
	DefaultSSHUserName        = "root"
	DefaultSSHUserPass        = "root"
	DefaultByteBufferSize     = 1024 * 1024 // 1MB

	hostNameCommand   = "/usr/bin/hostname"
	pathExistsCommand = "/usr/bin/test -e %s && echo 0 || echo 1"
	isDirCommand      = "/usr/bin/test -d %s && echo 0 || echo 1"
	lsCommand         = "/usr/bin/ls %s"
	mkdirCommand      = "/usr/bin/mkdir -p %s"
	rmCommand         = "/usr/bin/rm -rf %s"
	catCommand        = "/usr/bin/cat %s"
	touchCommand      = "/usr/bin/touch %s"
	cpCommand         = "/usr/bin/cp -r %s %s"
	mvCommand         = "/usr/bin/mv %s %s"
	chownCommand      = "/usr/bin/chown -R %s:%s %s"
	chmodCommand      = "/usr/bin/chmod -R %s %s"

	sudoPrefix = "sudo "
)

type SSHConfig struct {
	HostIp   string
	PortNum  int
	UserName string
	UserPass string
	useSudo  bool
}

func NewSSHConfig(hostIP string, portNum int, userName string, userPass string, useSudo bool) *SSHConfig {
	return &SSHConfig{
		hostIP,
		portNum,
		userName,
		userPass,
		useSudo,
	}
}

func NewSSHConfigWithDefault(hostIP string) *SSHConfig {
	return &SSHConfig{
		hostIP,
		DefaultSSHPortNum,
		DefaultSSHUserName,
		DefaultSSHUserPass,
		false,
	}
}

type SSHConn struct {
	Config     *SSHConfig
	SSHClient  *ssh.Client
	SFTPClient *sftp.Client
}

// NewSSHConn returns a new *SSHConn
func NewSSHConn(hostIP string, portNum int, userName, userPass string, useSudo bool) (*SSHConn, error) {
	return newSSHConnWithConfig(NewSSHConfig(hostIP, portNum, userName, userPass, useSudo))
}

// newSSHConnWithConfig returns *SSHConn with given config
func newSSHConnWithConfig(config *SSHConfig) (*SSHConn, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
	)

	// get auth method
	auth = append(auth, ssh.Password(config.UserPass))

	hostKeyCallBack := func(host string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	clientConfig = &ssh.ClientConfig{
		User:            config.UserName,
		Auth:            auth,
		Timeout:         DefaultSSHTimeout,
		HostKeyCallback: hostKeyCallBack,
	}

	// connect to ssh
	addr = fmt.Sprintf("%s:%d", config.HostIp, config.PortNum)
	sshClient, err := ssh.Dial(constant.TransportProtocolTCP, addr, clientConfig)
	if err != nil {
		return nil, errors.Trace(err)
	}

	// create sftp client
	sftpClient, err = sftp.NewClient(sshClient)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &SSHConn{
		config,
		sshClient,
		sftpClient,
	}, nil
}

// Close closes connections with the remote host
func (conn *SSHConn) Close() error {
	err := conn.SFTPClient.Close()
	if err != nil {
		return errors.Trace(err)
	}

	return errors.Trace(conn.SSHClient.Close())
}

// SetUseSudo sets if use sudo or not
func (conn *SSHConn) SetUseSudo(useSudo bool) {
	conn.Config.useSudo = useSudo
}

// ExecuteCommand executes shell command on the remote host
func (conn *SSHConn) ExecuteCommand(cmd string) (string, error) {
	return conn.executeCommand(cmd)
}

// ExecuteCommandWithoutOutput executes a command without output
func (conn *SSHConn) ExecuteCommandWithoutOutput(cmd string) error {
	_, err := conn.executeCommand(cmd)

	return err
}

func (conn *SSHConn) executeCommand(cmd string) (string, error) {
	var (
		stdOutBuffer bytes.Buffer
		stdErrBuffer bytes.Buffer
	)

	// create ssh session
	sshSession, err := conn.SSHClient.NewSession()
	if err != nil {
		return constant.EmptyString, errors.Trace(err)
	}
	defer func() { _ = sshSession.Close() }()

	sshSession.Stdout = &stdOutBuffer
	sshSession.Stderr = &stdErrBuffer

	// prepare command
	if conn.Config.useSudo && !strings.HasPrefix(cmd, sudoPrefix) {
		cmd = sudoPrefix + cmd
	}
	// run command
	err = sshSession.Run(cmd)
	if err != nil {
		if stdErrBuffer.String() != constant.EmptyString {
			err = errors.Errorf("%s%+v", stdErrBuffer.String(), errors.Trace(err))
		}
	}

	output := stdOutBuffer.String() + stdErrBuffer.String()
	return strings.TrimSpace(output), errors.Trace(err)
}

// GetHostName returns hostname of remote host
func (conn *SSHConn) GetHostName() (string, error) {
	return conn.ExecuteCommand(hostNameCommand)
}

// PathExists returns if given path exists
func (conn *SSHConn) PathExists(path string) (bool, error) {
	cmd := fmt.Sprintf(pathExistsCommand, strings.TrimSpace(path))
	output, err := conn.ExecuteCommand(cmd)
	if err != nil {
		return false, err
	}

	return output == strconv.Itoa(constant.ZeroInt), nil
}

// IsDir returns if given path on the remote host is a directory or not
func (conn *SSHConn) IsDir(path string) (bool, error) {
	cmd := fmt.Sprintf(isDirCommand, strings.TrimSpace(path))
	output, err := conn.ExecuteCommand(cmd)
	if err != nil {
		return false, err
	}

	return output == strconv.Itoa(constant.ZeroInt), nil
}

// ListPath returns subdirectories and files of given path on the remote host, it returns a slice of sub paths
func (conn *SSHConn) ListPath(path string) ([]string, error) {
	var subPathList []string

	cmd := fmt.Sprintf(lsCommand, strings.TrimSpace(path))
	subPathStr, err := conn.ExecuteCommand(cmd)
	if err != nil {
		return nil, err
	}

	subPathStr = strings.TrimSpace(subPathStr)
	if subPathStr != constant.EmptyString {
		subPathList = strings.Split(subPathStr, DefaultSplitStr)
	}

	return subPathList, nil
}

// ReadDir returns subdirectories and files of given directory on the remote host, it returns a slice of os.FileInfo
func (conn *SSHConn) ReadDir(dirName string) ([]os.FileInfo, error) {
	var fileInfoList []os.FileInfo

	dirName = strings.TrimSpace(dirName)
	isDir, err := conn.IsDir(dirName)
	if err != nil {
		return nil, err
	}
	if !isDir {
		return nil, errors.Errorf("it's not a directory. dir name: %s", dirName)
	}

	subPathList, err := conn.ListPath(dirName)
	if err != nil {
		return nil, err
	}
	for _, subPath := range subPathList {
		if subPath != constant.EmptyString {
			fileNameAbs := filepath.Join(dirName, subPath)
			fileInfo, err := conn.SFTPClient.Stat(fileNameAbs)
			if err != nil {
				return nil, errors.Trace(err)
			}

			fileInfoList = append(fileInfoList, fileInfo)
		}
	}

	return fileInfoList, nil
}

// GetFileInfo

// MkdirAll creates a directory named path, along with any necessary parents, on the remote host, it will act like shell command "mkdir -p $path"
func (conn *SSHConn) MkdirAll(path string) error {
	return conn.ExecuteCommandWithoutOutput(fmt.Sprintf(mkdirCommand, strings.TrimSpace(path)))
}

// RemoveAll removes given path on the remote host, it will act like shell command "rm -rf $path",
func (conn *SSHConn) RemoveAll(path string) error {
	return conn.ExecuteCommandWithoutOutput(fmt.Sprintf(rmCommand, strings.TrimSpace(path)))
}

// Cat returns the content of the given file on the remote host, it will act like shell command "cat $path"
func (conn *SSHConn) Cat(path string) (string, error) {
	return conn.ExecuteCommand(fmt.Sprintf(catCommand, strings.TrimSpace(path)))
}

// Touch touches the given path on the remote host, it will act like shell command "touch $path"
func (conn *SSHConn) Touch(path string) error {
	return conn.ExecuteCommandWithoutOutput(fmt.Sprintf(touchCommand, strings.TrimSpace(path)))
}

// Copy copies a file or directory on the remote host, it will act like shell command "copy -r $src $dest"
func (conn *SSHConn) Copy(src, dest string) error {
	return conn.ExecuteCommandWithoutOutput(fmt.Sprintf(cpCommand, strings.TrimSpace(src), strings.TrimSpace(dest)))
}

// Move moves a file or directory on the remote host, it will act like shell command "mv $src $dest"
func (conn *SSHConn) Move(src, dest string) error {
	return conn.ExecuteCommandWithoutOutput(fmt.Sprintf(mvCommand, strings.TrimSpace(src), strings.TrimSpace(dest)))
}

// Chown changes the owner and group of the given path on the remote host, it will act like shell command "chown -R $user:$group $path"
func (conn *SSHConn) Chown(path, user, group string) error {
	return conn.ExecuteCommandWithoutOutput(fmt.Sprintf(chownCommand, user, group, strings.TrimSpace(path)))
}

// Chmod changes the mode of the given path on the remote host, it will act like shell command "chmod -R $mode $path"
func (conn *SSHConn) Chmod(path string, mode string) error {
	return conn.ExecuteCommandWithoutOutput(fmt.Sprintf(chmodCommand, mode, strings.TrimSpace(path)))
}

// IsEmptyDir returns if  given directory is empty or not on the remote host
func (conn *SSHConn) IsEmptyDir(dirName string) (bool, error) {
	subPathList, err := conn.ListPath(strings.TrimSpace(dirName))
	if err != nil {
		return false, err
	}

	return len(subPathList) == constant.ZeroInt, nil
}

// CopyFile copy file content from source to destination, it doesn't care about which one is local or remote
func (conn *SSHConn) CopyFile(fileSource io.Reader, fileDest io.Writer, bufferSize int) error {
	if bufferSize <= constant.ZeroInt {
		bufferSize = DefaultByteBufferSize
	}

	buf := make([]byte, bufferSize)

	for {
		n, err := fileSource.Read(buf)
		if err != nil && err != io.EOF {
			return errors.Trace(err)
		}

		if n == constant.ZeroInt {
			break
		}

		_, err = fileDest.Write(buf[constant.ZeroInt:n])
		if err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}

// CopySingleFileFromRemote copies one single file from remote to local.
// if tmpDir is not empty and is different with the parent directory of fileNameSource,
// it will copy the file to tmpDir on the remote host first, then transfer it to the local host,
// after that, it will remove the temporary file on the remote host automatically.
// it is your responsibility to assure that tmpDir exists and is a directory,
// the connection should have enough privilege to the tmpDir
// and the tmpDir has enough space to hold the file temporarily,
// note that only the first tmpDir will be used.
func (conn *SSHConn) CopySingleFileFromRemote(fileNameSource string, fileNameDest string, tmpDir ...string) error {
	var (
		fileDest   *os.File
		fileSource *sftp.File
	)

	fileNameSource = strings.TrimSpace(fileNameSource)
	fileNameDest = strings.TrimSpace(fileNameDest)

	if fileNameDest == constant.EmptyString {
		fileNameDest = fileNameSource
	}

	// check if source path is a directory
	isDir, err := conn.IsDir(fileNameSource)
	if err != nil {
		return err
	}
	if isDir {
		return errors.Errorf("it's not a file. file name: %s", fileNameSource)
	}

	// check if parent path of destination exists
	fileNameDestParent := filepath.Dir(fileNameDest)
	pathExists, err := PathExists(fileNameDestParent)
	if err != nil {
		return err
	}
	if !pathExists {
		return errors.Errorf("parent path of destination does not exist nor have privilege. path: %s", fileNameDest)
	}

	fileNameSourceParent := filepath.Dir(fileNameSource)
	fileNameSourceBase := filepath.Base(fileNameSource)
	// check if destination path is a directory
	pathExists, err = PathExists(fileNameDest)
	if err != nil {
		return err
	}
	if pathExists {
		isDir, err = IsDir(fileNameDest)
		if err != nil {
			return err
		}
		if isDir {
			fileNameDest = filepath.Join(fileNameDest, fileNameSourceBase)
		}
	}

	var td string
	if len(tmpDir) > constant.ZeroInt {
		// only use the first one
		td = tmpDir[constant.ZeroInt]
	}

	tmpFile := fileNameSource + constant.DotString + utilrand.String(defaultRandStringLength)
	if td != constant.EmptyString && td != fileNameSourceParent {
		isDir, err = conn.IsDir(td)
		if err != nil {
			return err
		}
		if !isDir {
			return errors.Errorf("tmpDir is not a directory. path: %s", td)
		}
		tmpFile = filepath.Join(td, fileNameSourceBase) + constant.DotString + utilrand.String(defaultRandStringLength)
		err = conn.Copy(fileNameSource, tmpFile)
		if err != nil {
			return err
		}
		// remove tmp file
		defer func() { _ = conn.RemoveAll(tmpFile) }()
		// chmod
		err = conn.Chmod(tmpFile, constant.DefaultAllFileModeStr)
		if err != nil {
			return err
		}
	}

	fileSource, err = conn.SFTPClient.Open(tmpFile)
	if err != nil {
		return errors.Trace(err)
	}
	defer func() { _ = fileSource.Close() }()

	fileDest, err = os.Create(fileNameDest)
	if err != nil {
		return errors.Trace(err)
	}
	defer func() { _ = fileDest.Close() }()

	err = conn.CopyFile(fileSource, fileDest, DefaultByteBufferSize)
	if err != nil {
		return err
	}

	return nil
}

// CopySingleFileToRemote copies one single file from local to remote
// if tmpDir is not empty and is different with the parent directory of fileNameDest,
// it will transfer the file to tmpDir on the remote host first, then copy it to fileNameDest,
// after that, it will remove the temporary file on the remote host automatically.
// it is your responsibility to assure that tmpDir exists and is a directory,
// the connection should have enough privilege to the tmpDir
// and the tmpDir has enough space to hold the file temporarily,
// note that only the first tmpDir will be used.
func (conn *SSHConn) CopySingleFileToRemote(fileNameSource string, fileNameDest string, tmpDir ...string) error {
	var (
		fileSource *os.File
		fileDest   *sftp.File
	)

	fileNameSource = strings.TrimSpace(fileNameSource)
	fileNameDest = strings.TrimSpace(fileNameDest)

	if fileNameDest == constant.EmptyString {
		fileNameDest = fileNameSource
	}

	// check if source path is a directory
	isDir, err := IsDir(fileNameSource)
	if err != nil {
		return err
	}
	if isDir {
		return errors.Errorf("it's not a file. file name: %s", fileNameSource)
	}

	// check if parent path of destination exists
	fileNameDestParent := filepath.Dir(fileNameDest)
	fileNameDestBase := filepath.Base(fileNameDest)
	pathExists, err := conn.PathExists(fileNameDestParent)
	if err != nil {
		return nil
	}
	if !pathExists {
		return errors.Errorf("parent path of destination does not exist. path: %s", fileNameDest)
	}

	// check if destination path is a directory
	pathExists, err = conn.PathExists(fileNameDest)
	if err != nil {
		return err
	}
	if pathExists {
		isDir, err = conn.IsDir(fileNameDest)
		if err != nil {
			return err
		}
		if isDir {
			fileNameSourceBase := filepath.Base(fileNameSource)
			fileNameDest = filepath.Join(fileNameDest, fileNameSourceBase)
		}
	}

	fileSource, err = os.Open(fileNameSource)
	if err != nil {
		return errors.Trace(err)
	}
	defer func() { _ = fileSource.Close() }()

	var (
		td          string
		usedTmpFlag bool
	)
	if len(tmpDir) > constant.ZeroInt {
		// only use the first one
		td = tmpDir[constant.ZeroInt]
	}

	tmpFile := fileNameDest + constant.DotString + utilrand.String(defaultRandStringLength)
	if td != constant.EmptyString && td != fileNameDestParent {
		isDir, err = conn.IsDir(td)
		if err != nil {
			return err
		}
		if !isDir {
			return errors.Errorf("tmpDir is not a directory. path: %s", td)
		}
		tmpFile = filepath.Join(td, fileNameDestBase) + constant.DotString + utilrand.String(defaultRandStringLength)
		usedTmpFlag = true
	}

	fileDest, err = conn.SFTPClient.Create(tmpFile)
	if err != nil {
		return errors.Trace(err)
	}
	defer func() { _ = fileDest.Close() }()
	// transfer data to the temporary file
	err = conn.CopyFile(fileSource, fileDest, DefaultByteBufferSize)
	if err != nil {
		return err
	}

	if usedTmpFlag {
		// copy temporary file to the file dest
		err = conn.Copy(tmpFile, fileNameDest)
		if err != nil {
			return err
		}
		// remove tmp file
		return conn.RemoveAll(tmpFile)
	}

	return nil
}

// GetPathDirMapRemote reads all subdirectories and files of given directory on the remote host
// and calculate the relative path of rootPath,
// then map the absolute path of subdirectory names and file names as keys, relative paths as values to fileDirMap
func (conn *SSHConn) GetPathDirMapRemote(dirName, rootPath string) (map[string]string, error) {
	fileDirMap := make(map[string]string)

	err := conn.getPathDirMapRemote(fileDirMap, dirName, rootPath)
	if err != nil {
		return nil, err
	}

	return fileDirMap, nil
}

func (conn *SSHConn) getPathDirMapRemote(fileDirMap map[string]string, dirName, rootPath string) error {
	dirName = strings.TrimSpace(dirName)
	rootPath = strings.TrimSpace(rootPath)

	subPaths, err := conn.ListPath(dirName)
	if err != nil {
		return err
	}

	if subPaths == nil {
		// it's an empty directory
		fileDirMap[dirName] = constant.EmptyString
	}

	for _, subPath := range subPaths {
		subPathAbs := filepath.Join(dirName, subPath)

		isDir, err := conn.IsDir(subPathAbs)
		if err != nil {
			return err
		}

		if isDir {
			// call recursively
			err = conn.getPathDirMapRemote(fileDirMap, subPathAbs, rootPath)
			if err != nil {
				return err
			}
		} else {
			// get relative path with root path
			fileNameRel, err := filepath.Rel(rootPath, subPathAbs)
			if err != nil {
				return errors.Trace(err)
			}

			fileDirMap[subPathAbs] = fileNameRel
		}
	}

	return nil
}

// CopyDirFromRemote copies a directory with all subdirectories and files from remote to local
func (conn *SSHConn) CopyDirFromRemote(dirNameSource, dirNameDest string, tmpDir ...string) error {
	dirNameSource = strings.TrimSpace(dirNameSource)
	dirNameDest = strings.TrimSpace(dirNameDest)

	// check if source path is a directory
	isDir, err := conn.IsDir(dirNameSource)
	if err != nil {
		return err
	}
	if !isDir {
		return errors.Errorf("it's not a directory. dir name: %s", dirNameSource)
	}

	// check if parent path of destination exists
	dirNameDestParent := filepath.Dir(dirNameDest)
	pathExists, err := PathExists(dirNameDestParent)
	if err != nil {
		return nil
	}
	if !pathExists {
		return errors.Errorf("parent path of destination does not exists. path: %s", dirNameDest)
	}

	pathSourceBase := filepath.Base(dirNameSource)
	pathDestBase := filepath.Base(dirNameDest)

	// get new destination path to act like shell command "scp -r"
	if pathSourceBase != pathDestBase {
		pathExists, err = PathExists(dirNameDest)
		if err != nil {
			return err
		}
		if pathExists {
			isDir, err = IsDir(dirNameDest)
			if err != nil {
				return err
			}
			if !isDir {
				return errors.Errorf("it's not a directory. path: %s", dirNameDest)
			}

			dirNameDest = filepath.Join(dirNameDest, pathSourceBase)
		}
	}

	// get map of source path and relative path with destination directory
	pathDirMap, err := conn.GetPathDirMapRemote(dirNameSource, dirNameSource)
	if err != nil {
		return err
	}

	for pathName, relDir := range pathDirMap {
		if relDir == constant.EmptyString {
			// it's an empty directory, we just need to create it
			relDirSource, err := filepath.Rel(dirNameSource, pathName)
			if err != nil {
				return errors.Trace(err)
			}

			dirDestAbs := filepath.Join(dirNameDest, relDirSource)
			err = os.MkdirAll(dirDestAbs, constant.DefaultExecFileMode)
			if err != nil {
				return errors.Trace(err)
			}
			continue
		}

		relDir = filepath.Dir(relDir)
		DirDestAbs := filepath.Join(dirNameDest, relDir)
		err = os.MkdirAll(DirDestAbs, constant.DefaultExecFileMode)
		if err != nil {
			return errors.Trace(err)
		}

		fileNameDest := GetFileNameDest(pathName, DirDestAbs)
		// copy file from remote
		err = conn.CopySingleFileFromRemote(pathName, fileNameDest, tmpDir...)
		if err != nil {
			return err
		}
	}

	return nil
}

// CopyDirToRemote copies a directory with all subdirectories and files from local to remote
func (conn *SSHConn) CopyDirToRemote(dirNameSource, dirNameDest string, tmpDir ...string) error {
	dirNameSource = strings.TrimSpace(dirNameSource)
	dirNameDest = strings.TrimSpace(dirNameDest)

	// check if source path is a directory
	isDir, err := IsDir(dirNameSource)
	if err != nil {
		return err
	}
	if !isDir {
		return errors.Errorf("it's not a directory. dir name: %s", dirNameSource)
	}

	// check if parent path of destination exists
	dirNameDestParent := filepath.Dir(dirNameDest)
	pathExists, err := conn.PathExists(dirNameDestParent)
	if err != nil {
		return err
	}
	if !pathExists {
		return errors.Errorf("parent path of destination does not exsists. path: %s", dirNameDest)
	}

	pathSourceBase := filepath.Base(dirNameSource)
	pathDestBase := filepath.Base(dirNameDest)

	// get new destination path to act like shell command "scp -r"
	if pathSourceBase != pathDestBase {
		pathExists, err = conn.PathExists(dirNameDest)
		if err != nil {
			return err
		}
		if pathExists {
			isDir, err = conn.IsDir(dirNameDest)
			if err != nil {
				return err
			}
			if !isDir {
				return errors.Errorf("it's not a directory. dir name: %s", dirNameDest)
			}

			dirNameDest = filepath.Join(dirNameDest, pathSourceBase)
		}
	}

	// get map of source path and relative path with destination directory
	pathDirMap, err := GetPathDirMapLocal(dirNameSource, dirNameSource)
	if err != nil {
		return err
	}

	for pathName, relDir := range pathDirMap {
		if relDir == constant.EmptyString {
			//
			relDirSource, err := filepath.Rel(dirNameSource, pathName)
			if err != nil {
				return errors.Trace(err)
			}

			dirDestAbs := filepath.Join(dirNameDest, relDirSource)
			err = conn.MkdirAll(dirDestAbs)
			if err != nil {
				return errors.Trace(err)
			}

			continue
		}

		relDir = filepath.Dir(relDir)
		DirDestAbs := filepath.Join(dirNameDest, relDir)
		err = conn.MkdirAll(DirDestAbs)
		if err != nil {
			return errors.Trace(err)
		}

		fileNameDest := GetFileNameDest(pathName, DirDestAbs)
		err = conn.CopySingleFileToRemote(pathName, fileNameDest, tmpDir...)
		if err != nil {
			return err
		}
	}

	return nil
}

// CopyFromRemote copies no matter a directory or a file from remote to local
func (conn *SSHConn) CopyFromRemote(pathSource, pathDest string, tmpDir ...string) (err error) {
	pathSource = strings.TrimSpace(pathSource)
	pathDest = strings.TrimSpace(pathDest)

	// check if source path is a directory
	isDir, err := conn.IsDir(pathSource)
	if err != nil {
		return err
	}
	if isDir {
		return conn.CopyDirFromRemote(pathSource, pathDest, tmpDir...)
	}

	return conn.CopySingleFileFromRemote(pathSource, pathDest, tmpDir...)
}

// CopyToRemote copies no matter a directory or a file from local to remote
func (conn *SSHConn) CopyToRemote(pathSource, pathDest string, tmpDir ...string) (err error) {
	pathSource = strings.TrimSpace(pathSource)
	pathDest = strings.TrimSpace(pathDest)

	isDir, err := IsDir(pathSource)
	if err != nil {
		return err
	}
	if isDir {
		return conn.CopyDirToRemote(pathSource, pathDest, tmpDir...)
	}

	return conn.CopySingleFileToRemote(pathSource, pathDest, tmpDir...)
}
