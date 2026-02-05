package common

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/pingcap/errors"
	"github.com/tidwall/gjson"

	"github.com/romberli/go-util/constant"
)

const (
	// sensitive keyword
	DefaultSensitivePassKeyword   = "pass"
	DefaultSensitiveSecretKeyword = "secret"
	DefaultSensitivePwdKeyword    = "pwd"
)

var (
	DefaultSensitiveKeywords = []string{
		DefaultSensitivePassKeyword,
		DefaultSensitiveSecretKeyword,
		DefaultSensitivePwdKeyword,
	}
)

// KeyExists checks if the key exists in the json data
func KeyExists(data []byte, keys ...string) bool {
	_, _, _, err := jsonparser.Get(data, keys...)
	if err != nil {
		return false
	}

	return true
}

// KeyPathExists checks if the key path exists in the json data
func KeyPathExists(data []byte, path string) bool {
	_, _, _, err := jsonparser.Get(data, strings.Split(path, constant.DotString)...)
	if err != nil {
		return false
	}

	return true
}

// GetLength returns the length of the json array, the value of the path should be a json array
func GetLength(data []byte, path string) (int, error) {
	if !gjson.ValidBytes(data) {
		return constant.ZeroInt, errors.Errorf("invalid json data. data: %s", BytesToString(data))
	}

	result := gjson.GetBytes(data, path)
	if !result.Exists() {
		return constant.ZeroInt, errors.Errorf("path not found. path: %s", path)
	}
	if !result.IsArray() {
		return constant.ZeroInt, errors.Errorf("path is not an array. path: %s", path)
	}

	return len(result.Array()), nil
}

// SerializeBytes serializes the struct to map[string]interface{} and []byte to string
func SerializeBytes(v interface{}) interface{} {
	val := reflect.ValueOf(v)

	switch val.Kind() {
	case reflect.Ptr:
		if val.IsNil() {
			return nil
		}
		elem := val.Elem().Interface()
		return SerializeBytes(elem)
	case reflect.Struct:
		serialized := make(map[string]interface{})
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			if !val.Field(i).CanInterface() {
				// TODO: for now, simply skip the field that cannot be set
				continue
			}
			// get json tag
			tag := field.Tag.Get(constant.DefaultJSONTag)
			// if tag is "-", skip this field
			if tag == constant.DashString {
				continue
			}
			key := field.Name
			if tag != constant.EmptyString {
				tagParts := strings.Split(tag, constant.CommaString)
				if tagParts[constant.ZeroInt] != constant.EmptyString {
					key = tagParts[constant.ZeroInt]
				}
			}

			fieldValue := val.Field(i).Interface()
			serialized[key] = SerializeBytes(fieldValue)
		}
		return serialized
	case reflect.Slice:
		if val.Type().Elem().Kind() == reflect.Uint8 {
			// convert []byte to string
			return string(val.Bytes())
		}
		// other slice types
		serialized := make([]interface{}, val.Len())
		for i := 0; i < val.Len(); i++ {
			serialized[i] = SerializeBytes(val.Index(i).Interface())
		}
		return serialized
	case reflect.Map:
		serialized := make(map[string]interface{})
		for _, key := range val.MapKeys() {
			keyStr := fmt.Sprintf("%v", key)
			serialized[keyStr] = SerializeBytes(val.MapIndex(key).Interface())
		}
		return serialized
	default:
		return v
	}
}

// DeserializeBytes deserializes the map[string]interface{} to struct and string to []byte
func DeserializeBytes(v interface{}, t reflect.Type) interface{} {
	val := reflect.ValueOf(v)

	switch t.Kind() {
	case reflect.Ptr:
		if v == nil || val.IsNil() {
			return reflect.Zero(t).Interface()
		}
		elemType := t.Elem()
		elemValue := DeserializeBytes(val.Interface(), elemType)
		ptrValue := reflect.New(elemType)
		ptrValue.Elem().Set(reflect.ValueOf(elemValue))
		return ptrValue.Interface()
	case reflect.Struct:
		structValue := reflect.New(t).Elem()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if !structValue.Field(i).CanSet() {
				// TODO: for now, simply skip the field that cannot be set
				continue
			}
			// get json tag
			tag := field.Tag.Get(constant.DefaultJSONTag)
			// if tag is "-", skip this field
			if tag == constant.DashString {
				continue
			}
			key := field.Name
			if tag != constant.EmptyString {
				tagParts := strings.Split(tag, constant.CommaString)
				if tagParts[constant.ZeroInt] != constant.EmptyString {
					key = tagParts[constant.ZeroInt]
				}
			}

			fieldValue := DeserializeBytes(val.MapIndex(reflect.ValueOf(key)).Interface(), field.Type)
			convertedValue := reflect.ValueOf(fieldValue).Convert(field.Type)
			structValue.Field(i).Set(convertedValue)
		}
		return structValue.Interface()
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			// convert string to []byte
			return []byte(val.String())
		}
		// other slice types
		sliceValue := reflect.MakeSlice(t, val.Len(), val.Len())
		for i := 0; i < val.Len(); i++ {
			sliceValue.Index(i).Set(reflect.ValueOf(DeserializeBytes(val.Index(i).Interface(), t.Elem())))
		}
		return sliceValue.Interface()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val.Kind() == reflect.Float64 {
			return int(val.Float())
		}
		// convert to int
		return int(val.Float())
	case reflect.Float32, reflect.Float64:
		return val.Float()
	case reflect.String:
		return val.String()
	case reflect.Map:
		mapType := reflect.MapOf(t.Key(), t.Elem())
		mapValue := reflect.MakeMap(mapType)
		for _, key := range val.MapKeys() {
			mapKey := reflect.ValueOf(key.Interface()).Convert(t.Key())
			mapElem := DeserializeBytes(val.MapIndex(key).Interface(), t.Elem())
			mapValue.SetMapIndex(mapKey, reflect.ValueOf(mapElem))
		}
		return mapValue.Interface()
	default:
		return v
	}
}

// MaskJSON masks the sensitive fields in the json body
func MaskJSON(jsonBytes []byte, sensitiveFields []string, excludes ...string) ([]byte, error) {
	if len(jsonBytes) == constant.ZeroInt {
		return jsonBytes, nil
	}

	var data interface{}
	err := json.Unmarshal(jsonBytes, &data)
	if err != nil {
		return jsonBytes, errors.Trace(err)
	}

	maskedValue := maskValue(data, sensitiveFields, excludes...)

	result, err := json.Marshal(maskedValue)
	if err != nil {
		return jsonBytes, errors.Trace(err)
	}

	return result, nil
}

func maskValue(value interface{}, sensitiveFields []string, excludes ...string) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		for key, val := range v {
			if isSensitiveField(key, sensitiveFields, excludes...) {
				v[key] = constant.DefaultMaskedValue
			} else {
				// mask value recursively
				v[key] = maskValue(val, sensitiveFields, excludes...)
			}
		}
		return v

	case []interface{}:
		for i, item := range v {
			v[i] = maskValue(item, sensitiveFields, excludes...)
		}
		return v

	default:
		return v
	}
}

// isSensitiveField checks if the field name contains any of the sensitive fields
func isSensitiveField(fieldName string, sensitiveFields []string, excludes ...string) bool {
	lowerField := strings.ToLower(fieldName)
	for _, exclude := range excludes {
		if strings.Contains(lowerField, exclude) {
			return false
		}
	}
	for _, sensitiveField := range sensitiveFields {
		if strings.Contains(lowerField, sensitiveField) {
			return true
		}
	}

	return false
}
