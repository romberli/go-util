package clickhouse

import (
	"context"
	"database/sql"

	"github.com/pingcap/errors"
	"github.com/romberli/go-util/middleware"
	"github.com/romberli/go-util/middleware/sql/statement"
)

var _ middleware.Statement = (*Statement)(nil)

type Statement struct {
	stmt  *sql.Stmt
	tx    *sql.Tx
	query string
}

func NewStatement(stmt *sql.Stmt, tx *sql.Tx, query string) middleware.Statement {
	return newStatement(stmt, tx, query)
}

func newStatement(stmt *sql.Stmt, tx *sql.Tx, query string) *Statement {
	return &Statement{
		stmt:  stmt,
		tx:    tx,
		query: query,
	}
}

// GetQuery returns the sql of the statement
func (stmt *Statement) GetQuery() string {
	return stmt.query
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
	// get sql type
	sqlType := statement.GetType(stmt.GetQuery())
	if sqlType == statement.Select {
		// this is a select sql
		rows, err := stmt.stmt.QueryContext(ctx, args...)
		if err != nil {
			return nil, errors.Trace(err)
		}
		defer func() { _ = rows.Close() }()

		return NewResult(rows)
	}
	// this is not a select sql
	_, err := stmt.stmt.ExecContext(ctx, args...)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return NewEmptyResult(), nil
}

// Commit commits the transaction
func (stmt *Statement) Commit() error {
	return stmt.tx.Commit()
}

// Rollback rollbacks the transaction
func (stmt *Statement) Rollback() error {
	return stmt.tx.Rollback()
}

// Close closes the statement
func (stmt *Statement) Close() error {
	return stmt.stmt.Close()
}
