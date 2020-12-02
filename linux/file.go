package linux

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/pingcap/errors"

	"github.com/romberli/go-util/constant"
)

const (
	HostNameCommand         = "hostname"
	LsCommand               = "ls"
	DefaultEstimateLineSize = 1024
	MinStartPosition        = 0
)

// SyscallMode returns file mode which could be used at syscall
func SyscallMode(fileMode os.FileMode) (fileModeSys uint32) {
	fileModeSys |= uint32(fileMode.Perm())

	if fileMode&os.ModeSetuid != 0 {
		fileModeSys |= syscall.S_ISUID
	}
	if fileMode&os.ModeSetgid != 0 {
		fileModeSys |= syscall.S_ISGID
	}
	if fileMode&os.ModeSticky != 0 {
		fileModeSys |= syscall.S_ISVTX
	}

	return
}

// IsDir returns if given path is a directory or not
func IsDir(path string) (isDir bool, err error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer func() { _ = file.Close() }()

	fileInfo, err := file.Stat()
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}

// ReadDir returns subdirectories and files of given directory on the remote host, it returns a slice of os.FileInfo
func Readdir(dirName string) (fileInfoList []os.FileInfo, err error) {
	isDir, err := IsDir(dirName)
	if err != nil {
		return nil, err
	}
	if !isDir {
		return nil, errors.New(fmt.Sprintf("it's NOT a directory. dir name: %s", dirName))
	}

	file, err := os.Open(dirName)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	fileInfoList, err = file.Readdir(int(constant.MinUInt))
	if err != nil {
		return nil, err
	}

	return fileInfoList, nil
}

// IsEmptyDir returns if given directory is empty or not
func IsEmptyDir(dirName string) (isEmpty bool, err error) {
	fileInfoList, err := Readdir(dirName)
	if err != nil {
		return false, err
	}
	if len(fileInfoList) == 0 {
		isEmpty = true
	}

	return isEmpty, nil
}

// GetFileDirMapLocal reads all subdirectories and files of given directory and calculate the relative path of rootPath,
// then map the absolute path of subdirectory names and file names as keys, relative paths as values to fileDirMap
func GetFileDirMapLocal(fileDirMap map[string]string, dirName, rootPath string) (err error) {
	fileInfoList, err := Readdir(dirName)
	if err != nil {
		return err
	}
	if len(fileInfoList) == 0 {
		// it's an empty directory
		fileDirMap[dirName] = constant.EmptyString
	}

	for _, fileInfo := range fileInfoList {
		fileName := fileInfo.Name()
		fileNameAbs := filepath.Join(dirName, fileName)

		if fileInfo.IsDir() {
			err = GetFileDirMapLocal(fileDirMap, fileNameAbs, rootPath)
			if err != nil {
				return err
			}
		} else {
			fileNameRel, err := filepath.Rel(rootPath, fileNameAbs)
			if err != nil {
				return err
			}

			fileDirMap[fileNameAbs] = fileNameRel
		}
	}

	return nil
}

// GetFileNameDest returns the destination file name
func GetFileNameDest(fileNameSource, dirDest string) string {
	fileNameBase := filepath.Base(fileNameSource)

	return filepath.Join(dirDest, fileNameBase)
}

// TailN try get the latest n line of the file.
func TailN(fileName string, n int) (lines []string, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.AddStack(err)
	}
	defer func() { _ = file.Close() }()

	estimateLineSize := DefaultEstimateLineSize

	stat, err := os.Stat(fileName)
	if err != nil {
		return nil, errors.AddStack(err)
	}

	start := int(stat.Size()) - n*estimateLineSize
	if start < MinStartPosition {
		start = MinStartPosition
	}

	_, err = file.Seek(int64(start), MinStartPosition /*means relative to the origin of the file*/)
	if err != nil {
		return nil, errors.AddStack(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}

	return
}
