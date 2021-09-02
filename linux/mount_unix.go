// +build amd64,darwin linux

package linux

import (
	"os"
	"path/filepath"
	"syscall"

	"github.com/romberli/go-util/constant"
)

// FindMountPoint returns the mount point of the given path,
// note that this function only works on unix-like system
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
