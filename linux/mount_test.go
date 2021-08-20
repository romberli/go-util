package linux

import (
	"testing"

	"github.com/romberli/go-util/constant"
	"github.com/stretchr/testify/assert"
)

const (
	defaultPath       = "/boot/grub"
	defaultMountPoint = "/boot"

	// mount point
	data   = "/data"
	mysql  = "/data/mysql"
	binlog = "/data/mysql/binlog"
	other  = "/other/dir"
	// path
	mysqlDataDir = "/data/mysql/data"
	someOtherDir = "/some/other/dir"
)

func TestAll(t *testing.T) {
	TestFindMountPoint(t)
	TestMatchMountPoint(t)
}

func TestFindMountPoint(t *testing.T) {
	asst := assert.New(t)

	mountPoint, err := FindMountPoint(defaultPath)
	asst.Nil(err, "test FindMountPoint() failed")
	asst.Equal(defaultMountPoint, mountPoint, "test FindMountPoint() failed")
	t.Logf("path: %s, mount point: %s", defaultPath, mountPoint)
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
