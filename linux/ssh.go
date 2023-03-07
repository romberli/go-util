package linux

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/pingcap/errors"
	"github.com/pkg/sftp"
	"github.com/romberli/go-util/constant"
	"golang.org/x/crypto/ssh"
)

const (
	DefaultSplitStr           = "\n"
	DefaultSuccessReturnValue = 0
	DefaultFailedReturnValue  = 1
	DefaultSSHTimeout         = 10 * time.Second
	DefaultSSHPortNum         = 22
	DefaultSSHUserName        = "root"
	DefaultSSHUserPass        = "root"
	DefaultByteBufferSize     = 1024 * 1024 // 1MB

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
	Config    *SSHConfig
	SSHClient *ssh.Client
	*sftp.Client
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
	err := conn.Client.Close()
	if err != nil {
		return errors.Trace(err)
	}

	return errors.Trace(conn.SSHClient.Close())
}

// ExecuteCommand executes shell command on the remote host
func (conn *SSHConn) ExecuteCommand(cmd string) (string, error) {
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
	if conn.Config.useSudo {
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
	return output, errors.Trace(err)
}

// GetHostName returns hostname of remote host
func (conn *SSHConn) GetHostName() (string, error) {
	hostName, err := conn.ExecuteCommand(HostNameCommand)

	return hostName, err
}

// PathExists returns if given path exists
func (conn *SSHConn) PathExists(path string) (bool, error) {
	_, err := conn.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, errors.Trace(err)
}

// IsDir returns if given path on the remote host is a directory or not
func (conn *SSHConn) IsDir(path string) (bool, error) {
	path = strings.TrimSpace(path)

	info, err := conn.Stat(path)
	if err != nil {
		return false, errors.Trace(err)
	}

	return info.IsDir(), nil
}

// ListPath returns subdirectories and files of given path on the remote host, it returns a slice of sub paths
func (conn *SSHConn) ListPath(path string) ([]string, error) {
	var subPathList []string

	cmd := fmt.Sprintf("%s %s", LsCommand, strings.TrimSpace(path))
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
			fileInfo, err := conn.Stat(fileNameAbs)
			if err != nil {
				return nil, errors.Trace(err)
			}

			fileInfoList = append(fileInfoList, fileInfo)
		}
	}

	return fileInfoList, nil
}

// RemoveAll removes given path on the remote host, it will act like shell command "rm -rf $path",
// except that it will raise an error when something goes wrong.
func (conn *SSHConn) RemoveAll(path string) error {
	path = strings.TrimSpace(path)
	isDir, err := conn.IsDir(path)
	if err != nil {
		return err
	}

	if isDir {
		isEmpty, err := conn.IsEmptyDir(path)
		if err != nil {
			return err
		}

		if !isEmpty {
			subPathList, err := conn.ListPath(path)
			if err != nil {
				return err
			}
			for _, subPath := range subPathList {
				subPathAbs := filepath.Join(path, subPath)
				err = conn.RemoveAll(subPathAbs)
				if err != nil {
					return err
				}
			}
		}

		err = conn.RemoveDirectory(path)
		if err != nil {
			return errors.Trace(err)
		}
	} else {
		err = conn.Remove(path)
		if err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}

// IsEmptyDir returns if  given directory is empty or not
func (conn *SSHConn) IsEmptyDir(dirName string) (bool, error) {
	dirName = strings.TrimSpace(dirName)
	fileInfoList, err := conn.ReadDir(dirName)
	if err != nil {
		return false, err
	}

	return fileInfoList == nil, nil
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

// CopySingleFileFromRemote copies one single file from remote to local
func (conn *SSHConn) CopySingleFileFromRemote(fileNameSource string, fileNameDest string) error {
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
		return errors.Errorf("parent path of destination does not exist. path: %s", fileNameDest)
	}

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
			fileNameSourceBase := filepath.Base(fileNameSource)
			fileNameDest = filepath.Join(fileNameDest, fileNameSourceBase)
		}
	}

	fileSource, err = conn.Open(fileNameSource)
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
func (conn *SSHConn) CopySingleFileToRemote(fileNameSource string, fileNameDest string) error {
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

	fileDest, err = conn.Create(fileNameDest)
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

// CopyFileListFromRemote copies given files from remote to local
func (conn *SSHConn) CopyFileListFromRemote(fileListSource []string, FileDirDest string) error {
	FileDirDest = strings.TrimSpace(FileDirDest)
	if FileDirDest == constant.EmptyString {
		return errors.New("file destination directory should not an empty string")
	}

	pathExists, err := PathExists(FileDirDest)
	if err != nil {
		return err
	}

	if !pathExists {
		_, err = os.Create(FileDirDest)
		if err != nil {
			return errors.Trace(err)
		}
	}

	for _, fileNameSource := range fileListSource {
		fileNameSource = strings.TrimSpace(fileNameSource)
		fileNameDest := path.Base(fileNameSource)

		err = conn.CopySingleFileFromRemote(fileNameSource, path.Join(FileDirDest, fileNameDest))
		if err != nil {
			return err
		}
	}

	return nil
}

// CopyFileListFromRemoteWithNewName copies file from remote to local,
// it copies file contents and rename files to given file names
func (conn *SSHConn) CopyFileListFromRemoteWithNewName(fileListSource []string, fileListDest []string) (err error) {
	if len(fileListSource) != len(fileListDest) {
		return errors.Errorf("the length of source and destination file list must be exactly same. source length: %d, destination length: %d",
			len(fileListSource), len(fileListDest))
	}

	for i, fileNameSource := range fileListSource {
		fileNameDest := fileListDest[i]
		fileNameDest = strings.TrimSpace(fileNameDest)

		if fileNameDest == constant.EmptyString {
			return errors.New("destination file name should not be empty")
		}

		err = conn.CopySingleFileFromRemote(strings.TrimSpace(fileNameSource), fileNameDest)
		if err != nil {
			return err
		}
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

	fileInfoList, err := conn.ReadDir(dirName)
	if err != nil {
		return err
	}

	if fileInfoList == nil {
		// it's an empty directory
		fileDirMap[dirName] = constant.EmptyString
	}

	for _, fileInfo := range fileInfoList {
		fileName := fileInfo.Name()
		fileNameAbs := filepath.Join(dirName, fileName)

		if fileInfo.IsDir() {
			// call recursively
			err = conn.getPathDirMapRemote(fileDirMap, fileNameAbs, rootPath)
			if err != nil {
				return err
			}
		} else {
			// get relative path with root path
			fileNameRel, err := filepath.Rel(rootPath, fileNameAbs)
			if err != nil {
				return errors.Trace(err)
			}

			fileDirMap[fileNameAbs] = fileNameRel
		}
	}

	return nil
}

// CopyDirFromRemote copies a directory with all subdirectories and files from remote to local
func (conn *SSHConn) CopyDirFromRemote(dirNameSource, dirNameDest string) error {
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
		err = conn.CopySingleFileFromRemote(pathName, fileNameDest)
		if err != nil {
			return err
		}
	}

	return nil
}

// CopyDirToRemote copies a directory with all subdirectories and files from local to remote
func (conn *SSHConn) CopyDirToRemote(dirNameSource, dirNameDest string) error {
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
		err = conn.CopySingleFileToRemote(pathName, fileNameDest)
		if err != nil {
			return err
		}
	}

	return nil
}

// CopyFromRemote copies no matter a directory or a file from remote to local
func (conn *SSHConn) CopyFromRemote(pathSource, pathDest string) (err error) {
	pathSource = strings.TrimSpace(pathSource)
	pathDest = strings.TrimSpace(pathDest)

	// check if source path is a directory
	isDir, err := conn.IsDir(pathSource)
	if err != nil {
		return err
	}
	if isDir {
		return conn.CopyDirFromRemote(pathSource, pathDest)
	}

	return conn.CopySingleFileFromRemote(pathSource, pathDest)
}

// CopyToRemote copies no matter a directory or a file from local to remote
func (conn *SSHConn) CopyToRemote(pathSource, pathDest string) (err error) {
	pathSource = strings.TrimSpace(pathSource)
	pathDest = strings.TrimSpace(pathDest)

	isDir, err := IsDir(pathSource)
	if err != nil {
		return err
	}
	if isDir {
		return conn.CopyDirToRemote(pathSource, pathDest)
	}

	return conn.CopySingleFileToRemote(pathSource, pathDest)
}
