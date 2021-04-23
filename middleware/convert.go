package middleware

import (
	"database/sql/driver"

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
