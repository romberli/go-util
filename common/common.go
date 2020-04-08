package common

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/pkg/sftp"
)

func ConvertInterfaceToSliceInterface(in interface{}) ([]interface{}, error) {
	inType := reflect.TypeOf(in)
	inValue := reflect.ValueOf(in)

	if inType.Kind() != reflect.Slice {
		return nil, errors.New("argument must be array or slice")
	}

	inLength := inValue.Len()
	sliceInterface := make([]interface{}, inLength)

	for i := 0; i < inLength; i++ {
		sliceInterface[i] = inValue.Index(i).Interface()
	}

	return sliceInterface, nil
}

func ConvertInterfaceToMapInterfaceInterface(in interface{}) (map[interface{}]interface{}, error) {
	inType := reflect.TypeOf(in)
	inValue := reflect.ValueOf(in)

	if inType.Kind() != reflect.Map {
		return nil, errors.New("argument must be map")
	}

	inLength := inValue.Len()
	mapInterface := make(map[interface{}]interface{}, inLength)

	for _, key := range inValue.MapKeys() {
		mapInterface[key.Interface()] = inValue.MapIndex(key).Interface()
	}

	return mapInterface, nil
}

func ElementInSlice(e interface{}, s interface{}) (bool, error) {
	sType := reflect.TypeOf(s)
	sValue := reflect.ValueOf(s)

	if sType.Kind() != reflect.Slice {
		return false, errors.New("second argument must be array or slice")
	}

	for i := 0; i < sValue.Len(); i++ {
		if reflect.DeepEqual(e, sValue.Index(i).Interface()) {
			return true, nil
		}
	}

	return false, nil
}

func KeyInMap(k interface{}, m interface{}) (bool, error) {
	if reflect.TypeOf(m).Kind() != reflect.Map {
		return false, errors.New("second argument must be map")
	}

	iter := reflect.ValueOf(m).MapRange()
	for iter.Next() {
		if reflect.DeepEqual(k, iter.Key().Interface()) {
			return true, nil
		}
	}

	return false, nil
}

func ValueInMap(v interface{}, m interface{}) (bool, error) {
	if reflect.TypeOf(m).Kind() != reflect.Map {
		return false, errors.New("second argument must be map")
	}

	iter := reflect.ValueOf(m).MapRange()
	for iter.Next() {
		if reflect.DeepEqual(v, iter.Value().Interface()) {
			return true, nil
		}
	}

	return false, nil
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

func PathExists(in ...interface{}) (bool, error) {
	if len(in) == 0 {
		return false, errors.New("argument could not be nil")
	} else if len(in) == 1 {
		path := in[0]

		switch path.(type) {
		case string:
			return PathExistsLocal(path.(string))
		default:
			return false, errors.New("first argument must be string type, which presents a file or directory")
		}
	} else {
		path := in[0]
		client := in[1]

		switch path.(type) {
		case string:
		default:
			return false, errors.New("first argument must be string type, which presents a file or directory")
		}

		switch client.(type) {
		case nil:
			return false, errors.New("second argument could not be nil")
		case *sftp.Client:
			return PathExistsRemote(path.(string), client.(*sftp.Client))
		default:
			return false, errors.New(
				fmt.Sprintf("second argument must be *sftp.Client type instead of %s",
					reflect.TypeOf(client).Name()))
		}
	}
}

func TrimSpaceOfStructString(in interface{}) error {
	inType := reflect.TypeOf(in)
	inVal := reflect.ValueOf(in)

	if inType.Kind() == reflect.Ptr {
		inVal = inVal.Elem()
		inType = inVal.Type()
	} else {
		return errors.New("argument must be a pointer to struct")
	}

	for i := 0; i < inVal.NumField(); i++ {
		f := inVal.Field(i)
		switch f.Kind() {
		case reflect.String:
			if f.CanSet() {
				trimValue := strings.TrimSpace(f.String())
				f.Set(reflect.ValueOf(trimValue))
			}
		}
	}

	return nil
}
