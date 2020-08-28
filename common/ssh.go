package common

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type CloseConn interface {
	CloseConn()
}

const (
	DefaultStringZeroValue = ""
	DefaultIntZeroValue    = 0
	DefaultSSHTimeout      = 10 * time.Second
	DefaultSSHPortNum      = 22
	DefaultSSHUserName     = "root"
	DefaultSSHUserPass     = "shit"
	DefaultByteBufferSize  = 1024 * 1024 // 1MB

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
		DefaultSSHPortNum,
		DefaultSSHUserName,
		DefaultSSHUserPass,
	}
}

type MySSHConn struct {
	MyConn
	SSHClient  *ssh.Client
	SSHSession *ssh.Session
}

// NewMySSHConn returns *MySSHConn and error, first argument is mandatory which presents host ip,
// and the 3 flowing optional arguments which should be in exact order of port number, user name and user password
func NewMySSHConn(hostIp string, in ...interface{}) (sshConn *MySSHConn, err error) {
	var (
		myConn       *MyConn
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sshSession   *ssh.Session
	)

	argLen := len(in)
	switch argLen {
	case 0:
		hostIp = strings.TrimSpace(hostIp)
		if hostIp == "" {
			return nil, errors.New("host ip could not be empty")
		}

		myConn = NewMyConnWithDefaultValue(hostIp)
	case 3:
		var (
			portNumValue  int
			userNameValue string
			userPassValue string
		)

		portNum := in[0]
		userName := in[1]
		userPass := in[2]

		switch portNum.(type) {
		case nil:
			portNumValue = DefaultSSHPortNum
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			portNumValue = portNum.(int)
		default:
			return nil, errors.New(
				fmt.Sprintf("port number must be integer type instead of %s",
					reflect.TypeOf(portNum).Name()))
		}

		switch userName.(type) {
		case nil:
			userNameValue = DefaultSSHUserName
		case string:
			userNameValue = strings.TrimSpace(userName.(string))
			if userNameValue == "" {
				userNameValue = DefaultSSHUserName
			}
		default:
			return nil, errors.New(
				fmt.Sprintf("user name must be string type instead of %s",
					reflect.TypeOf(portNum).Name()))
		}

		switch userPass.(type) {
		case nil:
			userPassValue = DefaultSSHUserPass
		case string:
			userPassValue = strings.TrimSpace(userPass.(string))
			if userPassValue == "" {
				userPassValue = DefaultSSHUserPass
			}
		default:
			return nil, errors.New(
				fmt.Sprintf("user password must be string type instead of %s",
					reflect.TypeOf(portNum).Name()))
		}

		myConn = NewMyConn(hostIp, portNumValue, userNameValue, userPassValue)
	default:
		return nil, errors.New(fmt.Sprintf("optional argument number must be 0 or 3 instead of %d", argLen))
	}

	// get auth method
	auth = append(auth, ssh.Password(myConn.UserPass))

	// 这个是问你要不要验证远程主机，以保证安全性。这里不验证
	hostKeyCallBack := func(host string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	clientConfig = &ssh.ClientConfig{
		User:            myConn.UserName,
		Auth:            auth,
		Timeout:         DefaultSSHTimeout,
		HostKeyCallback: hostKeyCallBack,
	}

	// connect to ssh
	addr = fmt.Sprintf("%s:%d", myConn.HostIp, myConn.PortNum)
	sshClient, err = ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return nil, err
	}

	// create ssh session
	sshSession, err = sshClient.NewSession()
	if err != nil {
		return nil, err
	}

	sshConn = &MySSHConn{
		*myConn,
		sshClient,
		sshSession,
	}

	return sshConn, nil
}

func NewMySSHConnWithDefaultValue(hostIp string) (sshConn *MySSHConn, err error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sshSession   *ssh.Session
	)

	myConn := NewMyConnWithDefaultValue(hostIp)

	// get auth method
	auth = append(auth, ssh.Password(myConn.UserPass))

	// 这个是问你要不要验证远程主机，以保证安全性。这里不验证
	hostKeyCallBack := func(host string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	clientConfig = &ssh.ClientConfig{
		User:            myConn.UserName,
		Auth:            auth,
		Timeout:         DefaultSSHTimeout,
		HostKeyCallback: hostKeyCallBack,
	}

	// connect to ssh
	addr = fmt.Sprintf("%s:%d", myConn.HostIp, myConn.PortNum)
	sshClient, err = ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return nil, err
	}

	// create ssh session
	sshSession, err = sshClient.NewSession()
	if err != nil {
		return nil, err
	}

	sshConn = &MySSHConn{
		*myConn,
		sshClient,
		sshSession,
	}

	return sshConn, nil
}

func (conn *MySSHConn) ExecuteCommand(cmd string) (result int, stdOut string, stdErr string, err error) {
	var stdOutBuffer, stdErrBuffer bytes.Buffer

	conn.SSHSession.Stdout = &stdOutBuffer
	conn.SSHSession.Stderr = &stdErrBuffer

	// run command
	err = conn.SSHSession.Run(cmd)
	if err != nil {
		result = 1
	}

	return result, stdOutBuffer.String(), stdErrBuffer.String(), err
}

func (conn *MySSHConn) GetHostName() (hostName string, err error) {
	_, hostName, _, err = conn.ExecuteCommand("hostname")

	return hostName, err
}

func (conn *MySSHConn) CloseConn() error {
	return conn.SSHSession.Close()
}

type MySftpConn struct {
	MySSHConn
	SftpClient *sftp.Client
}

// NewMySftpConn returns *MySftpConn and error, first argument is mandatory which presents host ip,
// and the 3 flowing optional arguments which should be in exact order of port number, user name and user password
func NewMySftpConn(hostIp string, in ...interface{}) (sftpConn *MySftpConn, err error) {
	var (
		mySSHConn  *MySSHConn
		sftpClient *sftp.Client
	)

	if mySSHConn, err = NewMySSHConn(hostIp, in...); err != nil {
		return nil, err
	}

	if sftpClient, err = sftp.NewClient(mySSHConn.SSHClient); err != nil {
		return nil, err
	}

	sftpConn = &MySftpConn{
		*mySSHConn,
		sftpClient,
	}

	return sftpConn, nil
}

func NewMySftpConnWithDefaultValue(hostIp string) (sftpConn *MySftpConn, err error) {
	var mySSHConn *MySSHConn
	var sftpClient *sftp.Client

	mySSHConn, err = NewMySSHConn(hostIp)
	if err != nil {
		return nil, err
	}

	sftpClient, err = sftp.NewClient(mySSHConn.SSHClient)
	if err != nil {
		return nil, err
	}

	sftpConn = &MySftpConn{
		*mySSHConn,
		sftpClient,
	}

	return sftpConn, nil
}

func (conn *MySftpConn) CopyFile(fileSource io.Reader, fileDest io.Writer, bufferSize int) (err error) {
	var n int

	if bufferSize <= 0 {
		bufferSize = DefaultByteBufferSize
	}

	buf := make([]byte, bufferSize)

	for {
		n, err = fileSource.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		_, err = fileDest.Write(buf[0:n])
		if err != nil {
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

	fileSource, err = conn.SftpClient.Open(fileNameSource)
	if err != nil {
		return err
	}
	defer func() { _ = fileSource.Close() }()

	fileDest, err = os.Create(fileNameDest)
	if err != nil {
		return err
	}
	defer func() { _ = fileDest.Close() }()

	err = conn.CopyFile(fileSource, fileDest, DefaultByteBufferSize)
	if err != nil {
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

	fileSource, err = os.Open(fileNameSource)
	if err != nil {
		return err
	}
	defer func() { _ = fileSource.Close() }()

	fileDest, err = conn.SftpClient.Create(fileNameDest)
	if err != nil {
		return err
	}
	defer func() { _ = fileDest.Close() }()

	err = conn.CopyFile(fileSource, fileDest, DefaultByteBufferSize)
	if err != nil {
		return err
	}

	return nil
}

func (conn *MySftpConn) CopyFileListFromRemote(fileListSource []string, FileDirDest string) (err error) {
	var exists bool

	FileDirDest = strings.TrimSpace(FileDirDest)

	if FileDirDest == DefaultStringZeroValue {
		return errors.New("file destination directory should NOT an empty string.")
	}

	exists, err = PathExistsLocal(FileDirDest)
	if err != nil {
		return err
	}

	if !exists {
		_, err = os.Create(FileDirDest)
		if err != nil {
			return err
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

func (conn *MySftpConn) CopyFileListFromRemoteWithNewName(fileListSource []string, FileListDest []string) (err error) {
	if len(fileListSource) != len(FileListDest) {
		return errors.New("the length of source and destination list MUST be exactly same")
	}

	for i, fileNameSource := range fileListSource {
		fileNameSource = strings.TrimSpace(fileNameSource)

		fileNameDest := FileListDest[i]
		fileNameDest = strings.TrimSpace(fileNameDest)

		if fileNameDest == DefaultStringZeroValue {
			return errors.New("destination file name should not be an empty string")
		}

		err = conn.CopySingleFileFromRemote(fileNameSource, fileNameDest)
		if err != nil {
			return err
		}
	}

	return nil
}

func (conn *MySftpConn) CloseConn() (err error) {
	return conn.SftpClient.Close()
}
