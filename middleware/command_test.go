package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand_String(t *testing.T) {
	asst := assert.New(t)

	c := NewCommand("select * from table", 1, "test")
	asst.Equal("statement: select * from table, args: 1,test", c.String(), "TestCommand_String() failed")
	c = NewCommand("create user 'test' identified by 'test'")
	asst.Equal("statement: create user 'test' identified by 'xxxxxx', args: ", c.String(), "TestCommand_String() failed")
}
