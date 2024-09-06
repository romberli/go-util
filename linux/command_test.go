package linux

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	var (
		err     error
		workDir string
		cmd     string
		stdout  string
	)

	asst := assert.New(t)

	workDir = "/Users/romber/source_code/go/src/github.com/romberli/go-util"
	cmd = "ls -l /tmp"
	cmd = `go list -m -f '{{if not .Indirect}}{{.Path}}@{{.Version}}{{end}}' all`
	// cmd = `cd /Users/romber/source_code/go/src/github.com/romberli/go-util && go list -m -f '{{if not .Indirect}}{{.Path}}@{{.Version}}{{end}}' all`
	// cmd = `sh -c cd /Users/romber/source_code/go/src/github.com/romberli/go-util && go list -m -f '{{if not .Indirect}}{{.Path}}@{{.Version}}{{end}}' all`
	// cmd = "ls -l /tmp"

	// test command
	t.Log("==========test command started.==========")
	stdout, err = ExecuteCommand(cmd, WorkDirOption(workDir), UseSHCOption())
	t.Log(fmt.Sprintf("stdout: %s", stdout))
	asst.Nil(err, "test command failed.\ncmd: %s\n%v", cmd, err)

	// cmd = "ls -l /tmp/1234"
	// stdout, err = ExecuteCommand(cmd)
	// t.Log(fmt.Sprintf("stdout: %s", stdout))
	// asst.NotNil(err, "test command failed.\ncmd: %s\n%v", cmd, err)
	// t.Log("==========test command completed.==========\n")
}
