package middleware

import (
	"context"

	"github.com/romberli/go-util/middleware/result"
)

type Result interface {
	result.Raw
	result.Metadata
	result.Value
	result.Bool
	result.Number
	result.String
	result.Slice
	result.Map
	result.Unmarshaler
}

type Statement interface {
	// Close closes the statement
	Close() error
	// Execute executes given command and placeholders on the middleware
	Execute(args ...interface{}) (Result, error)
	// ExecuteContext executes given command and placeholders with context on the middleware
	ExecuteContext(ctx context.Context, args ...interface{}) (Result, error)
	// Commit commits the transaction
	Commit() error
	// Rollback rollbacks the transaction
	Rollback() error
}

type PoolConn interface {
	// Close returns connection back to the pool
	Close() error
	// Disconnect disconnects from the middleware, normally when using connection pool
	Disconnect() error
	// IsValid validates if connection is valid
	IsValid() bool
	// Prepare prepares a statement and returns a Statement
	Prepare(command string) (Statement, error)
	// PrepareContext prepares a statement with context and returns a Statement
	PrepareContext(ctx context.Context, command string) (Statement, error)
	// ExecuteInBatch executes given commands and placeholders on the middleware
	ExecuteInBatch(commands []*Command, isTransaction bool) ([]Result, error)
	// Execute executes given command and placeholders on the middleware
	Execute(command string, args ...interface{}) (Result, error)
	// ExecuteContext executes given command and placeholders with context on the middleware
	ExecuteContext(ctx context.Context, command string, args ...interface{}) (Result, error)
}

type Transaction interface {
	PoolConn
	// Begin begins a transaction
	Begin() error
	// Commit commits current transaction
	Commit() error
	// Rollback rollbacks current transaction
	Rollback() error
}

type Pool interface {
	// Close releases each connection in the pool
	Close() error
	// IsClosed returns if pool had been closed
	IsClosed() bool
	// Get gets a connection from the pool
	Get() (PoolConn, error)
	// Transaction returns a connection that could run multiple statements in the same transaction
	Transaction() (Transaction, error)
	// Supply creates given number of connections and add them to the pool
	Supply(num int) error
	// Release releases given number of connections, each connection will disconnect with the middleware
	Release(num int) error
}
