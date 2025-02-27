package middleware

import "fmt"

type Command struct {
	statement string
	args      []interface{}
}

// NewCommand returns a new *Command
func NewCommand(statement string, args ...interface{}) *Command {
	return &Command{
		statement: statement,
		args:      args,
	}
}

// GetStatement returns the statement of the command
func (c *Command) GetStatement() string {
	return c.statement
}

// GetArgs returns the args of the command
func (c *Command) GetArgs() []interface{} {
	return c.args
}

// String returns the sql of the command
func (c *Command) String() string {
	return fmt.Sprintf("statement: %s, args: %v", c.statement, c.args)
}
