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
	DefaultEstimateLineSize = 1024
	MinStartPosition        = 0
)

// SyscallMode returns file mode which could be used at syscall
func SyscallMode(fileMode os.FileMode) (fileModeSys uint32) {
	fileModeSys |= uint32(fileMode.Perm())

	if fileMode&os.ModeSetuid != constant.ZeroInt {
		fileModeSys |= syscall.S_ISUID
	}
	if fileMode&os.ModeSetgid != constant.ZeroInt {
		fileModeSys |= syscall.S_ISGID
	}
	if fileMode&os.ModeSticky != constant.ZeroInt {
		fileModeSys |= syscall.S_ISVTX
	}

	return
}

// PathExists returns if given path exists
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, errors.Trace(err)
}

// IsDir returns if given path is a directory or not
func IsDir(path string) (isDir bool, err error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, errors.Trace(err)
	}

	return fileInfo.IsDir(), nil
}

// ReadDir returns subdirectories and files of given directory on the host, it returns a slice of os.FileInfo
func ReadDir(dirName string) (fileInfoList []os.FileInfo, err error) {
	isDir, err := IsDir(dirName)
	if err != nil {
		return nil, err
	}
	if !isDir {
		return nil, errors.New(fmt.Sprintf("it's NOT a directory. dir name: %s", dirName))
	}

	file, err := os.Open(dirName)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer func() { _ = file.Close() }()

	fileInfoList, err = file.Readdir(constant.ZeroInt)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return fileInfoList, nil
}

// IsEmptyDir returns if given directory is empty or not
func IsEmptyDir(dirName string) (isEmpty bool, err error) {
	fileInfoList, err := ReadDir(dirName)
	if err != nil {
		return false, err
	}
	if len(fileInfoList) == constant.ZeroInt {
		isEmpty = true
	}

	return isEmpty, nil
}

// GetPathDirMapLocal reads all subdirectories and files of given directory and calculate the relative path of rootPath,
// then map the absolute path of subdirectory names and file names as keys, relative paths as values to fileDirMap
func GetPathDirMapLocal(dirName, rootPath string) (map[string]string, error) {
	pathDirMap := make(map[string]string)

	err := getPathDirMapLocal(pathDirMap, dirName, rootPath)
	if err != nil {
		return nil, err
	}

	return pathDirMap, nil
}

func getPathDirMapLocal(pathDirMap map[string]string, dirName, rootPath string) error {
	pathInfoList, err := ReadDir(dirName)
	if err != nil {
		return err
	}
	if len(pathInfoList) == constant.ZeroInt {
		// it's an empty directory
		pathDirMap[dirName] = constant.EmptyString
	}

	for _, pathInfo := range pathInfoList {
		pathName := pathInfo.Name()
		pathNameAbs := filepath.Join(dirName, pathName)

		if pathInfo.IsDir() {
			err = getPathDirMapLocal(pathDirMap, pathNameAbs, rootPath)
			if err != nil {
				return err
			}
		} else {
			pathRel, err := filepath.Rel(rootPath, pathNameAbs)
			if err != nil {
				return errors.Trace(err)
			}

			pathDirMap[pathNameAbs] = pathRel
		}
	}

	return nil
}

// GetFileNameDest returns the destination file name
func GetFileNameDest(fileNameSource, dirDest string) string {
	fileNameBase := filepath.Base(fileNameSource)

	return filepath.Join(dirDest, fileNameBase)
}

// TailN try to get the latest n line of the file.
func TailN(fileName string, n int) (lines []string, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer func() { _ = file.Close() }()

	estimateLineSize := DefaultEstimateLineSize

	stat, err := os.Stat(fileName)
	if err != nil {
		return nil, errors.Trace(err)
	}

	start := int(stat.Size()) - n*estimateLineSize
	if start < MinStartPosition {
		start = MinStartPosition
	}

	_, err = file.Seek(int64(start), MinStartPosition /*means relative to the origin of the file*/)
	if err != nil {
		return nil, errors.Trace(err)
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
