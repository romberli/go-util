package clickhouse

import (
	"database/sql/driver"

	"github.com/romberli/go-util/middleware/result"

	"github.com/romberli/go-util/middleware"
)

const middlewareType = "clickhouse"

var _ middleware.Result = (*Result)(nil)

type Result struct {
	Raw driver.Rows
	*result.Rows
	result.Metadata
	result.Map
}

// NewResult returns *Result, it builds from given rows
func NewResult(rows driver.Rows) *Result {
	return &Result{
		rows,
		result.NewRowsWithRows(rows),
		result.NewEmptyMetadata(middlewareType),
		result.NewEmptyMap(middlewareType),
	}
}

// NewEmptyResult returns an empty *Result
func NewEmptyResult() *Result {
	return &Result{}
}

// GetRaw returns the raw data of the result
func (r *Result) GetRaw() interface{} {
	return r.Raw
}
