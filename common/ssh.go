package common

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type CloseConn interface {
	CloseConn()
}

const (
	portNumSshDefault     = 22
	userNameSshDefault    = "root"
	userPassSshDefault    = "shit"
	byteBufferSizeDefault = 1024 * 1024 // 1MB

	copyFromDefault int = iota
	copyToDefault
)

type MyConn struct {
	HostIp   string
	PortNum  int
	UserName string
	UserPass string
}

func NewMyConn(hostIp string, portNum int, userName string, userPass string) (conn *MyConn) {
	return &MyConn{
		hostIp,
		portNum,
		userName,
		userPass,
	}
}

func NewMyConnWithDefaultValue(hostIp string) (conn *MyConn) {
	return &MyConn{
		hostIp,
		portNumSshDefault,
		userNameSshDefault,
		userPassSshDefault,
	}
}

type MySshConn struct {
	MyConn
	SshClient  *ssh.Client
	SshSession *ssh.Session
}

func NewMySshConn(hostIp string) (sshConn *MySshConn, err error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sshSession   *ssh.Session
	)

	myConn := NewMyConnWithDefaultValue(hostIp)

	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(myConn.UserPass))

	// 这个是问你要不要验证远程主机，以保证安全性。这里不验证
	hostKeyCallBack := func(host string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	clientConfig = &ssh.ClientConfig{
		User:            myConn.UserName,
		Auth:            auth,
		Timeout:         10 * time.Second,
		HostKeyCallback: hostKeyCallBack,
	}

	// connect to ssh
	addr = fmt.Sprintf("%s:%d", myConn.HostIp, myConn.PortNum)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create ssh session
	if sshSession, err = sshClient.NewSession(); err != nil {
		return nil, err
	}

	sshConn = &MySshConn{
		*myConn,
		sshClient,
		sshSession,
	}

	return sshConn, nil
}

func (conn *MySshConn) ExecuteCommand(cmd string) (result int, stdOut string, stdErr string, err error) {
	var stdOutBuffer, stdErrBuffer bytes.Buffer

	conn.SshSession.Stdout = &stdOutBuffer
	conn.SshSession.Stderr = &stdErrBuffer

	// run command
	conn.SshSession.Run(cmd)

	return result, stdOutBuffer.String(), stdErrBuffer.String(), err
}

func (conn *MySshConn) GetHostName() (hostName string, err error) {
	if _, hostName, _, err = conn.ExecuteCommand("hostname"); err != nil {
		return "", err
	}

	return hostName, err
}

func (conn *MySshConn) CloseConn() (err error) {
	if err = conn.SshSession.Close(); err != nil {
		return err
	}

	return nil
}

type MySftpConn struct {
	MySshConn
	SftpClient *sftp.Client
}

func NewMySftpConn(hostIp string) (sftpConn *MySftpConn, err error) {
	var mySshConn *MySshConn
	var sftpClient *sftp.Client

	if mySshConn, err = NewMySshConn(hostIp); err != nil {
		return nil, err
	}

	if sftpClient, err = sftp.NewClient(mySshConn.SshClient); err != nil {
		return nil, err
	}

	sftpConn = &MySftpConn{
		*mySshConn,
		sftpClient,
	}

	return sftpConn, nil
}

func (conn *MySftpConn) CopyFile(fileSource io.Reader, fileDest io.Writer, bufferSize int) (err error) {
	var n int

	if bufferSize <= 0 {
		bufferSize = byteBufferSizeDefault
	}

	buf := make([]byte, bufferSize)

	for {
		if n, err = fileSource.Read(buf); err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		if _, err = fileDest.Write(buf[0:n]); err != nil {
			return err
		}
	}

	return nil
}

func (conn *MySftpConn) CopySingleFileFromRemote(fileNameSource string, fileNameDest string) (err error) {
	var (
		fileDest   *os.File
		fileSource *sftp.File
	)

	fileNameSource = strings.TrimSpace(fileNameSource)
	fileNameDest = strings.TrimSpace(fileNameDest)

	if fileNameDest == "" {
		fileNameDest = fileNameSource
	}

	if fileSource, err = conn.SftpClient.Open(fileNameSource); err != nil {
		return err
	}
	defer fileSource.Close()

	if fileDest, err = os.Create(fileNameDest); err != nil {
		return err
	}
	defer fileDest.Close()

	if err = conn.CopyFile(fileSource, fileDest, byteBufferSizeDefault); err != nil {
		return err
	}

	return nil
}

func (conn *MySftpConn) CopySingleFileToRemote(fileNameSource string, fileNameDest string) (err error) {
	var (
		fileSource *os.File
		fileDest   *sftp.File
	)

	fileNameSource = strings.TrimSpace(fileNameSource)
	fileNameDest = strings.TrimSpace(fileNameDest)

	if fileNameDest == "" {
		fileNameDest = fileNameSource
	}

	if fileSource, err = os.Open(fileNameSource); err != nil {
		return err
	}
	defer fileSource.Close()

	if fileDest, err = conn.SftpClient.Create(fileNameDest); err != nil {
		return err
	}
	defer fileSource.Close()

	if err = conn.CopyFile(fileSource, fileDest, byteBufferSizeDefault); err != nil {
		return err
	}

	return nil
}

func (conn *MySftpConn) CopyFileListFromRemote(fileListSource []string, FileDirDest string) (err error) {
	var exists bool

	FileDirDest = strings.TrimSpace(FileDirDest)

	if FileDirDest == "" {
		return errors.New("file destination directory should NOT an empty string.")
	}

	if exists, err = PathExistsLocal(FileDirDest); err != nil {
		return err
	}

	if !exists {
		if _, err = os.Create(FileDirDest); err != nil {
			return err
		}
	}

	for _, fileNameSource := range fileListSource {
		fileNameSource = strings.TrimSpace(fileNameSource)
		fileNameDest := path.Base(fileNameSource)

		if err = conn.CopySingleFileFromRemote(fileNameSource, path.Join(FileDirDest, fileNameDest)); err != nil {
			return err
		}
	}

	return nil
}

func (conn *MySftpConn) CopyFileListFromRemoteWithNewName(fileListSource []string, FileListDest []string) (err error) {
	if len(fileListSource) != len(FileListDest) {
		return errors.New("the length of source and destination list MUST be exactly same")
	}

	for i, fileNameSource := range fileListSource {
		fileNameSource = strings.TrimSpace(fileNameSource)

		fileNameDest := FileListDest[i]
		fileNameDest = strings.TrimSpace(fileNameDest)

		if fileNameDest == "" {
			return errors.New("destination file name should not be an empty string")
		}

		if err = conn.CopySingleFileFromRemote(fileNameSource, fileNameDest); err != nil {
			return err
		}
	}

	return nil
}

func (conn *MySftpConn) CloseConn() (err error) {
	if err = conn.SftpClient.Close(); err != nil {
		return err
	}

	return nil
}
