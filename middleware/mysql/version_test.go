package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testVersionStr = "5.7.21-log"

func initVersion() Version {
	return NewVersion(5, 7, 21)
}

func equal(v1, v2 Version) bool {
	return v1.GetMajor() == v2.GetMajor() &&
		v1.GetMinor() == v2.GetMinor() &&
		v1.GetRelease() == v2.GetRelease() &&
		v1.String() == v2.String()
}

func TestParse(t *testing.T) {
	asst := assert.New(t)

	v1 := initVersion()
	v2, err := Parse(testVersionStr)
	asst.Nil(err, "test Parse() failed")
	asst.True(equal(v1, v2), "test Parse() failed")
}
