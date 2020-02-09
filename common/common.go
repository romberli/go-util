package common

import (
	"errors"
	"os"
	"reflect"

	"github.com/pkg/sftp"
)

func ConvertSliceToInterface(in interface{}) ([]interface{}, error) {
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

func ConvertMapToInterface(in interface{}) (map[interface{}]interface{}, error) {
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

func ValueInSlice(v interface{}, s interface{}) (bool, error) {
	sType := reflect.TypeOf(s)
	sValue := reflect.ValueOf(s)

	if sType.Kind() != reflect.Slice {
		return false, errors.New("second argument must be array or slice")
	}

	for i := 0; i < sValue.Len(); i++ {
		if reflect.DeepEqual(v, sValue.Index(i).Interface()) {
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

	return nil
}
