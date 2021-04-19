package mysql

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/romberli/go-util/constant"
)

const (
	defaultSeparator     = "."
	mysqlVersionSliceLen = 3
)

type Version interface {
	// GetMajor returns major
	GetMajor() int
	// GetMinor returns minor
	GetMinor() int
	// GetRelease returns release
	GetRelease() int
	// String returns string format of version
	String() string
}

type version struct {
	major   int
	minor   int
	release int
}

// NewVersion returns a new instance of Version
func NewVersion(major, minor, release int) Version {
	return &version{
		major,
		minor,
		release,
	}
}

// Parse parses given version string and returns a new instance of Version
func Parse(v string) (Version, error) {
	return parse(v)
}

// parse parses given version string and returns a new instance of version
func parse(v string) (*version, error) {
	vSlice := strings.Split(v, "-")
	versionSlice := strings.Split(vSlice[constant.ZeroInt], defaultSeparator)
	if len(versionSlice) > mysqlVersionSliceLen {
		return nil, errors.New(fmt.Sprintf("%s is not a valid mysql version string.", v))
	}

	ver := &version{}

	for i, vStr := range versionSlice {
		vInt, err := strconv.Atoi(vStr)
		if err != nil {
			return nil, err
		}

		switch i {
		case 0:
			ver.major = vInt
		case 1:
			ver.minor = vInt
		case 2:
			ver.release = vInt
		}
	}

	return ver, nil
}

// GetMajor returns major
func (v *version) GetMajor() int {
	return v.major
}

// GetMinor returns minor
func (v *version) GetMinor() int {
	return v.minor
}

// GetRelease returns release
func (v *version) GetRelease() int {
	return v.release
}

// String returns string format of version
func (v *version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.release)
}
