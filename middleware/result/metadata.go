package result

import (
	"errors"
	"fmt"

	"github.com/romberli/go-util/constant"
)

type Metadata interface {
	// LastInsertID returns the database's auto-generated ID
	// after, for example, an INSERT into a table with primary key.
	LastInsertID() (int, error)
	// RowsAffected returns the number of rows affected by the query.
	RowsAffected() (int, error)
}

var _ Metadata = (*EmptyMetadata)(nil)

type EmptyMetadata struct {
	MiddlewareType string
}

// NewEmptyMetadata returns *EmptyMetadata with given middleware type
func NewEmptyMetadata(middlewareType string) *EmptyMetadata {
	return &EmptyMetadata{middlewareType}
}

// LastInsertID always returns error, because middleware does not support Metadata,
// this function is only for implementing the middleware.Result interface
func (em *EmptyMetadata) LastInsertID() (int, error) {
	return constant.ZeroInt, errors.New(fmt.Sprintf("LastInsertID() for %s is not supported, never call this function", em.MiddlewareType))
}

// RowsAffected always returns error, because middleware does not support Metadata,
// this function is only for implementing the middleware.Result interface
func (em *EmptyMetadata) RowsAffected() (int, error) {
	return constant.ZeroInt, errors.New(fmt.Sprintf("RowsAffected() for %s is not supported, never call this function", em.MiddlewareType))
}
