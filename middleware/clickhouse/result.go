package clickhouse

import (
	"database/sql"

	"github.com/romberli/go-util/middleware"
	"github.com/romberli/go-util/middleware/result"
)

const middlewareType = "clickhouse"

var _ middleware.Result = (*Result)(nil)

type Result struct {
	Raw *sql.Rows
	*result.Rows
	result.Metadata
	result.Map
}

// NewResult returns *Result, it builds from given rows
func NewResult(rows *sql.Rows) (*Result, error) {
	r, err := result.NewRowsWithSQLRows(rows)
	if err != nil {
		return nil, err
	}

	return &Result{
		rows,
		r,
		result.NewEmptyMetadata(middlewareType),
		result.NewEmptyMap(middlewareType),
	}, nil
}

// NewEmptyResult returns an empty *Result
func NewEmptyResult() *Result {
	return &Result{}
}

// GetRaw returns the raw data of the result
func (r *Result) GetRaw() interface{} {
	return r.Raw
}
