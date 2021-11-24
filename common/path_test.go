package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testUnixAbsPath     = "/data/mysql/data/"
	testUnixRelPath     = "data/mysql/"
	testWindowsAbsPathA = `C:\data\mysql\data\`
	testWindowsAbsPathB = `C:/data/mysql/data/`
	testWindowsRelPath  = `data\mysql\data\`
)

func TestAll(t *testing.T) {
	TestIsAbs(t)
	TestIsAbsUnix(t)
	TestIsAbsWindows(t)
}
func TestIsAbs(t *testing.T) {
	asst := assert.New(t)

	asst.True(IsAbs(testUnixAbsPath), "test IsAbs() failed")
	asst.False(IsAbs(testUnixRelPath), "test IsAbs() failed")
	asst.True(IsAbs(testWindowsAbsPathA), "test IsAbs() failed")
	asst.True(IsAbs(testWindowsAbsPathB), "test IsAbs() failed")
	asst.False(IsAbs(testWindowsRelPath), "test IsAbs() failed")
}

func TestIsAbsUnix(t *testing.T) {
	asst := assert.New(t)

	asst.True(IsAbsUnix(testUnixAbsPath), "test IsAbsUnix() failed")
	asst.False(IsAbsUnix(testUnixRelPath), "test IsAbsUnix() failed")
	asst.False(IsAbsUnix(testWindowsAbsPathA), "test IsAbsUnix() failed")
	asst.False(IsAbsUnix(testWindowsAbsPathB), "test IsAbsUnix() failed")
	asst.False(IsAbsUnix(testWindowsRelPath), "test IsAbsUnix() failed")
}

func TestIsAbsWindows(t *testing.T) {
	asst := assert.New(t)

	asst.False(IsAbsWindows(testUnixAbsPath), "test IsAbsWindows() failed")
	asst.False(IsAbsWindows(testUnixRelPath), "test IsAbsWindows() failed")
	asst.True(IsAbsWindows(testWindowsAbsPathA), "test IsAbsWindows() failed")
	asst.True(IsAbsWindows(testWindowsAbsPathB), "test IsAbsWindows() failed")
	asst.False(IsAbsWindows(testWindowsRelPath), "test IsAbsWindows() failed")
}
