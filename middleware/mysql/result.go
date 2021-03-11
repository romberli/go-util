package mysql

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware"

	"github.com/siddontang/go-mysql/mysql"
)

const middlewareType = "mysql"

var _ middleware.Result = (*Result)(nil)

type Result struct {
	middleware.SliceResult
	middleware.MapResult
	*mysql.Result
}

func NewResult(r *mysql.Result) *Result {
	return &Result{
		middleware.NewEmptySliceResult(middlewareType),
		middleware.NewEmptyMapResult(middlewareType),
		r,
	}
}

// LastInsertID returns the database's auto-generated ID
// after, for example, an INSERT into a table with primary key.
func (r *Result) LastInsertID() (uint64, error) {
	return r.InsertId, nil
}

// RowsAffected returns the number of rows affected by the query.
func (r *Result) RowsAffected() (uint64, error) {
	return r.AffectedRows, nil
}

// RowNumber returns how many rows in the result
func (r *Result) RowNumber() int {
	return r.Result.RowNumber()
}

// ColumnNumber returns how many columns in the result
func (r *Result) ColumnNumber() int {
	return r.Result.ColumnNumber()
}

// GetValue returns interface{} type value of given row and column number
func (r *Result) GetValue(row, column int) (interface{}, error) {
	return r.Result.GetValue(row, column)
}

// ColumnExists check if column exists in the result
func (r *Result) ColumnExists(name string) bool {
	_, ok := r.FieldNames[name]
	return ok
}

// NameIndex returns index of given column
func (r *Result) NameIndex(name string) (int, error) {
	return r.Result.NameIndex(name)
}

// GetValueByName returns interface{} type value of given row number and column name
func (r *Result) GetValueByName(row int, name string) (interface{}, error) {
	return r.Result.GetValueByName(row, name)
}

// IsNull checks if value of given row and column index is nil
func (r *Result) IsNull(row, column int) (bool, error) {
	return r.Result.IsNull(row, column)
}

// IsNullByName checks if value of given row number and column name is nil
func (r *Result) IsNullByName(row int, name string) (bool, error) {
	return r.Result.IsNullByName(row, name)
}

// GetUint returns uint64 type value of given row and column number
func (r *Result) GetUint(row, column int) (uint64, error) {
	return r.Result.GetUint(row, column)
}

// GetUintByName returns uint64 type value of given row number and column name
func (r *Result) GetUintByName(row int, name string) (uint64, error) {
	return r.Result.GetUintByName(row, name)
}

// GetInt returns int64 type value of given row and column number
func (r *Result) GetInt(row, column int) (int64, error) {
	return r.Result.GetInt(row, column)
}

// GetIntByName returns int64 type value of given row number and column name
func (r *Result) GetIntByName(row int, name string) (int64, error) {
	return r.Result.GetIntByName(row, name)
}

// GetFloat returns float64 type value of given row and column number
func (r *Result) GetFloat(row, column int) (float64, error) {
	return r.Result.GetFloat(row, column)
}

// GetFloatByName returns float64 type value of given row number and column name
func (r *Result) GetFloatByName(row int, name string) (float64, error) {
	return r.Result.GetFloatByName(row, name)
}

// GetString returns string type value of given row and column number
func (r *Result) GetString(row, column int) (string, error) {
	return r.Result.GetString(row, column)
}

// GetStringByName returns string type value of given row number and column name
func (r *Result) GetStringByName(row int, name string) (string, error) {
	return r.Result.GetStringByName(row, name)
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
		case reflect.Bool:
			intVal, err := r.GetIntByName(row, columnName)
			if err != nil {
				return err
			}
			switch intVal {
			case 0:
				err = common.SetValueOfStruct(in, fieldName, false)
			case 1:
				err = common.SetValueOfStruct(in, fieldName, true)
			default:
				err = errors.New(fmt.Sprintf("bool type value should be either 0 or 1, %d is not valid", intVal))
			}
			if err != nil {
				return err
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			intVal, err := r.GetIntByName(row, columnName)
			if err != nil {
				return err
			}
			switch fieldKind {
			case reflect.Int:
				err = common.SetValueOfStruct(in, fieldName, int(intVal))
			case reflect.Int8:
				err = common.SetValueOfStruct(in, fieldName, int8(intVal))
			case reflect.Int16:
				err = common.SetValueOfStruct(in, fieldName, int16(intVal))
			case reflect.Int32:
				err = common.SetValueOfStruct(in, fieldName, int32(intVal))
			case reflect.Int64:
				err = common.SetValueOfStruct(in, fieldName, intVal)
			case reflect.Uint:
				err = common.SetValueOfStruct(in, fieldName, uint(intVal))
			case reflect.Uint8:
				err = common.SetValueOfStruct(in, fieldName, uint8(intVal))
			case reflect.Uint16:
				err = common.SetValueOfStruct(in, fieldName, uint16(intVal))
			case reflect.Uint32:
				err = common.SetValueOfStruct(in, fieldName, uint32(intVal))
			case reflect.Uint64:
				err = common.SetValueOfStruct(in, fieldName, uint64(intVal))
			}
			if err != nil {
				return err
			}
		case reflect.Float32, reflect.Float64:
			floatVal, err := r.GetFloatByName(row, columnName)
			if err != nil {
				return err
			}
			switch fieldKind {
			case reflect.Float32:
				err = common.SetValueOfStruct(in, fieldName, float32(floatVal))
			case reflect.Float64:
				err = common.SetValueOfStruct(in, fieldName, floatVal)
			}
			if err != nil {
				return err
			}
		case reflect.String:
			stringVal, err := r.GetStringByName(row, columnName)
			if err != nil {
				return err
			}
			err = common.SetValueOfStruct(in, fieldName, stringVal)
			if err != nil {
				return err
			}
		case reflect.Struct:
			stringVal, err := r.GetStringByName(row, columnName)
			if err != nil {
				return err
			}
			// for now, only support time.Time data type,
			// so if data type of field of struct is not time.Time,
			// it will return error
			t, err := time.ParseInLocation(constant.DefaultTimeLayout, stringVal, time.Local)
			if err != nil {
				return err
			}
			err = common.SetValueOfStruct(in, fieldName, t)
			if err != nil {
				return err
			}
		default:
			return errors.New(fmt.Sprintf("got unsupported reflect.Kind of data type: %s", fieldKind.String()))
		}
	}

	return nil
}
