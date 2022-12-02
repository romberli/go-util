package mysql

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pingcap/errors"
	"github.com/romberli/go-util/constant"
)

const (
	defaultSeparator     = "."
	mysqlVersionSliceLen = 3
	LessThan             = -1
	Equal                = 0
	GreaterThan          = 1
)

var Version8 = &version{8, 0, 0}

type Version interface {
	// GetMajor returns major
	GetMajor() int
	// GetMinor returns minor
	GetMinor() int
	// GetRelease returns release
	GetRelease() int
	// Compare compares two versions
	Compare(other Version) int
	// LessThan returns true if current version is less than other version
	LessThan(other Version) bool
	// Equal returns true if current version is equal to other version
	Equal(other Version) bool
	// GreaterThan returns true if current version is greater than other version
	GreaterThan(other Version) bool
	// IsMySQL8 returns true if current version is greater than or equal to 8.0.0
	IsMySQL8() bool
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
	vSlice := strings.Split(v, constant.DashString)
	versionSlice := strings.Split(vSlice[constant.ZeroInt], defaultSeparator)
	if len(versionSlice) > mysqlVersionSliceLen {
		return nil, errors.New(fmt.Sprintf("%s is not a valid mysql version string.", v))
	}

	ver := &version{}

	for i, vStr := range versionSlice {
		vInt, err := strconv.Atoi(vStr)
		if err != nil {
			return nil, errors.Trace(err)
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

func (v *version) Compare(other Version) int {
	if v.GetMajor() > other.GetMajor() {
		return GreaterThan
	} else if v.major < other.GetMajor() {
		return LessThan
	}

	if v.GetMinor() > other.GetMinor() {
		return GreaterThan
	} else if v.GetMinor() < other.GetMinor() {
		return LessThan
	}

	if v.GetRelease() > other.GetRelease() {
		return GreaterThan
	} else if v.GetRelease() < other.GetRelease() {
		return LessThan
	}

	return Equal
}

func (v *version) LessThan(other Version) bool {
	return v.Compare(other) == LessThan
}

func (v *version) Equal(other Version) bool {
	return v.Compare(other) == Equal
}

func (v *version) GreaterThan(other Version) bool {
	return v.Compare(other) == GreaterThan
}

func (v *version) IsMySQL8() bool {
	return !v.LessThan(Version8)
}

// String returns string format of version
func (v *version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.release)
}
