package prometheus

import (
	"database/sql/driver"
	"errors"
	"fmt"

	apiv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware"
	"github.com/romberli/go-util/middleware/result"
)

const (
	defaultColumnNum = 2
	middlewareType   = "prometheus"
	timestampColumn  = "timestamp"
	valueColumn      = "value"
	warningsColumn   = "warnings"
)

var _ middleware.Result = (*Result)(nil)

type RawData struct {
	Value    model.Value
	Warnings apiv1.Warnings
}

// NewRawData returns a new RawData
func NewRawData(value model.Value, warnings apiv1.Warnings) *RawData {
	return &RawData{
		Value:    value,
		Warnings: warnings,
	}
}

// GetValue returns value
func (rd *RawData) GetValue() model.Value {
	return rd.Value
}

// GetWarnings returns warnings
func (rd *RawData) GetWarnings() apiv1.Warnings {
	return rd.Warnings
}

func (rd *RawData) GetMatrix() (model.Matrix, error) {
	value, ok := rd.GetValue().(model.Matrix)
	if !ok {
		return nil, errors.New(fmt.Sprintf("type assertion failed, %T could not be converted to model.Matrix", rd.GetValue()))
	}

	return value, nil
}

func (rd *RawData) GetVector() (model.Vector, error) {
	value, ok := rd.GetValue().(model.Vector)
	if !ok {
		return nil, errors.New(fmt.Sprintf("type assertion failed, %T could not be converted to model.Vector", rd.GetValue()))
	}

	return value, nil
}

func (rd *RawData) GetScalar() (*model.Scalar, error) {
	value, ok := rd.GetValue().(*model.Scalar)
	if !ok {
		return nil, errors.New(fmt.Sprintf("type assertion failed, %T could not be converted to *model.Scalar", rd.GetValue()))
	}

	return value, nil
}

func (rd *RawData) GetString() (*model.String, error) {
	value, ok := rd.GetValue().(*model.String)
	if !ok {
		return nil, errors.New(fmt.Sprintf("type assertion failed, %T could not be converted to *model.String", rd.GetValue()))
	}

	return value, nil
}

type Result struct {
	Raw *RawData
	*result.Rows
	result.Metadata
	result.Map
}

// NewResult returns a new *Result with given value and warnings
// note that if return value is matrix type, only the first matrix will be processed,
// all others will be discarded, if a query returns more than one matrix,
// use GetRaw() function to get the raw data which is returned by prometheus go client package
func NewResult(value model.Value, warnings apiv1.Warnings) *Result {
	var values [][]driver.Value

	fieldSlice := []string{valueColumn, timestampColumn}
	fieldMap := map[string]int{valueColumn: 0, timestampColumn: 1}

	switch v := value.(type) {
	case *model.Scalar:
		if v != nil {
			values = make([][]driver.Value, 1)
			values[constant.ZeroInt] = make([]driver.Value, defaultColumnNum)

			values[constant.ZeroInt][0] = float64(v.Value)
			values[constant.ZeroInt][1] = v.Timestamp.Time()
		}
	case *model.String:
		if v != nil {
			values = make([][]driver.Value, 1)
			values[constant.ZeroInt] = make([]driver.Value, defaultColumnNum)

			values[constant.ZeroInt][0] = v.Value
			values[constant.ZeroInt][1] = v.Timestamp.Time()
		}
	case model.Vector:
		if v != nil && v.Len() > constant.ZeroInt {
			values = make([][]driver.Value, v.Len())

			for i := 0; i < v.Len(); i++ {
				values[i] = make([]driver.Value, defaultColumnNum)

				values[i][0] = float64(v[i].Value)
				values[i][1] = v[i].Timestamp.Time()
			}
		}
	case model.Matrix:
		// note that only the first matrix value will be processed,
		// if a query returns more than one matrix,
		// use GetRaw() function to get the raw data which is returned by prometheus go client package
		if v != nil && v.Len() > constant.ZeroInt {
			samplePairs := v[constant.ZeroInt].Values
			values = make([][]driver.Value, len(samplePairs))

			for i := 0; i < len(samplePairs); i++ {
				values[i] = make([]driver.Value, defaultColumnNum)

				values[i][0] = float64(samplePairs[i].Value)
				values[i][1] = samplePairs[i].Timestamp.Time()
			}
		}
	}

	return &Result{
		Raw:      NewRawData(value, warnings),
		Rows:     result.NewRows(fieldSlice, fieldMap, values),
		Metadata: result.NewEmptyMetadata(middlewareType),
		Map:      result.NewEmptyMap(middlewareType),
	}
}

// NewEmptyResult returns a new empty *Result
func NewEmptyResult() *Result {
	return &Result{}
}

// GetRaw returns the raw data of the result
func (r *Result) GetRaw() interface{} {
	return r.Raw
}
