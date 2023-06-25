package middleware

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pingcap/errors"
	"github.com/siddontang/go/hack"

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
func ConvertSliceToString(args ...interface{}) (string, error) {
	var result string

	if len(args) == constant.ZeroInt {
		return constant.EmptyString, errors.New("args should not be empty")
	}

	for _, arg := range args {
		argStr, err := ConvertToString(arg)
		if err != nil {
			return constant.EmptyString, err
		}
		result += argStr + constant.CommaString
	}

	return strings.Trim(result, constant.CommaString), nil
}

// GetQuestionMarks returns question marks with given num
func GetQuestionMarks(num int) string {
	var result string

	for i := constant.ZeroInt; i < num; i++ {
		result += constant.QuestionMarkString + constant.SpaceString
	}

	return strings.Trim(result, constant.QuestionMarkString+constant.SpaceString)
}

// ConvertToString converts an interface type argument to string,
// it's usually used to generate "in clause" of a select statement
func ConvertToString(arg interface{}) (string, error) {
	switch v := arg.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case []byte:
		return fmt.Sprintf("'%s'", hack.String(v)), nil
	case string:
		return fmt.Sprintf("'%s'", v), nil
	case time.Time:
		return fmt.Sprintf("'%s'", v.Format(constant.DefaultTimeLayout)), nil
	default:
		return constant.EmptyString, errors.Errorf("unsupported data type: %T", v)
	}
}
