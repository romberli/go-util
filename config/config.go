package config

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"

	"github.com/pingcap/errors"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
)

// WriteToBuffer loops each member of given input struct recursively,
// converts member variable names and concerning values to "key = value" string,
// and then write the string into buffer, it follows some rules:
//  1. if tag type is specified, which is optional, "key" will be replaced by the tag name
//  2. toml and ini type config files requires each key should has a value(but could be empty),
//     but some toml/ini like config files(for example: my.cnf) allows variables are only keys instead of key value pairs,
//     (for example: skip-name-resolve), in this case, you can specify constant.DefaultRandomString or
//     constant.DefaultRandomInt to those keys not having values, therefore when this function find out those constant values,
//     it will convert to key string ignore the value and also the equal mark
func WriteToBuffer(in interface{}, buffer *bytes.Buffer, tagType ...string) (err error) {
	var (
		tagTypeStr string
		tagName    string
		fieldStr   string
		line       string
	)

	kind := reflect.TypeOf(in).Kind()
	if kind != reflect.Ptr {
		return errors.Errorf("second argument must be a pointer to struct, %s is not valid", kind.String())
	}

	inVal := reflect.ValueOf(in).Elem()
	inType := inVal.Type()

	// check if tagType is valid
	optsLen := len(tagType)
	switch optsLen {
	case 0:
		tagTypeStr = constant.EmptyString
	case 1:
		tagTypeStr = tagType[constant.ZeroInt]
	default:
		return errors.Errorf("tagType should be either empty or only have 1 value. actual tagType length: %d", len(tagType))
	}

	// loop each member of the struct to get a big string
	for i := constant.ZeroInt; i < inVal.NumField(); i++ {
		field := inVal.Field(i)
		fieldType := reflect.TypeOf(field)

		if fieldType.Kind() == reflect.Ptr {
			// this filed is also a struct, we need to call recursively
			err = WriteToBuffer(field, buffer, tagType...)
			if err != nil {
				return err
			}
		} else {
			// this field is a normal value
			if tagTypeStr != constant.EmptyString {
				f := inType.Field(i)
				tagName = f.Tag.Get(tagTypeStr)
			} else {
				tagName = fieldType.Name()
			}

			fieldInterface := field.Interface()
			// convert field value to string
			fieldStr, err = common.ConvertNumberToString(fieldInterface)
			if err != nil {
				return err
			}

			line = tagName
			if fieldStr != constant.DefaultRandomString && fieldStr != strconv.Itoa(constant.DefaultRandomInt) {
				// this field has a value
				line += fmt.Sprintf(" = %s", fieldStr)
			}
			line += constant.CRLFString
			_, err = buffer.WriteString(line)
			if err != nil {
				return errors.Trace(err)
			}
		}
	}

	return nil
}

// ConvertToString convert struct to string
func ConvertToString(in interface{}, tagType ...string) (s string, err error) {
	var buffer bytes.Buffer

	err = WriteToBuffer(in, &buffer, tagType...)
	if err != nil {
		return constant.EmptyString, err
	}

	return buffer.String(), nil
}

// ConvertToStringWithTitle convert struct to string with given title
func ConvertToStringWithTitle(in interface{}, title string, tagType ...string) (s string, err error) {
	s, err = ConvertToString(in, tagType...)
	if err != nil {
		return constant.EmptyString, err
	}

	return title + constant.CRLFString + s, nil
}
