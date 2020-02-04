package common

import (
	"bytes"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"path"
	"time"
)

type CloseConn interface{}

const (
	portNumSshDefault      = 22
	userNameSshDefault     = "root"
	userPassSshDefault     = "shit"
	bytesBufferSizeDefault = 1024 * 1024 // 1MB
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
	Conn       MyConn
	SshClient  *ssh.Client
	SshSession *ssh.Session
}

func NewMySshConn(myConn MyConn) (sshConn *MySshConn, err error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sshSession   *ssh.Session
	)

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
		Timeout:         30 * time.Second,
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
		myConn,
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
	Conn       MySshConn
	SftpClient *sftp.Client
}

func NewMySftpConn(mySshConn MySshConn) (sftpConn *MySftpConn, err error) {
	var sftpClient *sftp.Client

	if sftpClient, err = sftp.NewClient(mySshConn.SshClient); err != nil {
		return nil, err
	}

	sftpConn = &MySftpConn{
		mySshConn,
		sftpClient,
	}

	return sftpConn, nil
}

func (conn *MySftpConn) CopySingleFileFromRemote(fileNameSource string, fileNameDest string) (err error) {
	var (
		n          int
		fileDest   *os.File
		fileSource *sftp.File
	)

	if fileSource, err = conn.SftpClient.Create(fileNameSource); err != nil {
		return err
	}
	defer fileSource.Close()

	if fileDest, err = os.Open(fileNameSource); err != nil {
		return err
	}
	defer fileDest.Close()

	buf := make([]byte, bytesBufferSizeDefault)

	for {
		if n, err = fileSource.Read(buf); err != nil {
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

func (conn *MySftpConn) CopySingleFileToRemote(fileNameSource string, fileNameDest string) (err error) {
	var (
		n          int
		fileSource *os.File
		fileDest   *sftp.File
	)

	if fileSource, err = os.Open(fileNameSource); err != nil {
		return err
	}
	defer fileSource.Close()

	if fileDest, err = conn.SftpClient.Create(fileNameDest); err != nil {
		return err
	}
	defer fileSource.Close()

	buf := make([]byte, bytesBufferSizeDefault)

	for {
		if n, err = fileSource.Read(buf); err != nil {
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

func (conn *MySftpConn) CopyFileListFromRemote(fileListSource []string, FileDirDest string) (err error) {
	for _, fileNameSource := range fileListSource {
		fileNameDest := path.Base(fileNameSource)
		if err = conn.CopySingleFileFromRemote(fileNameSource, path.Join(FileDirDest, fileNameDest)); err != nil {
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
