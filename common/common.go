package common

import (
	"errors"
	"os"
	"reflect"

	"github.com/pkg/sftp"
)

func StringInSlice(str string, slice []string) bool {
	for i := range slice {
		if slice[i] == str {
			return true
		}
	}

	return false
}

func StringInMap(str string, m map[string]string) bool {
	if _, ok := m[str]; ok {
		return true
	}

	return false
}

func PathExistsLocal(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func PathExistsRemote(path string, client *sftp.Client) (bool, error) {
	if _, err := client.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func TrimSpaceOfStructString(in interface{}) error {
	inType := reflect.TypeOf(in)
	inVal := reflect.ValueOf(in)

	if inType.Kind() == reflect.Ptr {
		inType = inType.Elem()
		inType = inType.Elem()
	} else {
		return errors.New("argument must be a pointer to struct")
	}

	for i := 0; i < inVal.NumField(); i++ {
		f := inVal.Field(i)
		switch f.Kind() {
		case reflect.String:
			trimValue := f.String()
			f.Set(reflect.ValueOf(trimValue))
		}
	}
}
