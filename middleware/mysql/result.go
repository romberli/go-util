package mysql

import (
	"database/sql/driver"
	"encoding/json"
	"reflect"

	"github.com/go-mysql-org/go-mysql/mysql"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware"
	"github.com/romberli/go-util/middleware/result"
)

const middlewareType = "mysql"

var _ middleware.Result = (*Result)(nil)

type Result struct {
	Raw          *mysql.Result `json:"-"`
	*result.Rows `json:"rows"`
	result.Map   `json:"-"`
}

// NewResult returns *Result
func NewResult(r *mysql.Result) *Result {
	if r == nil {
		columns := make([]string, constant.ZeroInt)
		filedNames := make(map[string]int)
		values := make([][]driver.Value, constant.ZeroInt)

		return &Result{
			&mysql.Result{Resultset: &mysql.Resultset{}},
			result.NewRows(columns, filedNames, values),
			result.NewEmptyMap(middlewareType),
		}
	}

	rowNum := r.RowNumber()

	var colNum int
	if r.Resultset != nil {
		colNum = r.ColumnNumber()
	}

	columns := make([]string, colNum)
	values := make([][]driver.Value, rowNum)
	fieldNames := make(map[string]int)

	if r.Resultset != nil {
		if r.Resultset.FieldNames != nil {
			fieldNames = r.Resultset.FieldNames
		}

		for fieldName, fieldIndex := range r.Resultset.FieldNames {
			columns[fieldIndex] = fieldName
		}

		for i := constant.ZeroInt; i < r.RowNumber(); i++ {
			values[i] = make([]driver.Value, colNum)

			for j := constant.ZeroInt; j < r.ColumnNumber(); j++ {
				values[i][j] = r.Resultset.Values[i][j].Value()
			}
		}
	}

	return &Result{
		r,
		result.NewRows(columns, fieldNames, values),
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

// MarshalJSON custom marshal method to handle []uint8 to string conversion
func (r *Result) MarshalJSON() ([]byte, error) {
	serialized := common.SerializeBytes(r)

	return json.Marshal(serialized)
}

// UnmarshalJSON custom unmarshal method to handle []uint8 to string conversion
func (r *Result) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	deserialized := common.DeserializeBytes(raw, reflect.TypeOf(*r))
	*r = deserialized.(Result)
	return nil
}
