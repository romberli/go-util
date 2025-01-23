package common

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/pingcap/errors"
	"github.com/romberli/dynamic-struct"

	json "github.com/json-iterator/go"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/types"
)

// ReverseString reverses a string
func ReverseString(s string) string {
	var reverse string
	for _, v := range s {
		reverse = string(v) + reverse
	}

	return reverse
}

// SetRandomValueToNil set each value in slice values if value is a random value
func SetRandomValueToNil(values ...interface{}) error {
	for i, value := range values {
		if value == nil {
			continue
		}

		isRandom, err := IsRandomValue(value)
		if err != nil {
			return err
		}

		if isRandom {
			values[i] = nil
		}
	}

	return nil
}

// IsRandomValue checks if given value is a random value
func IsRandomValue(value interface{}) (bool, error) {
	switch value.(type) {
	case int:
		if value.(int) == constant.DefaultRandomInt {
			return true, nil
		}
	case int8:
		if int(value.(int8)) == constant.DefaultRandomInt {
			return true, nil
		}
	case int16:
		if int(value.(int16)) == constant.DefaultRandomInt {
			return true, nil
		}
	case int32:
		if int(value.(int32)) == constant.DefaultRandomInt {
			return true, nil
		}
	case int64:
		if int(value.(int64)) == constant.DefaultRandomInt {
			return true, nil
		}
	case uint:
		if int(value.(uint)) == constant.DefaultRandomInt {
			return true, nil
		}
	case uint8:
		if int(value.(uint8)) == constant.DefaultRandomInt {
			return true, nil
		}
	case uint16:
		if int(value.(uint16)) == constant.DefaultRandomInt {
			return true, nil
		}
	case uint32:
		if int(value.(uint32)) == constant.DefaultRandomInt {
			return true, nil
		}
	case uint64:
		if int(value.(uint64)) == constant.DefaultRandomInt {
			return true, nil
		}
	case float32:
		if float64(value.(float32)) == constant.DefaultRandomFloat {
			return true, nil
		}
	case float64:
		if value.(float64) == constant.DefaultRandomFloat {
			return true, nil
		}
	case string:
		if value.(string) == constant.DefaultRandomString {
			return true, nil
		}
	case time.Time:
		if value.(time.Time).Format(constant.DefaultTimeLayout) == constant.DefaultRandomTimeString {
			return true, nil
		}
	default:
		val := reflect.ValueOf(value)
		kind := val.Kind()
		switch kind {
		case reflect.Ptr, reflect.Slice, reflect.Map:
			if val.IsNil() {
				return true, nil
			}
		default:
			valueType := reflect.TypeOf(value)
			if valueType.ConvertibleTo(reflect.TypeOf(constant.DefaultRandomInt)) {
				newVal := val.Convert(reflect.TypeOf(constant.DefaultRandomInt)).Interface().(int)
				if newVal == constant.DefaultRandomInt {
					return true, nil
				}

				return false, nil
			}

			if valueType.ConvertibleTo(reflect.TypeOf(constant.DefaultRandomFloat)) {
				newVal := val.Convert(reflect.TypeOf(constant.DefaultRandomFloat)).Interface().(float64)
				if newVal == constant.DefaultRandomFloat {
					return true, nil
				}

				return false, nil
			}

			if valueType.ConvertibleTo(reflect.TypeOf(constant.DefaultRandomString)) {
				newVal := val.Convert(reflect.TypeOf(constant.DefaultRandomString)).Interface().(string)
				if newVal == constant.DefaultRandomString {
					return true, nil
				}

				return false, nil
			}

			return false, errors.Errorf("unsupported data type: %T", value)
		}
	}

	return false, nil
}

// CombineMessageWithError returns a new string which combines given message and error
func CombineMessageWithError(message string, err error) string {
	if err == nil {
		return message
	}

	return fmt.Sprintf("%s\n%s", message, err.Error())
}

// StringKeyInMap checks if a string key is in the map
func StringKeyInMap(m map[string]string, str string) bool {
	if _, ok := m[str]; ok {
		return true
	}

	return false
}

// ElementEqualOrderInSlice checks if given elements are same in the slices,
// note that the order of elements in the slices must be same
func ElementEqualOrderInSlice[T types.Primitive](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true
}

// ElementEqualInSlice checks if given elements are same in the slices,
// note that the order of elements in the slices is not important
func ElementEqualInSlice[T types.Primitive](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}

	for _, v1 := range s1 {
		if !ElementInSlice(s2, v1) {
			return false
		}
	}

	return true
}

// ElementInSlice checks if given element is in the slice
func ElementInSlice[T types.Primitive](s []T, e T) bool {
	for _, element := range s {
		if element == e {
			return true
		}
	}

	return false
}

// ElementInSortedSlice checks if given element is in the sorted slice,
// note that the slice must be sorted in ascending order
func ElementInSortedSlice[T types.Number](s []T, e T) bool {
	if len(s) == constant.ZeroInt {
		return false
	}
	if e < s[constant.ZeroInt] || e > s[len(s)-constant.OneInt] {
		return false
	}

	for _, element := range s {
		if element == e {
			return true
		}

		if element > e {
			return false
		}
	}

	return false
}

// ElementInSliceInterface checks if given element is in the slice
func ElementInSliceInterface(s interface{}, e interface{}) (bool, error) {
	kind := reflect.TypeOf(s).Kind()
	sValue := reflect.ValueOf(s)

	if kind != reflect.Slice {
		return false, errors.Errorf("first argument must be array or slice, %s is not valid", kind.String())
	}

	for i := constant.ZeroInt; i < sValue.Len(); i++ {
		if reflect.DeepEqual(e, sValue.Index(i).Interface()) {
			return true, nil
		}
	}

	return false, nil
}

// KeyInMap checks if given key is in the map
func KeyInMap(m interface{}, k interface{}) (bool, error) {
	kind := reflect.TypeOf(m).Kind()
	if kind != reflect.Map {
		return false, errors.Errorf("first argument must be map, %s is not valid", kind.String())
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
func ValueInMap(m interface{}, v interface{}) (bool, error) {
	kind := reflect.TypeOf(m).Kind()
	if kind != reflect.Map {
		return false, errors.Errorf("first argument must be map, %s is not valid", kind.String())
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
	kind := reflect.TypeOf(in).Kind()
	if kind != reflect.Ptr {
		return errors.Errorf("first must be a pointer to struct, %s is not valid", kind.String())
	}

	inVal := reflect.ValueOf(in).Elem()

	for i := constant.ZeroInt; i < inVal.NumField(); i++ {
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
// the fields must exist and be exported, otherwise, it will return an error,
// the first argument must be a pointer to struct
func GetValueOfStruct(in interface{}, field string) (interface{}, error) {
	kind := reflect.TypeOf(in).Kind()
	if kind != reflect.Ptr {
		return nil, errors.Errorf("first must be a pointer to struct, %s is not valid", kind.String())
	}
	v := reflect.ValueOf(in).Elem().FieldByName(field)
	if !v.CanSet() {
		return nil, errors.Errorf("field %s can not be set, please check if this field is exported", field)
	}

	return v.Interface(), nil
}

// SetValueOfStruct sets value of specified field of input struct,
// the fields must exist and be exported, otherwise, it will return an error,
// the first argument must be a pointer to struct
// if value is nil, the field value will be set to ZERO value of the type
func SetValueOfStruct(in interface{}, field string, value interface{}) error {
	kind := reflect.TypeOf(in).Kind()
	if kind != reflect.Ptr {
		return errors.Errorf("first must be a pointer to struct, %s is not valid", kind.String())
	}

	v := reflect.ValueOf(in).Elem().FieldByName(field)
	if !v.IsValid() {
		return errors.Errorf("field does not exist. field: %s", field)
	}
	if !v.CanSet() {
		return errors.Errorf("field can not be set, please check if this field is exported. filed: %s", field)
	}

	vType := v.Type()
	valueType := reflect.TypeOf(value)

	if valueType == nil {
		// set zero value
		v.Set(reflect.Zero(vType))
		return nil
	}

	if vType != valueType {
		if valueType.ConvertibleTo(vType) {
			// convert value type to field type
			v.Set(reflect.ValueOf(value).Convert(vType))
			return nil
		}

		return errors.Errorf("types of field and value mismatched. field: %s, field type: %s, value type: %s",
			field, v.Type().String(), valueType.String())
	}

	// set value
	v.Set(reflect.ValueOf(value))

	return nil
}

// SetValuesWithMap sets values of input struct with given map,
// the fields of map must exist and be exported, otherwise, it will return an error,
// the first argument must be a pointer to struct
func SetValuesWithMap(in interface{}, fields map[string]interface{}) error {
	kind := reflect.TypeOf(in).Kind()
	if kind != reflect.Ptr {
		return errors.Errorf("first must be a pointer to struct, %s is not valid", kind.String())
	}

	for field, value := range fields {
		err := SetValueOfStruct(in, field, value)
		if err != nil {
			return err
		}
	}

	return nil
}

// SetValueOfStructByTag sets value of specified field of input struct,
// field in this function represents the tag of the struct field,
// the concerning struct field must exist and be exported, otherwise, it will return an error,
// the first argument must be a pointer to struct
// if value is nil, the field value will be set to ZERO value of the type
func SetValueOfStructByTag(in interface{}, field string, value interface{}, tag string) error {
	kind := reflect.TypeOf(in).Kind()
	if kind != reflect.Ptr {
		return errors.Errorf("first must be a pointer to struct, %s is not valid", kind.String())
	}

	if tag == constant.EmptyString {
		return errors.New("tag should not be empty")
	}

	inVal := reflect.ValueOf(in).Elem()

	for i := constant.ZeroInt; i < inVal.NumField(); i++ {
		fieldType := inVal.Type().Field(i)
		fieldTag := fieldType.Tag.Get(tag)
		if fieldTag == field {
			v := inVal.FieldByName(fieldType.Name)
			if !v.IsValid() {
				return errors.Errorf("field does not exist. field: %s", field)
			}
			if !v.CanSet() {
				return errors.Errorf("field can not be set, please check if this field is exported. filed: %s", field)
			}

			vType := v.Type()
			valueType := reflect.TypeOf(value)

			if valueType == nil {
				// set zero value
				v.Set(reflect.Zero(vType))
				return nil
			}

			if vType != valueType {
				if valueType.ConvertibleTo(vType) {
					// convert value type to field type
					v.Set(reflect.ValueOf(value).Convert(vType))
					return nil
				}

				return errors.Errorf("types of field and value mismatched. field: %s, field type: %s, value type: %s",
					field, v.Type().String(), valueType.String())
			}

			// set value
			v.Set(reflect.ValueOf(value))

			return nil
		}
	}

	return errors.Errorf("field does not exist in the struct with given tag. field: %s, tag: %s", field, tag)
}

// SetValuesWithMapByTag sets values of input struct with given map,
// the fields of map represents the tag of the struct field,
// the concerning struct field must exist and be exported, otherwise, it will return an error,
// the first argument must be a pointer to struct
func SetValuesWithMapByTag(in interface{}, fields map[string]interface{}, tag string) error {
	kind := reflect.TypeOf(in).Kind()
	if kind != reflect.Ptr {
		return errors.Errorf("first must be a pointer to struct, %s is not valid", kind.String())
	}

	for field, value := range fields {
		err := SetValueOfStructByTag(in, field, value, tag)
		if err != nil {
			return err
		}
	}

	return nil
}

// SetValuesWithMapAndRandom sets values of input struct with given map,
// if fields in struct does not exist in given map, some of them--depends on the data type--will be set with default value,
// the fields of map must exist and be exported, otherwise, it will return an error,
// the first argument must be a pointer to struct
func SetValuesWithMapAndRandom(in interface{}, fields map[string]interface{}) error {
	kind := reflect.TypeOf(in).Kind()
	if kind != reflect.Ptr {
		return errors.Errorf("first must be a pointer to struct, %s is not valid", kind.String())
	}

	inVal := reflect.ValueOf(in).Elem()
	inType := inVal.Type()
	for i := constant.ZeroInt; i < inVal.NumField(); i++ {
		fieldName := inType.Field(i).Name
		fieldValue, exists := fields[fieldName]
		if !exists {
			fieldValue = inVal.Field(i).Interface()
			// set default value
			switch fieldValue.(type) {
			case int, int32, int64, uint, uint32, uint64:
				fieldValue = constant.DefaultRandomInt
			case float32, float64:
				fieldValue = constant.DefaultRandomFloat
			case string:
				fieldValue = constant.DefaultRandomString
			case time.Time:
				fieldValue = constant.DefaultRandomTime
			default:
				// TODO: for now, do nothing here
				continue
			}
		}

		err := SetValueOfStruct(in, fieldName, fieldValue)
		if err != nil {
			return err
		}
	}

	return nil
}

// SetValuesWithMapAndRandomByTag sets values of input struct with given map,
// field in this function represents the tag of the struct field,
// if fields in struct does not exist in given map, some of them--depends on the data type--will be set with default value,
// the fields of map must exist and be exported, otherwise, it will return an error,
// the first argument must be a pointer to struct
func SetValuesWithMapAndRandomByTag(in interface{}, fields map[string]interface{}, tag string) error {
	kind := reflect.TypeOf(in).Kind()
	if kind != reflect.Ptr {
		return errors.Errorf("first must be a pointer to struct, %s is not valid", kind.String())
	}

	inVal := reflect.ValueOf(in).Elem()
	inType := inVal.Type()
	for i := constant.ZeroInt; i < inVal.NumField(); i++ {
		fieldName := inType.Field(i).Name
		fieldTag := inType.Field(i).Tag.Get(tag)
		fieldValue, exists := fields[fieldTag]
		if !exists {
			fieldValue = inVal.Field(i).Interface()
			// set default value
			switch fieldValue.(type) {
			case int, int32, int64, uint, uint32, uint64:
				fieldValue = constant.DefaultRandomInt
			case float32, float64:
				fieldValue = constant.DefaultRandomFloat
			case string:
				fieldValue = constant.DefaultRandomString
			case time.Time:
				fieldValue = constant.DefaultRandomTime
			default:
				// TODO: for now, do nothing here
				continue
			}
		}

		err := SetValueOfStruct(in, fieldName, fieldValue)
		if err != nil {
			return err
		}
	}

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
		v, err := ConvertToInt(value)
		if err != nil {
			return err
		}

		switch kind {
		case reflect.Int:
			err = SetValueOfStruct(in, field, v)
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
			err = SetValueOfStruct(in, field, uint64(v))
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
			return errors.Errorf("field %s does not exist", field)
		}

		k := fieldType.Type.Elem().Kind()
		b, ok := value.([]byte)
		if ok && json.Valid(b) {
			v, err := ConvertBytesToSlice(b, k)
			if err != nil {
				return err
			}

			err = SetValueOfStruct(in, field, v)
			if err != nil {
				return err
			}

			return nil
		}

		v, err := ConvertToSlice(value, k)
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
		if v == constant.EmptyString {
			return SetValueOfStruct(in, field, time.Time{})
		}
		t, err := time.ParseInLocation(constant.DefaultTimeLayout, v, time.Local)
		if err != nil {
			return err
		}
		err = SetValueOfStruct(in, field, t)
		if err != nil {
			return err
		}
	default:
		return errors.Errorf("unsupported data type: %s", kind.String())
	}

	return nil
}

// CopyStructWithFields returns a new struct with only specified fields
// NOTE:
//  1. tags and values of fields are exactly same
//  2. only exported and addressable fields will be copied
//  3. if any field in fields does not exist in the input struct, it returns error
//  4. if values in input struct is a pointer, then value in the new struct will point to the same object
//  5. returning struct is totally a new data type, so you could not use any (*type) assertion
//  6. if fields argument is empty, a new struct which contains the whole fields of input struct will be returned
//  7. technically, for convenience purpose, this function creates a new struct as same as input struct,
//     then removes fields that do not exist in the given fields
func CopyStructWithFields(in interface{}, fields ...string) (interface{}, error) {
	if len(fields) == constant.ZeroInt {
		return CopyStructWithoutFields(in)
	}

	kind := reflect.TypeOf(in).Kind()
	if kind != reflect.Ptr {
		return nil, errors.Errorf("first must be a pointer to struct, %s is not valid", kind.String())
	}

	var removeFields []string

	inVal := reflect.ValueOf(in).Elem()
	inType := inVal.Type()

	for i := constant.ZeroInt; i < inVal.NumField(); i++ {
		fieldName := inType.Field(i).Name
		ok := ElementInSlice(fields, fieldName)
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
		return nil, errors.Trace(err)
	}

	for _, field := range fields {
		newStruct = newStruct.RemoveField(field)
	}

	// generate new instance
	newInstance := newStruct.Build().New()
	newVal := reflect.ValueOf(newInstance).Elem()
	newType := newVal.Type()

	inVal := reflect.ValueOf(in).Elem()

	for i := constant.ZeroInt; i < newVal.NumField(); i++ {
		fieldType := newType.Field(i)
		fieldVal := newVal.Field(i)
		// set value
		newField := inVal.FieldByName(fieldType.Name)
		if newField.Type().Kind() == reflect.Interface {
			newField = reflect.New(newField.Elem().Type())
			newField.Elem().Set(inVal.FieldByName(fieldType.Name).Elem())
			newField = newField.Elem()
		}

		fieldVal.Set(newField)
	}

	return newInstance, nil
}

// MarshalStructWithFields marshals input struct using json.Marshal() with given fields,
// first argument must be a pointer to struct, not the struct itself
func MarshalStructWithFields(in interface{}, fields ...string) ([]byte, error) {
	return marshalStructWithFunc(in, CopyStructWithFields, fields...)
}

// MarshalStructWithoutFields marshals input struct using json.Marshal() without given fields,
// first argument must be a pointer to struct, not the struct itself
func MarshalStructWithoutFields(in interface{}, fields ...string) ([]byte, error) {
	return marshalStructWithFunc(in, CopyStructWithoutFields, fields...)
}

func marshalStructWithFunc(in interface{}, copyFunc func(interface{}, ...string) (interface{}, error), fields ...string) ([]byte, error) {
	kind := reflect.TypeOf(in).Kind()
	if kind != reflect.Ptr {
		return nil, errors.Errorf("first must be a pointer to struct, %s is not valid", kind.String())
	}

	if len(fields) == constant.ZeroInt {
		bytes, err := json.Marshal(in)

		return bytes, errors.Trace(err)
	}

	// generate a new struct with fields
	newInstance, err := copyFunc(in, fields...)
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(newInstance)

	return bytes, errors.Trace(err)
}

// MarshalStructWithTag marshals input struct using json.Marshal() with fields that contain given tag,
// first argument must be a pointer to struct, not the struct itself
func MarshalStructWithTag(in interface{}, tag string) ([]byte, error) {
	kind := reflect.TypeOf(in).Kind()
	if kind != reflect.Ptr {
		return nil, errors.Errorf("first must be a pointer to struct, %s is not valid", kind.String())
	}

	if tag == constant.EmptyString {
		return nil, errors.New("tag should not be empty")
	}

	inVal := reflect.ValueOf(in).Elem()

	var fields []string

	for i := constant.ZeroInt; i < inVal.NumField(); i++ {
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
// if key and tag are same, the field of the input struct will be the key of the new map,
// the value of the given map will be the value of the new map
func NewMapWithStructTag(m map[string]interface{}, in interface{}, tag string) (map[string]interface{}, error) {
	kind := reflect.TypeOf(in).Kind()
	if kind != reflect.Ptr {
		return nil, errors.Errorf("second must be a pointer to struct, %s is not valid", kind.String())
	}

	result := make(map[string]interface{})
	inVal := reflect.ValueOf(in).Elem()

Loop:
	for key := range m {
		for i := constant.ZeroInt; i < inVal.NumField(); i++ {
			fieldType := inVal.Type().Field(i)
			fieldTag := fieldType.Tag.Get(tag)

			if key == fieldTag {
				result[fieldType.Name] = m[key]
				continue Loop
			}
		}
	}
	return result, nil
}

// UnmarshalToMapWithStructTag returns a map as the result, it works as following logic:
// 1. unmarshals given data to a temporary map to get keys
// 2. unmarshals given data to the input struct, to get field names and values with appropriate data types
// 3. loop keys in the temporary map, loops tags of the input struct in a nested loop
// 4. if the key matches the tag, set field name as the key of result map, set field value as the value of the result map
func UnmarshalToMapWithStructTag(data []byte, in interface{}, tag string) (map[string]interface{}, error) {
	kind := reflect.TypeOf(in).Kind()
	if kind != reflect.Ptr {
		return nil, errors.Errorf("second argument must be a pointer to struct, %s is not valid", kind.String())
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
		return nil, errors.Trace(err)
	}
	// unmarshal to temporary map to get key names
	tmpMap := make(map[string]interface{})
	err = decoder.Unmarshal(data, &tmpMap)
	if err != nil {
		return nil, errors.Trace(err)
	}

	result := make(map[string]interface{})
	inVal := reflect.ValueOf(in).Elem()
Loop:
	for key := range tmpMap {
		for i := constant.ZeroInt; i < inVal.NumField(); i++ {
			fieldType := inVal.Type().Field(i)
			fieldTag := fieldType.Tag.Get(tag)
			if key == fieldTag {
				// set field name of struct as key, set value of temporary map as value of the result
				result[fieldType.Name] = inVal.Field(i).Interface()
				continue Loop
			}
		}
	}

	return result, nil
}

// UnmarshalToMapWithStructTagFromString converts given string to []byte and then call UnmarshalToMapWithStructTag() function
func UnmarshalToMapWithStructTagFromString(s string, in interface{}, tag string) (map[string]interface{}, error) {
	return UnmarshalToMapWithStructTag([]byte(s), in, tag)
}
