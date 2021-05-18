package mysql

import (
	"database/sql/driver"

	"github.com/romberli/go-util/middleware"
	"github.com/romberli/go-util/middleware/result"

	"github.com/go-mysql-org/go-mysql/mysql"
)

const middlewareType = "mysql"

var _ middleware.Result = (*Result)(nil)

type Result struct {
	Raw *mysql.Result
	*result.Rows
	result.Map
}

// NewResult returns *Result
func NewResult(r *mysql.Result) *Result {
	columns := make([]string, r.ColumnNumber())
	for fieldName, fieldIndex := range r.FieldNames {
		columns[fieldIndex] = fieldName
	}

	values := make([][]driver.Value, r.RowNumber())
	for i := 0; i < r.RowNumber(); i++ {
		values[i] = make([]driver.Value, r.ColumnNumber())

		for j := 0; j < r.ColumnNumber(); j++ {
			values[i][j] = r.Values[i][j].Value()
		}
	}

	return &Result{
		r,
		result.NewRows(columns, r.FieldNames, values),
		result.NewEmptyMap(middlewareType),
	}
}

// LastInsertID returns the database's auto-generated ID
// after, for example, an INSERT into a table with primary key.
func (r *Result) LastInsertID() (int, error) {
	return int(r.Raw.InsertId), nil
}

// RowsAffected returns the number of rows affected by the query.
func (r *Result) RowsAffected() (int, error) {
	return int(r.Raw.AffectedRows), nil
}

func (r *Result) GetRaw() interface{} {
	return r.Raw
}
