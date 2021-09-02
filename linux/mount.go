package linux

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/romberli/go-util/constant"
)

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
