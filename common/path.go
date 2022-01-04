package common

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/romberli/go-util/constant"
)

// IsAbs returns if given path is an absolute path,
// this function is not related to the platform,
// it works for both unix like and windows paths
func IsAbs(path string) bool {
	if filepath.IsAbs(path) {
		return true
	}
	if IsAbsUnix(path) {
		return true
	}
	if IsAbsWindows(path) {
		return true
	}
	return false
}

// IsAbsUnix returns if given path is an absolute path of unix like system
func IsAbsUnix(path string) bool {
	return strings.HasPrefix(path, constant.SlashString)
}

// IsAbsUnix returns if given path is an absolute path of windows system
func IsAbsWindows(path string) bool {
	if len(path) < 2 {
		return false
	}

	c := path[constant.ZeroInt]
	if ('a' <= c && c <= 'z' || 'A' <= c && c <= 'Z') && path[1] == ':' {
		if len(path) == 2 {
			return true
		}

		if path[2] == '/' || path[2] == '\\' {
			return true
		}
	}

	return false
}

// GetParentDir returns the parent directory of the given path, it will not change the separate character
func GetParentDir(absPath string) (string, error) {
	if !IsAbs(absPath) {
		return constant.EmptyString, errors.New(fmt.Sprintf("path must be an absolute path, %s is not valid", absPath))
	}

	for i := len(absPath) - 1; i > constant.ZeroInt; i-- {
		if absPath[i] == '/' || absPath[i] == '\\' {
			return absPath[:i], nil
		}
	}

	return constant.SlashString, nil
}
