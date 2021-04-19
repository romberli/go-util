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
	GetMajor() int
	GetMinor() int
	GetRelease() int
	String() string
}

type version struct {
	major   int
	minor   int
	release int
}

func NewVersion(major, minor, release int) Version {
	return &version{
		major,
		minor,
		release,
	}
}

func Parse(v string) (Version, error) {
	return parse(v)
}

func parse(v string) (*version, error) {
	vSlice := strings.Split(v, "-")
	versionSlice := strings.Split(vSlice[constant.ZeroInt], defaultSeparator)
	if len(versionSlice) > mysqlVersionSliceLen {
		return nil, errors.New(fmt.Sprintf("%s is not a valid mysql version string.", v))
	}

	ver := &version{}

	for i, vStr := range versionSlice {
		vNum, err := strconv.Atoi(vStr)
		if err != nil {
			return nil, err
		}

		switch i {
		case 0:
			ver.major = vNum
		case 1:
			ver.minor = vNum
		case 2:
			ver.release = vNum
		}
	}

	return ver, nil
}

func (v *version) GetMajor() int {
	return v.major
}

func (v *version) GetMinor() int {
	return v.minor
}

func (v *version) GetRelease() int {
	return v.release
}

func (v *version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.release)
}
