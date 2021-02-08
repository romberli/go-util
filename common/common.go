package common

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/jinzhu/now"
	"github.com/romberli/dynamic-struct"

	json "github.com/json-iterator/go"

	"github.com/romberli/go-util/constant"
)

// CombineMessageWithError returns a new string which combines given message and error
func CombineMessageWithError(message string, err error) string {
	if err == nil {
		return message
	}

	return fmt.Sprintf("%s\n%s", message, err.Error())
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

	if inType.Kind() != reflect.Ptr {
		return errors.New("argument must be a pointer to struct")
	}

	inVal := reflect.ValueOf(in).Elem()

	for i := 0; i < inVal.NumField(); i++ {
		f := inVal.Field(i)
		switch f.Kind() {
		case reflect.String:
			if f.CanSet() {
				trimValue := strings.TrimSpace(f.String())
				f.SetString(trimValue)
			}
		}
	}

	return nil
}

// GetValueOfStruct get value of specified field of input struct,
// the field must exist and be exported, otherwise, it will return an error,
// the first argument must be a pointer to struct
func GetValueOfStruct(in interface{}, field string) (interface{}, error) {
	if reflect.TypeOf(in).Kind() != reflect.Ptr {
		return nil, errors.New("first argument must be a pointer to struct")
	}
	v := reflect.ValueOf(in).Elem().FieldByName(field)
	if !v.CanSet() {
		return nil, errors.New(fmt.Sprintf("field %s can not be set, please check if this field is exported", field))
	}

	return v.Interface(), nil
}

// SetValueOfStruct set value of specified field of input struct,
// the field must exist and be exported, otherwise, it will return an error,
// the first argument must be a pointer to struct
func SetValueOfStruct(in interface{}, field string, value interface{}) error {
	if reflect.TypeOf(in).Kind() != reflect.Ptr {
		return errors.New("first argument must be a pointer to struct")
	}

	v := reflect.ValueOf(in).Elem().FieldByName(field)
	if !v.CanSet() {
		return errors.New(fmt.Sprintf("field %s can not be set, please check if this field is exported", field))
	}

	vType := v.Type()
	valueType := reflect.TypeOf(value)
	if vType != valueType {
		return errors.New(fmt.Sprintf("types of field %s(%s) and value(%s) mismatched",
			field, v.Type().String(), valueType.String()))
	}

	v.Set(reflect.ValueOf(value))

	return nil
}

func SetValueOfStructByKind(in interface{}, field string, value interface{}, kind reflect.Kind) error {
	switch kind {
	case reflect.Bool:
		v, err := ConvertToBool(value)
		if err != nil {
			return err
		}

		err = SetValueOfStruct(in, field, v)
		if err != nil {
			return err
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := ConvertToUint(value)
		if err != nil {
			return err
		}

		switch kind {
		case reflect.Int:
			err = SetValueOfStruct(in, field, int(v))
		case reflect.Int8:
			err = SetValueOfStruct(in, field, int8(v))
		case reflect.Int16:
			err = SetValueOfStruct(in, field, int16(v))
		case reflect.Int32:
			err = SetValueOfStruct(in, field, int32(v))
		case reflect.Int64:
			err = SetValueOfStruct(in, field, int64(v))
		case reflect.Uint:
			err = SetValueOfStruct(in, field, uint(v))
		case reflect.Uint8:
			err = SetValueOfStruct(in, field, uint8(v))
		case reflect.Uint16:
			err = SetValueOfStruct(in, field, uint16(v))
		case reflect.Uint32:
			err = SetValueOfStruct(in, field, uint32(v))
		case reflect.Uint64:
			err = SetValueOfStruct(in, field, v)
		}
		if err != nil {
			return err
		}
	case reflect.Float32, reflect.Float64:
		v, err := ConvertToFloat(value)
		if err != nil {
			return err
		}

		switch kind {
		case reflect.Float32:
			err = SetValueOfStruct(in, field, float32(v))
		case reflect.Float64:
			err = SetValueOfStruct(in, field, v)
		}
		if err != nil {
			return err
		}
	case reflect.String:
		v, err := ConvertToString(value)
		if err != nil {
			return err
		}

		err = SetValueOfStruct(in, field, v)
		if err != nil {
			return err
		}
	case reflect.Slice:
		fieldType, ok := reflect.ValueOf(in).Elem().Type().FieldByName(field)
		if !ok {
			return errors.New(fmt.Sprintf("field %s does not exist", field))
		}

		v, err := ConvertToSlice(value, fieldType.Type.Elem().Kind())
		if err != nil {
			return err
		}

		err = SetValueOfStruct(in, field, v)
		if err != nil {
			return err
		}
	case reflect.Struct:
		v, err := ConvertToString(value)
		if err != nil {
			return err
		}

		now.TimeFormats = append(now.TimeFormats, constant.DefaultTimeLayout)
		t, err := now.Parse(v)
		if err != nil {
			return err
		}
		err = SetValueOfStruct(in, field, t)
		if err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("got unsupported reflect.Kind of data type: %s", kind.String()))
	}

	return nil
}

// CopyStructWithFields returns a new struct with only specified fields
// NOTE:
// 1. tags and values of fields are exactly same
// 2. only exported and addressable fields will be copied
// 3. if any field in fields does not exist in the input struct, it returns error
// 4. if values in input struct is a pointer, then value in the new struct will point to the same object
// 5. returning struct is totally a new data type, so you could not use any (*type) assertion
// 6. if fields argument is empty, a new struct which contains the whole fields of input struct will be returned
// 7. technically, for convenience purpose, this function creates a new struct as same as input struct,
//    then removes fields that does not exist in the given fields
func CopyStructWithFields(in interface{}, fields ...string) (interface{}, error) {
	if len(fields) == 0 {
		return CopyStructWithoutFields(in)
	}

	if reflect.TypeOf(in).Kind() != reflect.Ptr {
		return nil, errors.New("first argument must be a pointer to struct")
	}

	var removeFields []string

	inVal := reflect.ValueOf(in).Elem()
	inType := inVal.Type()

	for i := 0; i < inVal.NumField(); i++ {
		fieldName := inType.Field(i).Name
		ok, err := ElementInSlice(fieldName, fields)
		if err != nil {
			return nil, err
		}
		if !ok {
			removeFields = append(removeFields, fieldName)
		}
	}

	return CopyStructWithoutFields(in, removeFields...)
}

// CopyStructWithoutFields returns a new struct without specified fields, there are something you should know.
// NOTE:
// 1. tags and values of remaining fields are exactly same
// 2. only exported and addressable fields will be copied
// 3. if any field in fields does not exist in the input struct, it simply ignores
// 4. if values in input struct is a pointer, then value in the new struct will point to the same object
// 5. returning struct is totally a new data type, so you could not use any (*type) assertion
// 6. if fields argument is empty, a new struct which contains the whole fields of input struct will be returned
func CopyStructWithoutFields(in interface{}, fields ...string) (interface{}, error) {
	newStruct, err := dynamicstruct.MergeStructsWithSettableFields(in)
	if err != nil {
		return nil, err
	}

	for _, field := range fields {
		newStruct = newStruct.RemoveField(field)
	}

	// generate new instance
	newInstance := newStruct.Build().New()
	newVal := reflect.ValueOf(newInstance).Elem()
	newType := newVal.Type()

	inVal := reflect.ValueOf(in).Elem()

	for i := 0; i < newVal.NumField(); i++ {
		fieldType := newType.Field(i)
		fieldVal := newVal.Field(i)
		// set value
		fieldVal.Set(inVal.FieldByName(fieldType.Name))
	}

	return newInstance, nil
}

// MarshalStructWithFields marshals input struct using json.Marshal() with given fields,
// first argument must be a pointer to struct, not the struct itself
func MarshalStructWithFields(in interface{}, fields ...string) ([]byte, error) {
	if reflect.TypeOf(in).Kind() != reflect.Ptr {
		return nil, errors.New("first argument must be a pointer to struct")
	}

	if len(fields) == 0 {
		return json.Marshal(in)
	}

	// generate a new struct with given fields
	newInstance, err := CopyStructWithFields(in, fields...)
	if err != nil {
		return nil, err
	}

	return json.Marshal(newInstance)
}

// MarshalStructWithFields marshals input struct using json.Marshal() without given fields,
// first argument must be a pointer to struct, not the struct itself
func MarshalStructWithoutFields(in interface{}, fields ...string) ([]byte, error) {
	if reflect.TypeOf(in).Kind() != reflect.Ptr {
		return nil, errors.New("first argument must be a pointer to struct")
	}

	if len(fields) == 0 {
		return json.Marshal(in)
	}

	// generate a new struct with fields
	newInstance, err := CopyStructWithoutFields(in, fields...)
	if err != nil {
		return nil, err
	}

	return json.Marshal(newInstance)
}

// MarshalStructWithTag marshals input struct using json.Marshal() with fields that contain given tag,
// first argument must be a pointer to struct, not the struct itself
func MarshalStructWithTag(in interface{}, tag string) ([]byte, error) {
	if reflect.TypeOf(in).Kind() != reflect.Ptr {
		return nil, errors.New("first argument must be a pointer to struct")
	}
	if tag == constant.EmptyString {
		return nil, errors.New("tag should not be empty")
	}

	inVal := reflect.ValueOf(in).Elem()

	var fields []string

	for i := 0; i < inVal.NumField(); i++ {
		fieldType := inVal.Type().Field(i)
		fieldTag := fieldType.Tag.Get(tag)
		if fieldTag != constant.EmptyString {
			fields = append(fields, fieldType.Name)
		}
	}

	return MarshalStructWithFields(in, fields...)
}

// NewMapWithStructTag returns a new map as the result,
// it loops the keys of given map and tags of the struct,
// if key and tag are same, the field of the the input struct will be the key of the new map,
// the value of the given map will be the value of the new map,
// if any key in the given map could not match any tag in the struct,
// it will return error, so make sure that each key the given map could match a field tag in the struct
func NewMapWithStructTag(m map[string]interface{}, in interface{}, tag string) (map[string]interface{}, error) {
	if reflect.TypeOf(in).Kind() != reflect.Ptr {
		return nil, errors.New("second argument must be a pointer to struct")
	}

	result := make(map[string]interface{})
	inVal := reflect.ValueOf(in).Elem()

Loop:
	for key := range m {
		for i := 0; i < inVal.NumField(); i++ {
			fieldType := inVal.Type().Field(i)
			fieldTag := fieldType.Tag.Get(tag)

			if key == fieldTag {
				result[fieldType.Name] = m[key]
				continue Loop
			}
		}
		// this means there is no relevant tag in the struct, should return error
		return nil, errors.New(fmt.Sprintf("key %s could not match any tag in the struct", key))
	}

	return result, nil
}

// UnmarshalToMapWithStructTag returns a map as the result, it works as following logic:
// 1. unmarshals given data to to a temporary map to get keys
// 2. unmarshals given data to the input struct, to get field names and values with appropriate data types
// 3. loops keys in the temporary map, loops tags of the input struct in a nested
// 4. if the key matches the tag, set field name as the key of result map, set field value as the value of the result map
// 5. if any key in the temporary map can not match any field tag of the struct, it returns error,
//    so make sure that each key of the given data could match a field tag in the struct
func UnmarshalToMapWithStructTag(data []byte, in interface{}, tag string) (map[string]interface{}, error) {
	if reflect.TypeOf(in).Kind() != reflect.Ptr {
		return nil, errors.New("second argument must be a pointer to struct")
	}

	// get new decoder to unmarshal with specified tag
	decoder := json.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		TagKey:                 tag,
	}.Froze()
	// unmarshal to struct to get appropriate data type
	err := decoder.Unmarshal(data, &in)
	if err != nil {
		return nil, err
	}
	// unmarshal to temporary map to get key names
	tmpMap := make(map[string]interface{})
	err = decoder.Unmarshal(data, &tmpMap)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	inVal := reflect.ValueOf(in).Elem()
Loop:
	for key := range tmpMap {
		for i := 0; i < inVal.NumField(); i++ {
			fieldType := inVal.Type().Field(i)
			fieldTag := fieldType.Tag.Get(tag)
			if key == fieldTag {
				// set field name of struct as key, set value of temporary map as value of the result
				result[fieldType.Name] = inVal.Field(i).Interface()
				continue Loop
			}
		}
		// this means there is no relevant tag in the struct, should return error
		return nil, errors.New(fmt.Sprintf("key %s could not match any tag in the struct", key))
	}

	return result, nil
}

// UnmarshalToMapWithStructTag converts given string to []byte and then call UnmarshalToMapWithStructTag() function
func UnmarshalToMapWithStructTagFromString(s string, in interface{}, tag string) (map[string]interface{}, error) {
	return UnmarshalToMapWithStructTag([]byte(s), in, tag)
}
