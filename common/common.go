package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/romberli/dynamic-struct"

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

// CopyStructWithFields returns a new struct with only specified fields
// NOTE:
// 1. tags and values of fields are exactly same
// 2. only exported and addressable fields will be copied
// 3. if any field in fields does not exist in the input struct, it returns error
// 4. if values in input struct is a pointer, then value in the new struct will point to the same object
// 5. returning struct is totally a new data type, so you could not use any (*type) assertion
// 6. technically, for convenience purpose, this function creates a new struct as same as input struct,
//    then removes fields that does not exist in the given fields
func CopyStructWithFields(in interface{}, fields ...string) (interface{}, error) {
	if len(fields) == 0 {
		return in, nil
	}

	if reflect.TypeOf(in).Kind() != reflect.Ptr {
		return nil, errors.New("first argument must be a pointer to struct")
	}

	var removeFields []string

	inVal := reflect.ValueOf(in).Elem()
	inType := inVal.Type()

	for i := 0; i < inVal.NumField(); i++ {
		inField := inType.Field(i).Name
		ok, err := ElementInSlice(inField, fields)
		if err != nil {
			return nil, err
		}
		if !ok {
			removeFields = append(removeFields, inField)
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
func CopyStructWithoutFields(in interface{}, fields ...string) (interface{}, error) {
	if len(fields) == 0 {
		return in, nil
	}

	newStruct, err := dynamicstruct.MergeStructsWithSettableFields(in)
	if err != nil {
		return nil, err
	}

	for _, field := range fields {
		newStruct = newStruct.RemoveField(field)
	}

	// generate new instance
	newInstance := newStruct.Build().New()
	newValue := reflect.ValueOf(newInstance).Elem()
	newType := newValue.Type()

	inputValue := reflect.ValueOf(in).Elem()

	for i := 0; i < newValue.NumField(); i++ {
		fType := newType.Field(i)
		fValue := newValue.Field(i)
		// set value
		fValue.Set(inputValue.FieldByName(fType.Name))
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

// NewMapWithStructTag returns a new map, it loops the keys of given map and tags of the struct,
// if key and tag are same, the field of the the input struct will be the key of the new map,
// the value of the given map will be the value of the new map,
// if any key in the given map could not match any tag in the struct,
// it will return error, so make sure that each key the given map could match a field tag in the struct
func NewMapWithStructTag(m map[string]interface{}, in interface{}, tag string) (map[string]interface{}, error) {
	if reflect.TypeOf(in).Kind() != reflect.Ptr {
		return nil, errors.New("second argument must be a pointer to struct")
	}

	newMap := make(map[string]interface{})
	inVal := reflect.ValueOf(in).Elem()

Loop:
	for key := range m {
		for i := 0; i < inVal.NumField(); i++ {
			fieldType := inVal.Type().Field(i)
			fieldTag := fieldType.Tag.Get(tag)

			if key == fieldTag {
				newMap[fieldType.Name] = m[key]
				continue Loop
			}
		}
		// this means there is no relevant tag in the struct, should return error
		return nil, errors.New(fmt.Sprintf("key %s could not match any tag in the struct", key))
	}

	return newMap, nil
}
