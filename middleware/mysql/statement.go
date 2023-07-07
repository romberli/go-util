package mysql

import (
	"context"

	"github.com/go-mysql-org/go-mysql/client"

	"github.com/romberli/go-util/middleware"
)

const (
	commitSQL   = "COMMIT"
	rollbackSQL = "ROLLBACK"
)

var _ middleware.Statement = (*Statement)(nil)

type Statement struct {
	*client.Stmt
}

// NewStatement returns a new *Statement with given *client.Stmt
func NewStatement(stmt *client.Stmt) *Statement {
	return &Statement{stmt}
}

// Close closes the statement
func (stmt *Statement) Close() error {
	return stmt.Stmt.Close()
}

// Execute executes given sql and placeholders and returns a result
func (stmt *Statement) Execute(args ...interface{}) (middleware.Result, error) {
	return stmt.executeContext(context.Background(), args...)
}

// ExecuteContext executes given sql and placeholders with context and returns a result
func (stmt *Statement) ExecuteContext(ctx context.Context, args ...interface{}) (middleware.Result, error) {
	return stmt.executeContext(ctx, args...)
}

// executeContext executes given sql and placeholders with context and returns a result
func (stmt *Statement) executeContext(ctx context.Context, args ...interface{}) (*Result, error) {
	r, err := stmt.Stmt.Execute(args...)
	if err != nil {
		return nil, err
	}

	return NewResult(r), nil
}

// Commit commits the transaction
func (stmt *Statement) Commit() error {
	_, err := stmt.Stmt.Execute(commitSQL)

	return err
}

// Rollback rollbacks the transaction
func (stmt *Statement) Rollback() error {
	_, err := stmt.Execute(rollbackSQL)

	return err
}
