package middleware

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"

	"github.com/romberli/go-util/constant"
)

// ConvertArgsToNamedValues converts args to named values
func ConvertArgsToNamedValues(args ...interface{}) []driver.NamedValue {
	namedValues := make([]driver.NamedValue, len(args))

	for i, arg := range args {
		namedValues[i] = driver.NamedValue{
			Name:    constant.EmptyString,
			Ordinal: i + 1,
			Value:   driver.Value(arg),
		}
	}

	return namedValues
}

// ConvertSliceToString converts args to string,
// it's usually used to generate "in clause" of a select statement
func ConvertSliceToString(args []interface{}) (string, error) {
	var result string

	if len(args) == constant.ZeroInt {
		return constant.EmptyString, errors.New("args should not be empty")
	}

	first := args[0]
	result, err := ConvertToString(first)
	if err != nil {
		return constant.EmptyString, err
	}

	for i := 1; i < len(args); i++ {
		argStr, err := ConvertToString(args[i])
		if err != nil {
			return constant.EmptyString, err
		}
		result += fmt.Sprintf(", %s", argStr)
	}

	return result, nil
}

// ConvertToString converts an interface type argument to string,
// it's usually used to generate "in clause" of a select statement
func ConvertToString(arg interface{}) (string, error) {
	switch arg.(type) {
	case string:
		return fmt.Sprintf("'%s'", arg), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", arg), nil
	case float32, float64:
		return fmt.Sprintf("%f", arg), nil
	default:
		return constant.EmptyString, errors.New(fmt.Sprintf("only support string, integer, float type convertion, %s is not valid", reflect.TypeOf(arg).String()))
	}
}
