package clickhouse

import (
	"context"

	"github.com/ClickHouse/clickhouse-go"
	"github.com/romberli/go-util/middleware"
	"github.com/romberli/go-util/middleware/sql/statement"
)

var _ middleware.Statement = (*Statement)(nil)

type Statement struct {
	clickhouse.Stmt
}

func NewStatement(stmt clickhouse.Stmt) *Statement {
	return &Statement{stmt}
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
		rows, err := stmt.QueryContext(ctx, middleware.ConvertArgsToNamedValues(args...))
		if err != nil {
			return nil, err
		}
		defer func() { _ = rows.Close() }()

		return NewResult(rows), nil
	}
	// this is not a select sql
	_, err := stmt.ExecContext(ctx, middleware.ConvertArgsToNamedValues(args...))
	if err != nil {
		return nil, err
	}

	return NewEmptyResult(), nil
}
