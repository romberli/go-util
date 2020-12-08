package config

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
)

func WriteToBuffer(in interface{}, buffer *bytes.Buffer, tagType ...string) (err error) {
	var (
		tagTypeStr string
		tagName    string
		fieldStr   string
		line       string
	)

	// check if v is a struct
	inType := reflect.TypeOf(in)
	inVal := reflect.ValueOf(in)
	if inType.Kind() == reflect.Ptr {
		inVal = inVal.Elem()
		inType = inVal.Type()
	} else {
		return errors.New("can NOT parse non-pointer struct")
	}

	// check if tagType is valid
	optsLen := len(tagType)
	switch optsLen {
	case 0:
		tagTypeStr = constant.EmptyString
	case 1:
		tagTypeStr = tagType[0]
	default:
		return errors.New(fmt.Sprintf(
			"tagType should be either empty or only 1 value. actual tagType length: %d", len(tagType)))
	}

	// loop each member of the struct to get a big string
	for i := 0; i < inVal.NumField(); i++ {
		field := inVal.Field(i)
		fieldType := reflect.TypeOf(field)

		if fieldType.Kind() == reflect.Ptr {
			// this filed is a struct, we need to call recursively
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
			fieldStr, err = common.ConvertNumberToString(fieldInterface)
			if err != nil {
				return err
			}

			line = fmt.Sprintf("%s", tagName)
			if fieldStr != constant.DefaultRandomString && fieldStr != strconv.Itoa(constant.DefaultRandomInt) {
				line = fmt.Sprintf("%s = %s", tagName, fieldStr)
			}
			line += constant.CRLFString
			_, err = buffer.WriteString(line)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ConvertToString(in interface{}, tagType ...string) (s string, err error) {
	var buffer bytes.Buffer

	err = WriteToBuffer(in, &buffer, tagType...)
	if err != nil {
		return constant.EmptyString, err
	}

	return buffer.String(), nil
}

func ConvertToStringWithTitle(in interface{}, title string, tagType ...string) (s string, err error) {
	s, err = ConvertToString(in, tagType...)
	if err != nil {
		return constant.EmptyString, err
	}

	return title + constant.CRLFString + s, nil
}
