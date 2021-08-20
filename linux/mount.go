package linux

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/romberli/go-util/constant"
)

// FindMountPoint returns the mount point of the given path
func FindMountPoint(path string) (string, error) {
	pathStat, err := os.Stat(path)
	if err != nil {
		return constant.EmptyString, err
	}

	pathStat.Sys()
	dev := pathStat.Sys().(*syscall.Stat_t).Dev

	for path != constant.RootDir {
		parent := filepath.Dir(path)

		parentStat, err := os.Stat(parent)
		if err != nil {
			return constant.EmptyString, err
		}

		parentDev := parentStat.Sys().(*syscall.Stat_t).Dev

		if dev != parentDev {
			break
		}

		path = parent
	}

	return path, nil
}

// MatchMountPoint matches mount point of given path in the mount point slice,
// if nothing matched, it returns "/" as default mount point
func MatchMountPoint(path string, mountPoints []string) (string, error) {
	if !filepath.IsAbs(path) {
		return constant.EmptyString, errors.New(fmt.Sprintf("path must be an absolute path, %s is not valid", path))
	}
	if path == constant.RootDir {
		return constant.RootDir, nil
	}

	for _, mountPoint := range mountPoints {
		if mountPoint == path {
			return mountPoint, nil
		}
	}

	return MatchMountPoint(filepath.Dir(path), mountPoints)
}
