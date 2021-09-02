// +build amd64,darwin linux

package linux

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	defaultPath       = "/boot/grub"
	defaultMountPoint = "/boot"
)

func TestMountUnixAll(t *testing.T) {
	TestFindMountPoint(t)
}

func TestFindMountPoint(t *testing.T) {
	asst := assert.New(t)

	mountPoint, err := FindMountPoint(defaultPath)
	asst.Nil(err, "test FindMountPoint() failed")
	asst.Equal(defaultMountPoint, mountPoint, "test FindMountPoint() failed")
	t.Logf("path: %s, mount point: %s", defaultPath, mountPoint)
}
