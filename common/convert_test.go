package common

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/constant"
)

func TestConvert_All(t *testing.T) {
	TestConver_StringToBytes(t)
	TestConvert_BytesToString(t)
	TestConvert_ConvertSliceToString(t)
	TestConvert_ConvertInterfaceSliceToString(t)
	TestConvert_JSONNumber(t)
}

func TestConver_StringToBytes(t *testing.T) {
	asst := assert.New(t)

	// empty string
	b := StringToBytes(constant.EmptyString)
	asst.Equal([]byte(nil), b, "test StringToBytes() failed")
}

func TestConvert_BytesToString(t *testing.T) {
	asst := assert.New(t)

	// nil
	s := BytesToString(nil)
	asst.Equal(constant.EmptyString, s, "test BytesToString() failed")
}

func TestConvert_ConvertSliceToString(t *testing.T) {
	asst := assert.New(t)

	// bool slice
	bl := []bool{true, false, true}
	s := "true,false,true"
	asst.Equal(s, ConvertSliceToString(bl, constant.CommaString), "test ConvertSliceToString() failed")
	// string slice
	sl := []string{"a", "b", "c"}
	s = "a,b,c"
	asst.Equal(s, ConvertSliceToString(sl, constant.CommaString), "test ConvertSliceToString() failed")
	// int slice
	il := []int{1, 2, 3}
	s = "1,2,3"
	asst.Equal(s, ConvertSliceToString(il, constant.CommaString), "test ConvertSliceToString() failed")
	// float64 slice
	fl := []float64{1.1, 2.2, 3.3}
	s = "1.1,2.2,3.3"
	asst.Equal(s, ConvertSliceToString(fl, constant.CommaString), "test ConvertSliceToString() failed")
}

func TestConvert_ConvertInterfaceSliceToString(t *testing.T) {
	asst := assert.New(t)

	// bool slice
	bl := []interface{}{true, false, true}
	s := "true,false,true"
	asst.Equal(s, ConvertInterfaceSliceToString(bl, constant.CommaString), "test ConvertSliceToString() failed")
	// string slice
	sl := []interface{}{"a", "b", "c"}
	s = "a,b,c"
	asst.Equal(s, ConvertInterfaceSliceToString(sl, constant.CommaString), "test ConvertSliceToString() failed")
	// int slice
	il := []interface{}{1, 2, 3}
	s = "1,2,3"
	asst.Equal(s, ConvertInterfaceSliceToString(il, constant.CommaString), "test ConvertSliceToString() failed")
	// float64 slice
	fl := []interface{}{1.1, 2.2, 3.3}
	s = "1.1,2.2,3.3"
	asst.Equal(s, ConvertInterfaceSliceToString(fl, constant.CommaString), "test ConvertSliceToString() failed")

}

func TestConvert_JSONNumber(t *testing.T) {
	asst := assert.New(t)

	// large integer keeps precision
	v, err := ConvertToInt64(json.Number("9007199254740993"))
	asst.Nil(err, "test ConvertToInt64() failed")
	asst.Equal(int64(9007199254740993), v, "test ConvertToInt64() failed")

	// integer overflow returns error instead of silent wrapping
	_, err = ConvertToInt64(json.Number("18446744073709551616"))
	asst.NotNil(err, "test ConvertToInt64() failed")

	f, err := ConvertToFloat(json.Number("1.5"))
	asst.Nil(err, "test ConvertToFloat() failed")
	asst.Equal(1.5, f, "test ConvertToFloat() failed")

	// string keeps the original literal
	s, err := ConvertToString(json.Number("1234567890123456789"))
	asst.Nil(err, "test ConvertToString() failed")
	asst.Equal("1234567890123456789", s, "test ConvertToString() failed")

	s, err = ConvertNumberToString(json.Number("9007199254740993"))
	asst.Nil(err, "test ConvertNumberToString() failed")
	asst.Equal("9007199254740993", s, "test ConvertNumberToString() failed")

	b, err := ConvertToBool(json.Number("1"))
	asst.Nil(err, "test ConvertToBool() failed")
	asst.True(b, "test ConvertToBool() failed")
}
