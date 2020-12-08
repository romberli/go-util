package common

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/romberli/go-util/constant"
)

// StringToBytes converts string type to byte slice
func StringToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}

	return *(*[]byte)(unsafe.Pointer(&h))
}

// BytesToString converts byte slice type to string
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// ConvertNumberToString tries to convert number to string,
// if input is neither number type nor string, it will return error
func ConvertNumberToString(in interface{}) (string, error) {
	inType := reflect.TypeOf(in)

	switch inType.Kind() {
	case reflect.String:
		return in.(string), nil
	case reflect.Bool:
		if in.(bool) == true {
			return constant.TrueString, nil
		}

		return constant.FalseString, nil
	case reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%v", in), nil
	default:
		return constant.EmptyString, errors.New(
			fmt.Sprintf("convert %s to string is not supported. ONLY accept string, float, int, bool.",
				inType.String()))
	}
}

// ConvertInterfaceToSliceInterface converts input data which must be slice type to interface slice,
// it means each element in the slice is interface type.
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

// ConvertInterfaceToMapInterfaceInterface converts input data which must be map type to map interface interface,
// it means each pair of key and value in the map will be interface type
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

// ElementInSlice checks if given element is in the slice
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

// KeyInMap checks if given key is in the map
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

// ValueInMap checks if given value is in the map
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

// TrimSpaceOfStructString trims spaces of each member variable of the struct
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
