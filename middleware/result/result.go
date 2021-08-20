package result

import (
	"database/sql/driver"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
)

const mapColumnNum = 2

type Rows struct {
	FieldSlice []string
	FieldMap   map[string]int
	Values     [][]driver.Value
}

// NewRows returns *Rows
func NewRows(fieldSlice []string, fieldMap map[string]int, values [][]driver.Value) *Rows {
	return &Rows{
		fieldSlice,
		fieldMap,
		values,
	}
}

// NewRowsWithRows returns *Rows, it builds from given rows
func NewRowsWithRows(rows driver.Rows) *Rows {
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

	return &Rows{
		columns,
		fieldMap,
		values,
	}
}

// NewEmptyRows returns an empty *Rows
func NewEmptyRows() *Rows {
	return &Rows{}
}

// RowNumber returns how many rows in the result
func (r *Rows) RowNumber() int {
	return len(r.Values)
}

// ColumnNumber returns how many columns in the result
func (r *Rows) ColumnNumber() int {
	return len(r.FieldSlice)
}

// GetValue returns interface{} type value of given row and column number
func (r *Rows) GetValue(row, column int) (interface{}, error) {
	if row >= len(r.Values) || row < constant.ZeroInt {
		return nil, errors.Errorf("invalid row index %d", row)
	}

	if column >= len(r.FieldSlice) || column < constant.ZeroInt {
		return nil, errors.Errorf("invalid column index %d", column)
	}

	return r.Values[row][column], nil
}

// ColumnExists check if column exists in the result
func (r *Rows) ColumnExists(name string) bool {
	_, ok := r.FieldMap[name]

	return ok
}

// NameIndex returns index of given column
func (r *Rows) NameIndex(name string) (int, error) {
	column, ok := r.FieldMap[name]
	if ok {
		return column, nil
	}

	return constant.ZeroInt, errors.Errorf("invalid field name %s", name)
}

// GetValueByName returns interface{} type value of given row number and column name
func (r *Rows) GetValueByName(row int, name string) (interface{}, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return nil, err
	}

	return r.GetValue(row, column)
}

// IsNull checks if value of given row and column index is nil
func (r *Rows) IsNull(row, column int) (bool, error) {
	value, err := r.GetValue(row, column)
	if err != nil {
		return false, err
	}

	return value == nil, nil
}

// IsNullByName checks if value of given row number and column name is nil
func (r *Rows) IsNullByName(row int, name string) (bool, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return false, err
	}

	return r.IsNull(row, column)
}

// GetUint returns uint type value of given row and column number
func (r *Rows) GetUint(row, column int) (uint, error) {
	value, err := r.GetValue(row, column)
	if err != nil {
		return constant.ZeroInt, err
	}

	return common.ConvertToUint(value)
}

// GetUintByName returns uint type value of given row number and column name
func (r *Rows) GetUintByName(row int, name string) (uint, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return constant.ZeroInt, err
	}

	return r.GetUint(row, column)
}

// GetInt returns int type value of given row and column number
func (r *Rows) GetInt(row, column int) (int, error) {
	value, err := r.GetValue(row, column)
	if err != nil {
		return constant.ZeroInt, err
	}

	return common.ConvertToInt(value)
}

// GetIntByName returns int type value of given row number and column name
func (r *Rows) GetIntByName(row int, name string) (int, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return constant.ZeroInt, err
	}

	return r.GetInt(row, column)
}

// GetFloat returns float64 type value of given row and column number
func (r *Rows) GetFloat(row, column int) (float64, error) {
	value, err := r.GetValue(row, column)
	if err != nil {
		return constant.ZeroInt, err
	}

	return common.ConvertToFloat(value)
}

// GetFloatByName returns float64 type value of given row number and column name
func (r *Rows) GetFloatByName(row int, name string) (float64, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return constant.ZeroInt, err
	}

	return r.GetFloat(row, column)
}

// GetString returns string type value of given row and column number
func (r *Rows) GetString(row, column int) (string, error) {
	value, err := r.GetValue(row, column)
	if err != nil {
		return constant.EmptyString, err
	}

	return common.ConvertToString(value)
}

// GetStringByName returns string type value of given row number and column name
func (r *Rows) GetStringByName(row int, name string) (string, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return constant.EmptyString, err
	}

	return r.GetString(row, column)
}

// GetSlice returns []interface type value of given row and column number
func (r *Rows) GetSlice(row, column int) ([]interface{}, error) {
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

// GetSliceByName returns []interface type value of given row number and column name
func (r *Rows) GetSliceByName(row int, name string) ([]interface{}, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return nil, err
	}

	return r.GetSlice(row, column)
}

// GetUintSlice returns []uint type value of given row and column number
func (r *Rows) GetUintSlice(row, column int) ([]uint, error) {
	value, err := r.GetValue(row, column)
	if err != nil {
		return nil, err
	}

	res, err := common.ConvertToSlice(value, reflect.Uint)
	if err != nil {
		return nil, err
	}

	return res.([]uint), nil
}

// GetUintSliceByName returns []uint64 type value of given row number and column name
func (r *Rows) GetUintSliceByName(row int, name string) ([]uint, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return nil, err
	}

	return r.GetUintSlice(row, column)
}

// GetIntSlice returns []uint64 type value of given row and column number
func (r *Rows) GetIntSlice(row, column int) ([]int, error) {
	value, err := r.GetValue(row, column)
	if err != nil {
		return nil, err
	}

	res, err := common.ConvertToSlice(value, reflect.Int)
	if err != nil {
		return nil, err
	}

	return res.([]int), nil
}

// GetIntSliceByName returns []int64 type value of given row number and column name
func (r *Rows) GetIntSliceByName(row int, name string) ([]int, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return nil, err
	}

	return r.GetIntSlice(row, column)
}

// GetFloatSlice returns []float64 type value of given row and column number
func (r *Rows) GetFloatSlice(row, column int) ([]float64, error) {
	value, err := r.GetValue(row, column)
	if err != nil {
		return nil, err
	}

	res, err := common.ConvertToSlice(value, reflect.Float64)
	if err != nil {
		return nil, err
	}

	return res.([]float64), nil
}

// GetFloatSliceByName returns []float64 type value of given row number and column name
func (r *Rows) GetFloatSliceByName(row int, name string) ([]float64, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return nil, err
	}

	return r.GetFloatSlice(row, column)
}

// GetStringSlice returns []string type value of given row and column number
func (r *Rows) GetStringSlice(row, column int) ([]string, error) {
	value, err := r.GetValue(row, column)
	if err != nil {
		return nil, err
	}

	res, err := common.ConvertToSlice(value, reflect.String)
	if err != nil {
		return nil, err
	}

	return res.([]string), nil
}

// GetStringSliceByName returns []string type value of given row number and column name
func (r *Rows) GetStringSliceByName(row int, name string) ([]string, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return nil, err
	}

	return r.GetStringSlice(row, column)
}

// MapToStructSlice maps each row to a struct of the first argument,
// first argument must be a slice of pointers to structs,
// each row in the result maps to a struct in the slice,
// each column in the row maps to a field of the struct,
// tag argument is the tag of the field, it represents the column name,
// if there is no such tag in the field, this field will be ignored,
// so set tag to each field that need to be mapped,
// using "middleware" as the tag is recommended.
func (r *Rows) MapToStructSlice(in interface{}, tag string) error {
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
func (r *Rows) MapToStructByRowIndex(in interface{}, row int, tag string) error {
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
func (r *Rows) mapToStructByRowIndex(in interface{}, row int, tag string) error {
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

// MapToMapStringInterface maps rows to map[string]interface{}
func (r *Rows) MapToMapStringInterface() (map[string]interface{}, error) {
	if len(r.FieldSlice) != mapColumnNum {
		return nil, errors.New(fmt.Sprintf("to use this function, number of field must be %d, %d is not valid", mapColumnNum, len(r.FieldSlice)))
	}

	dataMap := make(map[string]interface{}, r.RowNumber())

	for i := constant.ZeroInt; i < r.RowNumber(); i++ {
		for _, field := range r.FieldSlice {
			value, err := r.GetValueByName(i, field)
			if err != nil {
				return nil, err
			}

			dataMap[field] = value
		}
	}

	return dataMap, nil
}
