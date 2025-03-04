package middleware

import (
	"fmt"
	"regexp"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
)

type Command struct {
	Statement string        `middleware:"statement" json:"statement"`
	Args      []interface{} `middleware:"args" json:"args"`
}

// NewCommand returns a new *Command
func NewCommand(statement string, args ...interface{}) *Command {
	return &Command{
		Statement: statement,
		Args:      args,
	}
}

// String returns the sql of the command
func (c *Command) String() string {
	return fmt.Sprintf("statement: %s, args: %v",
		c.GetMaskedStatement(), common.ConvertInterfaceSliceToString(c.Args, constant.CommaString))
}

func (c *Command) GetMaskedStatement() string {
	re := regexp.MustCompile(`(?i)(IDENTIFIED BY\s*')([^']+)(')`)

	return re.ReplaceAllString(c.Statement, "${1}xxxxxx${3}")
}
