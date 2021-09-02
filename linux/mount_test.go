package linux

import (
	"testing"

	"github.com/romberli/go-util/constant"
	"github.com/stretchr/testify/assert"
)

const (
	// mount point
	data   = "/data"
	mysql  = "/data/mysql"
	binlog = "/data/mysql/binlog"
	other  = "/other/dir"
	// path
	mysqlDataDir = "/data/mysql/data/"
	someOtherDir = "/some/other/dir"
)

func TestMountAll(t *testing.T) {
	TestMatchMountPoint(t)
}

func TestMatchMountPoint(t *testing.T) {
	asst := assert.New(t)

	mountPoints := []string{constant.RootDir, data, mysql, binlog, other}

	mountPoint, err := MatchMountPoint(mysqlDataDir, mountPoints)
	asst.Nil(err, "test MatchMountPoint() failed")
	asst.Equal(mysql, mountPoint, "test MatchMountPoint() failed")

	mountPoint, err = MatchMountPoint(someOtherDir, mountPoints)
	asst.Nil(err, "test MatchMountPoint() failed")
	asst.Equal(constant.RootDir, mountPoint, "test MatchMountPoint() failed")
}
