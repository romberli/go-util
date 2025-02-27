package middleware

import "fmt"

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
	return fmt.Sprintf("statement: %s, args: %v", c.Statement, c.Args)
}
