package common

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/constant"
)

func TestConvert_All(t *testing.T) {
	TestConver_StringToBytes(t)
	TestConvert_BytesToString(t)
	TestConvert_ConvertSliceToString(t)
	TestConvert_ConvertInterfaceSliceToString(t)
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
