package prometheus

import (
	"database/sql/driver"

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

type RawData map[string]interface{}

func NewRawData(value model.Value, warnings apiv1.Warnings) RawData {
	return map[string]interface{}{valueColumn: value, warningsColumn: warnings}
}

func (rd RawData) GetValue() model.Value {
	return rd[valueColumn].(model.Value)
}

func (rd RawData) GetWarnings() apiv1.Warnings {
	return rd[warningsColumn].(apiv1.Warnings)
}

type Result struct {
	Raw RawData
	*result.Rows
	result.Metadata
	result.Map
}

func NewResult(value model.Value, warnings apiv1.Warnings) *Result {
	var values [][]driver.Value

	fieldSlice := []string{valueColumn, timestampColumn}
	fieldMap := map[string]int{valueColumn: 0, timestampColumn: 1}
	row := make([]driver.Value, defaultColumnNum)

	switch v := value.(type) {
	case *model.Scalar:
		row[0] = float64(v.Value)
		row[1] = v.Timestamp.Time()
		values = append(values, row)
	case *model.String:
		row[0] = v.Value
		row[1] = v.Timestamp.Time()
		values = append(values, row)
	case model.Vector:
		for i := 0; i < len(v); i++ {
			sample := v[i]
			row[0] = float64(sample.Value)
			row[1] = sample.Timestamp.Time()
			values = append(values, row)
		}
	case model.Matrix:
		samplePairs := v[constant.ZeroInt].Values
		for i := 0; i < len(samplePairs); i++ {
			sp := samplePairs[i]
			row[0] = float64(sp.Value)
			row[1] = sp.Timestamp.Time()
			values = append(values, row)
		}
	}

	return &Result{
		Raw:      NewRawData(value, warnings),
		Rows:     result.NewRows(fieldSlice, fieldMap, values),
		Metadata: result.NewEmptyMetadata(middlewareType),
		Map:      result.NewEmptyMap(middlewareType),
	}
}

func NewEmptyResult() *Result {
	return &Result{}
}
