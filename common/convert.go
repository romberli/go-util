package common

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/pingcap/errors"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/types"
)

// StringToBytes converts string type to byte slice
func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// BytesToString converts byte slice type to string
func BytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
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
		return constant.EmptyString, errors.Errorf("convert %s to string is not supported. ONLY accept string, float, int, bool.",
			inType.String())
	}
}

// ConvertInterfaceToString converts struct pointer to string
func ConvertInterfaceToString(in interface{}) string {
	val := reflect.ValueOf(in)

	if val.Kind() != reflect.Ptr {
		return fmt.Sprintf("%v", in)
	}
	if val.IsNil() {
		return "<nil>"
	}

	elem := val.Elem()

	if elem.Kind() != reflect.Struct {
		return fmt.Sprintf("%v", in)
	}

	var builder strings.Builder
	builder.WriteString(elem.Type().Name())
	builder.WriteString(constant.LeftBraceString)

	for i := constant.ZeroInt; i < elem.NumField(); i++ {
		field := elem.Type().Field(i)
		fieldValue := elem.Field(i)
		if field.PkgPath != constant.EmptyString {
			continue
		}
		if i > constant.ZeroInt {
			builder.WriteString(constant.CommaString + constant.SpaceString)
		}

		builder.WriteString(field.Name)
		builder.WriteString(constant.ColonString)

		switch fieldValue.Kind() {
		case reflect.Slice, reflect.Array:
			builder.WriteString(formatSlice(fieldValue))
		case reflect.Map:
			builder.WriteString(formatMap(fieldValue))
		case reflect.Ptr:
			if fieldValue.IsNil() {
				builder.WriteString("<nil>")
			} else {
				builder.WriteString(fmt.Sprintf("%v", fieldValue.Interface()))
			}
		default:
			builder.WriteString(fmt.Sprintf("%v", fieldValue.Interface()))
		}
	}

	builder.WriteString(constant.RightBraceString)

	return builder.String()
}

// formatSlice formats slice type
func formatSlice(v reflect.Value) string {
	if v.Len() == constant.ZeroInt {
		return constant.EmptyListString
	}

	var builder strings.Builder
	builder.WriteString(constant.LeftBracketString)

	for i := constant.ZeroInt; i < v.Len(); i++ {
		if i > constant.ZeroInt {
			builder.WriteString(constant.CommaString + constant.SpaceString)
		}
		builder.WriteString(fmt.Sprintf("%v", v.Index(i).Interface()))
	}

	builder.WriteString(constant.RightBracketString)
	return builder.String()
}

// formatMap formats map type
func formatMap(v reflect.Value) string {
	if v.Len() == constant.ZeroInt {
		return constant.EmptyMapString
	}

	var builder strings.Builder
	builder.WriteString(constant.LeftBraceString)

	keys := v.MapKeys()
	for i, key := range keys {
		if i > 0 {
			builder.WriteString(constant.CommaString + constant.SpaceString)
		}
		value := v.MapIndex(key)
		builder.WriteString(fmt.Sprintf("%v:%v", key.Interface(), value.Interface()))
	}

	builder.WriteString(constant.RightBraceString)
	return builder.String()
}

// ConvertInterfaceToSliceInterface converts input data which must be slice type to interface slice,
// it means each element in the slice is interface type.
func ConvertInterfaceToSliceInterface(in interface{}) ([]interface{}, error) {
	inType := reflect.TypeOf(in)
	inValue := reflect.ValueOf(in)

	if inType.Kind() != reflect.Slice {
		return nil, errors.New("the argument must be array or slice")
	}

	inLength := inValue.Len()
	sliceInterface := make([]interface{}, inLength)

	for i := 0; i < inLength; i++ {
		sliceInterface[i] = inValue.Index(i).Interface()
	}

	return sliceInterface, nil
}

// ConvertInterfaceToMapInterfaceInterface converts input data which must be map type to map[interface]interface,
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

func ConvertBytesToSlice(bytes []byte, kind reflect.Kind) (interface{}, error) {
	switch kind {
	case reflect.Bool:
		result := make([]bool, constant.ZeroInt)
		err := json.Unmarshal(bytes, &result)
		if err != nil {
			return nil, errors.Trace(err)
		}

		return result, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result := make([]int, constant.ZeroInt)
		err := json.Unmarshal(bytes, &result)
		if err != nil {
			return nil, errors.Trace(err)
		}

		return result, nil
	case reflect.Float32, reflect.Float64:
		result := make([]float64, constant.ZeroInt)
		err := json.Unmarshal(bytes, &result)
		if err != nil {
			return nil, errors.Trace(err)
		}

		return result, nil
	case reflect.String:
		result := make([]string, constant.ZeroInt)
		err := json.Unmarshal(bytes, &result)
		if err != nil {
			return nil, errors.Trace(err)
		}

		return result, nil
	default:
		return nil, errors.Errorf("unsupported data type: %s", kind.String())
	}
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
			return false, errors.Errorf("bool value of integer type should be either 0 or 1, %d is not valid", v)
		}
	case []byte, string:
		v, err := ConvertToString(in)
		if err != nil {
			return false, err
		}

		switch v {
		case constant.TrueString:
			return true, nil
		case constant.FalseString:
			return false, nil
		default:
			return false, errors.Errorf("bool value of string type should be either true or false, %s is not valid", in.(string))
		}
	default:
		return false, errors.Errorf("can not convert to a valid bool value, %v is not valid", in)
	}
}

func ConvertToInt64(in interface{}) (int64, error) {
	switch v := in.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case string:
		value, err := strconv.ParseInt(v, 10, 64)

		return int64(value), errors.Trace(err)
	case []byte:
		value, err := strconv.ParseInt(string(v), 10, 64)

		return int64(value), errors.Trace(err)
	case nil:
		return int64(constant.ZeroInt), nil
	default:
		return int64(constant.ZeroInt), errors.Errorf("unsupported data type: %T", v)
	}
}

func ConvertToUint64(in interface{}) (uint64, error) {
	value, err := ConvertToInt64(in)
	if err != nil {
		return constant.ZeroInt, err
	}

	return uint64(value), nil
}

func ConvertToInt(in interface{}) (int, error) {
	value, err := ConvertToInt64(in)
	if err != nil {
		return constant.ZeroInt, err
	}

	return int(value), nil
}

func ConvertToUint(in interface{}) (uint, error) {
	value, err := ConvertToInt64(in)
	if err != nil {
		return uint(constant.ZeroInt), err
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
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return constant.ZeroInt, errors.Trace(err)
		}

		return value, nil
	case []byte:
		value, err := strconv.ParseFloat(string(v), 64)
		if err != nil {
			return constant.ZeroInt, errors.Trace(err)
		}

		return value, nil
	case nil:
		return constant.ZeroInt, nil
	default:
		return constant.ZeroInt, errors.Errorf("unsupported data type: %T", v)
	}
}

func ConvertToString(in interface{}) (string, error) {
	switch v := in.(type) {
	case string:
		return v, nil
	case []byte:
		return BytesToString(v), nil
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case nil:
		return constant.EmptyString, nil
	case time.Time:
		return v.Format(constant.DefaultTimeLayout), nil
	default:
		return constant.EmptyString, errors.Errorf("unsupported data type: %T", v)
	}
}

func ConvertToSlice(in interface{}, kind reflect.Kind) (interface{}, error) {
	inKind := reflect.TypeOf(in).Kind()
	if inKind != reflect.String && inKind != reflect.Slice {
		return nil, errors.Errorf("value must be a string or a slice, not %s", inKind.String())
	}

	if inKind == reflect.String {
		inStr, ok := in.(string)
		if !ok {
			return nil, errors.New("convert to string failed")
		}
		inBytes := []byte(inStr)
		if json.Valid(inBytes) {
			val, err := ConvertBytesToSlice(inBytes, kind)
			if err != nil {
				return nil, err
			}

			return val, nil
		} else {
			return nil, errors.New("invalid json string")
		}
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
			value, err := ConvertToUint64(element)
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
			value, err := ConvertToUint64(element)
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
			value, err := ConvertToUint64(element)
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
			value, err := ConvertToUint64(element)
			if err != nil {
				return nil, err
			}
			result[i] = value
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
			value, err := ConvertToInt64(element)
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
			value, err := ConvertToInt64(element)
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
			value, err := ConvertToInt64(element)
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
			value, err := ConvertToInt64(element)
			if err != nil {
				return nil, err
			}
			result[i] = value
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
		return nil, errors.Errorf("kind must be one of [reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64, reflect.String], %s is not valid", kind.String())
	}
}

// ConvertSliceToString converts the slice to string, it uses c as the separation character
func ConvertSliceToString[T types.Primitive](s []T, c string) string {
	var result string

	for _, v := range s {
		result += fmt.Sprintf("%v%s", v, c)
	}

	return strings.Trim(result, c)
}

// ConvertInterfaceSliceToString converts the interface slice to string, it uses c as the separation character
func ConvertInterfaceSliceToString(s []interface{}, c string) string {
	var result string

	for _, v := range s {
		result += fmt.Sprintf("%v%s", v, c)
	}

	return strings.Trim(result, c)
}
