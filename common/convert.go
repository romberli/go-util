package common

import (
	"fmt"
	"github.com/pingcap/errors"
	"github.com/siddontang/go/hack"

	"reflect"
	"strconv"
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

func ConvertToBool(in interface{}) (bool, error) {
	switch in.(type) {
	case bool:
		return in.(bool), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		v, err := ConvertToInt(in)
		if err != nil {
			return false, err
		}

		switch v {
		case 0:
			return false, nil
		case 1:
			return true, nil
		default:
			return false, errors.New(fmt.Sprintf("bool type value should be either 0 or 1, %d is not valid", v))
		}
	default:
		return false, errors.New(fmt.Sprintf("can not convert to a valid bool value, %v is not valid", in))
	}
}

func ConvertToInt(in interface{}) (int, error) {
	switch v := in.(type) {
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		return int(v), nil
	case float32:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		value, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return constant.ZeroInt, err
		}

		return int(value), nil
	case []byte:
		value, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return constant.ZeroInt, err
		}

		return int(value), nil
	case nil:
		return constant.ZeroInt, nil
	default:
		return constant.ZeroInt, errors.Errorf("data type is %T", v)
	}
}

func ConvertToUint(in interface{}) (uint, error) {
	value, err := ConvertToInt(in)
	if err != nil {
		return constant.ZeroInt, err
	}

	return uint(value), nil
}

func ConvertToFloat(in interface{}) (float64, error) {
	switch v := in.(type) {
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(v, 64)
	case []byte:
		return strconv.ParseFloat(string(v), 64)
	case nil:
		return constant.ZeroInt, nil
	default:
		return constant.ZeroInt, errors.Errorf("data type is %T", v)
	}
}

func ConvertToString(in interface{}) (string, error) {
	switch v := in.(type) {
	case string:
		return v, nil
	case []byte:
		return hack.String(v), nil
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case nil:
		return constant.EmptyString, nil
	default:
		return constant.EmptyString, errors.Errorf("data type is %T", v)
	}
}

func ConvertToSlice(in interface{}, kind reflect.Kind) (interface{}, error) {
	inKind := reflect.TypeOf(in).Kind()
	if inKind != reflect.Slice {
		return nil, errors.New(fmt.Sprintf("value must be a slice, not %s", inKind.String()))
	}

	inVal := reflect.ValueOf(in)

	switch kind {
	case reflect.Uint:
		result := make([]uint, inVal.Len())
		for i := 0; i < inVal.Len(); i++ {
			element := inVal.Index(i).Interface()
			value, err := ConvertToUint(element)
			if err != nil {
				return nil, err
			}
			result[i] = value
		}

		return result, nil
	case reflect.Uint8:
		result := make([]uint8, inVal.Len())
		for i := 0; i < inVal.Len(); i++ {
			element := inVal.Index(i).Interface()
			value, err := ConvertToUint(element)
			if err != nil {
				return nil, err
			}
			result[i] = uint8(value)
		}

		return result, nil
	case reflect.Uint16:
		result := make([]uint16, inVal.Len())
		for i := 0; i < inVal.Len(); i++ {
			element := inVal.Index(i).Interface()
			value, err := ConvertToUint(element)
			if err != nil {
				return nil, err
			}
			result[i] = uint16(value)
		}

		return result, nil
	case reflect.Uint32:
		result := make([]uint32, inVal.Len())
		for i := 0; i < inVal.Len(); i++ {
			element := inVal.Index(i).Interface()
			value, err := ConvertToUint(element)
			if err != nil {
				return nil, err
			}
			result[i] = uint32(value)
		}

		return result, nil
	case reflect.Uint64:
		result := make([]uint64, inVal.Len())
		for i := 0; i < inVal.Len(); i++ {
			element := inVal.Index(i).Interface()
			value, err := ConvertToUint(element)
			if err != nil {
				return nil, err
			}
			result[i] = uint64(value)
		}

		return result, nil
	case reflect.Int:
		result := make([]int, inVal.Len())
		for i := 0; i < inVal.Len(); i++ {
			element := inVal.Index(i).Interface()
			value, err := ConvertToInt(element)
			if err != nil {
				return nil, err
			}
			result[i] = value
		}

		return result, nil
	case reflect.Int8:
		result := make([]int8, inVal.Len())
		for i := 0; i < inVal.Len(); i++ {
			element := inVal.Index(i).Interface()
			value, err := ConvertToInt(element)
			if err != nil {
				return nil, err
			}
			result[i] = int8(value)
		}

		return result, nil
	case reflect.Int16:
		result := make([]int16, inVal.Len())
		for i := 0; i < inVal.Len(); i++ {
			element := inVal.Index(i).Interface()
			value, err := ConvertToInt(element)
			if err != nil {
				return nil, err
			}
			result[i] = int16(value)
		}

		return result, nil
	case reflect.Int32:
		result := make([]int32, inVal.Len())
		for i := 0; i < inVal.Len(); i++ {
			element := inVal.Index(i).Interface()
			value, err := ConvertToInt(element)
			if err != nil {
				return nil, err
			}
			result[i] = int32(value)
		}

		return result, nil
	case reflect.Int64:
		result := make([]int64, inVal.Len())
		for i := 0; i < inVal.Len(); i++ {
			element := inVal.Index(i).Interface()
			value, err := ConvertToInt(element)
			if err != nil {
				return nil, err
			}
			result[i] = int64(value)
		}

		return result, nil
	case reflect.Float32:
		result := make([]float32, inVal.Len())
		for i := 0; i < inVal.Len(); i++ {
			element := inVal.Index(i).Interface()
			value, err := ConvertToFloat(element)
			if err != nil {
				return nil, err
			}
			result[i] = float32(value)
		}

		return result, nil
	case reflect.Float64:
		result := make([]float64, inVal.Len())
		for i := 0; i < inVal.Len(); i++ {
			element := inVal.Index(i).Interface()
			value, err := ConvertToFloat(element)
			if err != nil {
				return nil, err
			}
			result[i] = value
		}

		return result, nil
	case reflect.String:
		result := make([]string, inVal.Len())
		for i := 0; i < inVal.Len(); i++ {
			element := inVal.Index(i).Interface()
			value, err := ConvertToString(element)
			if err != nil {
				return nil, err
			}
			result[i] = value
		}

		return result, nil
	default:
		return nil, errors.New(fmt.Sprintf("kind must be one of [reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64, reflect.String], %s is not valid", kind.String()))
	}
}
