package clickhouse

import (
	"database/sql/driver"
	"fmt"
	"reflect"

	"github.com/pingcap/errors"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware"
)

var _ middleware.Result = (*Result)(nil)

type Result struct {
	FieldSlice []string
	FieldMap   map[string]int
	Values     [][]driver.Value
}

// NewResult returns *Result, it builds from given rows
func NewResult(rows driver.Rows) *Result {
	var values [][]driver.Value

	columns := rows.Columns()

	fieldMap := make(map[string]int)
	for i, column := range columns {
		fieldMap[column] = i
	}

	row := make([]driver.Value, len(columns))
	for rows.Next(row) == nil {
		r := make([]driver.Value, len(columns))
		// copy to a new slice, therefore if row is changed at next loop,
		// returning result will not be impact
		_ = copy(r, row)
		values = append(values, r)
	}

	return &Result{
		columns,
		fieldMap,
		values,
	}
}

func NewEmptyResult() *Result {
	return &Result{}
}

func (r *Result) LastInsertID() (uint64, error) {
	return constant.ZeroInt, errors.New("LastInsertID() is not supported")
}

func (r *Result) RowsAffected() (uint64, error) {
	return constant.ZeroInt, errors.New("RowsAffected() is not supported")
}

func (r *Result) RowNumber() int {
	return len(r.Values)
}

func (r *Result) ColumnNumber() int {
	return len(r.FieldSlice)
}

func (r *Result) GetValue(row, column int) (interface{}, error) {
	if row >= len(r.Values) || row < constant.ZeroInt {
		return nil, errors.Errorf("invalid row index %d", row)
	}

	if column >= len(r.FieldSlice) || column < constant.ZeroInt {
		return nil, errors.Errorf("invalid column index %d", column)
	}

	return r.Values[row][column], nil
}

func (r *Result) ColumnExists(name string) bool {
	_, ok := r.FieldMap[name]

	return ok
}

func (r *Result) NameIndex(name string) (int, error) {
	column, ok := r.FieldMap[name]
	if ok {
		return column, nil
	}

	return constant.ZeroInt, errors.Errorf("invalid field name %s", name)
}

func (r *Result) GetValueByName(row int, name string) (interface{}, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return r.GetValue(row, column)
}

func (r *Result) IsNull(row, column int) (bool, error) {
	value, err := r.GetValue(row, column)
	if err != nil {
		return false, err
	}

	return value == nil, nil
}

func (r *Result) IsNullByName(row int, name string) (bool, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return false, err
	}

	return r.IsNull(row, column)
}

func (r *Result) GetUint(row, column int) (uint64, error) {
	value, err := r.GetValue(row, column)
	if err != nil {
		return constant.ZeroInt, err
	}

	return common.ConvertToUint(value)
}

func (r *Result) GetUintByName(row int, name string) (uint64, error) {
	if column, err := r.NameIndex(name); err != nil {
		return constant.ZeroInt, err
	} else {
		return r.GetUint(row, column)
	}
}

func (r *Result) GetInt(row, column int) (int64, error) {
	value, err := r.GetUint(row, column)
	if err != nil {
		return constant.ZeroInt, err
	}

	return int64(value), nil
}

func (r *Result) GetIntByName(row int, name string) (int64, error) {
	value, err := r.GetUintByName(row, name)
	if err != nil {
		return constant.ZeroInt, err
	}

	return int64(value), nil
}

func (r *Result) GetFloat(row, column int) (float64, error) {
	value, err := r.GetValue(row, column)
	if err != nil {
		return constant.ZeroInt, err
	}

	return common.ConvertToFloat(value)
}

func (r *Result) GetFloatByName(row int, name string) (float64, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return constant.ZeroInt, err
	}

	return r.GetFloat(row, column)
}

func (r *Result) GetString(row, column int) (string, error) {
	value, err := r.GetValue(row, column)
	if err != nil {
		return constant.EmptyString, err
	}

	return common.ConvertToString(value)
}

func (r *Result) GetStringByName(row int, name string) (string, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return constant.EmptyString, err
	}

	return r.GetString(row, column)
}

func (r *Result) GetSlice(row, column int) ([]interface{}, error) {
	value, err := r.GetValue(row, column)
	if err != nil {
		return nil, err
	}

	valueKind := reflect.TypeOf(value).Kind()
	if valueKind != reflect.Slice {
		return nil, errors.New(fmt.Sprintf("value must be a slice, not %s", valueKind.String()))
	}

	valueOf := reflect.ValueOf(value)
	v := make([]interface{}, valueOf.Len())

	for i := 0; i < valueOf.Len(); i++ {
		v[i] = valueOf.Index(i).Interface()
	}

	return v, nil
}

func (r *Result) GetSliceByName(row int, name string) ([]interface{}, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return nil, err
	}

	return r.GetSlice(row, column)
}

func (r *Result) GetUintSlice(row, column int) ([]uint64, error) {
	value, err := r.GetValue(row, column)
	if err != nil {
		return nil, err
	}

	result, err := common.ConvertToSlice(value, reflect.Uint)
	if err != nil {
		return nil, err
	}

	return result.([]uint64), nil
}

func (r *Result) GetUintSliceByName(row int, name string) ([]uint64, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return nil, err
	}

	return r.GetUintSlice(row, column)
}

func (r *Result) GetIntSlice(row, column int) ([]int64, error) {
	value, err := r.GetValue(row, column)
	if err != nil {
		return nil, err
	}

	result, err := common.ConvertToSlice(value, reflect.Int)
	if err != nil {
		return nil, err
	}

	return result.([]int64), nil
}

func (r *Result) GetIntSliceByName(row int, name string) ([]int64, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return nil, err
	}

	return r.GetIntSlice(row, column)
}

func (r *Result) GetFloatSlice(row, column int) ([]float64, error) {
	value, err := r.GetValue(row, column)
	if err != nil {
		return nil, err
	}

	result, err := common.ConvertToSlice(value, reflect.Float64)
	if err != nil {
		return nil, err
	}

	return result.([]float64), nil
}

func (r *Result) GetFloatSliceByName(row int, name string) ([]float64, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return nil, err
	}

	return r.GetFloatSlice(row, column)
}

func (r *Result) GetStringSlice(row, column int) ([]string, error) {
	value, err := r.GetValue(row, column)
	if err != nil {
		return nil, err
	}

	result, err := common.ConvertToSlice(value, reflect.String)
	if err != nil {
		return nil, err
	}

	return result.([]string), nil
}

func (r *Result) GetStringSliceByName(row int, name string) ([]string, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return nil, err
	}

	return r.GetStringSlice(row, column)
}

func (r *Result) GetMap(row, column int) (map[string]interface{}, error) {
	return nil, errors.New("map data type is not supported, never call this function")
}

func (r *Result) GetMapUint(row, column int) (map[string]uint64, error) {
	return nil, errors.New("map data type is not supported, never call this function")
}

func (r *Result) GetMapInt(row, column int) (map[string]int64, error) {
	return nil, errors.New("map data type is not supported, never call this function")
}

func (r *Result) GetMapFloat(row, column int) (map[string]float64, error) {
	return nil, errors.New("map data type is not supported, never call this function")
}

func (r *Result) GetMapString(row, column int) (map[string]string, error) {
	return nil, errors.New("map data type is not supported, never call this function")
}

// MapToStructSlice maps each row to a struct of the first argument,
// first argument must be a slice of pointers to structs,
// each row in the result maps to a struct in the slice,
// each column in the row maps to a field of the struct,
// tag argument is the tag of the field, it represents the column name,
// if there is no such tag in the field, this field will be ignored,
// so set tag to each field that need to be mapped,
// using "middleware" as the tag is recommended.
func (r *Result) MapToStructSlice(in interface{}, tag string) error {
	if reflect.TypeOf(in).Kind() != reflect.Slice {
		return errors.New("first argument must be a slice of pointers to struct")
	}
	if tag == constant.EmptyString {
		return errors.New("tag argument could not be empty")
	}

	inVal := reflect.ValueOf(in)
	rowNum := r.RowNumber()
	length := inVal.Len()
	if rowNum != length {
		return errors.New(fmt.Sprintf("number of rows(%d) is not equal to length of the slice(%d)", rowNum, length))
	}

	for i := 0; i < length; i++ {
		value := inVal.Index(i).Interface()
		err := r.mapToStructByRowIndex(value, i, tag)
		if err != nil {
			return err
		}
	}

	return nil
}

// MapToStructByRowIndex maps row of given index result to the struct
// first argument must be a pointer to struct,
// each column in the row maps to a field of the struct,
// tag argument is the tag of the field, it represents the column name,
// if there is no such tag in the field, this field will be ignored,
// so set tag to each field that need to be mapped,
// using "middleware" as the tag is recommended.
func (r *Result) MapToStructByRowIndex(in interface{}, row int, tag string) error {
	if tag == constant.EmptyString {
		return errors.New("tag argument could not be empty")
	}

	return r.mapToStructByRowIndex(in, row, tag)
}

// mapToStructByRowIndex maps row of given index result to the struct
// first argument must be a pointer to struct,
// each column in the row maps to a field of the struct,
// tag argument is the tag of the field, it represents the column name,
// if there is no such tag in the field, this field will be ignored,
// so set tag to each field that need to be mapped,
// using "middleware" as the tag is recommended.
func (r *Result) mapToStructByRowIndex(in interface{}, row int, tag string) error {
	if reflect.TypeOf(in).Kind() != reflect.Ptr {
		return errors.New("first argument must be a pointer to struct")
	}

	inVal := reflect.ValueOf(in).Elem()
	inType := inVal.Type()

	for i := 0; i < inVal.NumField(); i++ {
		fieldType := inType.Field(i)
		fieldName := fieldType.Name
		columnName := fieldType.Tag.Get(tag)
		if columnName == constant.EmptyString {
			// no such tag, ignore this field
			continue
		}

		// get value with row number and column name
		fieldKind := fieldType.Type.Kind()
		switch fieldKind {
		case reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64,
			reflect.String,
			reflect.Struct,
			reflect.Slice:
			value, err := r.GetValueByName(row, columnName)
			if err != nil {
				return err
			}

			err = common.SetValueOfStructByKind(in, fieldName, value, fieldKind)
			if err != nil {
				return err
			}
		default:
			return errors.New(fmt.Sprintf("got unsupported reflect.Kind of data type: %s", fieldKind.String()))
		}
	}

	return nil
}
