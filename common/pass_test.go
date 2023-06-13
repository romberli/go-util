package common

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/constant"
)

func checkRandomString(s string) (int, int, int, int) {
	var (
		NormalCharNum  int
		DigitCharNum   int
		SpecialCharNum int
		UnknownCharNum int
	)
	for _, c := range s {
		var isSpecialChar bool

		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
			NormalCharNum++
			continue
		}
		if c >= '0' && c <= '9' {
			DigitCharNum++
			continue
		}
		for _, sc := range DefaultSpecialCharString {
			if c == sc {
				SpecialCharNum++
				isSpecialChar = true
			}
		}

		if !isSpecialChar {
			UnknownCharNum++
		}
	}

	return NormalCharNum, DigitCharNum, SpecialCharNum, UnknownCharNum
}

func TestPass_All(t *testing.T) {
	TestGetRandomStringWithDefault(t)
	TestGetRandomNormalCharString(t)
	TestGetRandomDigitalString(t)
	TestGetRandomSpecialString(t)
	TestGetRandomString(t)
	TestGetRandomBytes(t)
}

func TestGetRandomStringWithDefault(t *testing.T) {
	asst := assert.New(t)

	s := GetRandomStringWithDefault()
	NormalCharNum, DigitCharNum, SpecialCharNum, UnknownCharNum := checkRandomString(s)
	asst.Equal(DefaultNormalCharNum, NormalCharNum, "test GetRandomStringWithDefault() failed")
	asst.Equal(DefaultDigitalNum, DigitCharNum, "test GetRandomStringWithDefault() failed")
	asst.Equal(DefaultSpecialNum, SpecialCharNum, "test GetRandomStringWithDefault() failed")
	asst.Equal(constant.ZeroInt, UnknownCharNum, "test GetRandomStringWithDefault() failed")
}

func TestGetRandomNormalCharString(t *testing.T) {
	asst := assert.New(t)

	s := GetRandomNormalCharString(DefaultNormalCharNum)
	NormalCharNum, DigitCharNum, SpecialCharNum, UnknownCharNum := checkRandomString(s)
	asst.Equal(DefaultNormalCharNum, NormalCharNum, "test GetRandomNormalCharString() failed")
	asst.Equal(constant.ZeroInt, DigitCharNum, "test GetRandomNormalCharString() failed")
	asst.Equal(constant.ZeroInt, SpecialCharNum, "test GetRandomNormalCharString() failed")
	asst.Equal(constant.ZeroInt, UnknownCharNum, "test GetRandomNormalCharString() failed")
}

func TestGetRandomDigitalString(t *testing.T) {
	asst := assert.New(t)

	s := GetRandomDigitalString(DefaultDigitalNum)
	NormalCharNum, DigitCharNum, SpecialCharNum, UnknownCharNum := checkRandomString(s)
	asst.Equal(constant.ZeroInt, NormalCharNum, "test GetRandomDigitalString() failed")
	asst.Equal(DefaultDigitalNum, DigitCharNum, "test GetRandomDigitalString() failed")
	asst.Equal(constant.ZeroInt, SpecialCharNum, "test GetRandomDigitalString() failed")
	asst.Equal(constant.ZeroInt, UnknownCharNum, "test GetRandomDigitalString() failed")
}

func TestGetRandomSpecialString(t *testing.T) {
	asst := assert.New(t)

	s := GetRandomSpecialString(DefaultSpecialNum)
	NormalCharNum, DigitCharNum, SpecialCharNum, UnknownCharNum := checkRandomString(s)
	asst.Equal(constant.ZeroInt, NormalCharNum, "test GetRandomSpecialString() failed")
	asst.Equal(constant.ZeroInt, DigitCharNum, "test GetRandomSpecialString() failed")
	asst.Equal(DefaultSpecialNum, SpecialCharNum, "test GetRandomSpecialString() failed")
	asst.Equal(constant.ZeroInt, UnknownCharNum, "test GetRandomSpecialString() failed")
}

func TestGetRandomString(t *testing.T) {
	asst := assert.New(t)

	s := GetRandomString(DefaultNormalCharString+DefaultDigitalCharString+DefaultSpecialCharString, DefaultNormalCharNum+DefaultDigitalNum+DefaultSpecialNum)
	NormalCharNum, DigitCharNum, SpecialCharNum, UnknownCharNum := checkRandomString(s)
	asst.NotEmptyf(s, "test GetRandomString() failed")

	t.Logf("randomString: %s, normalCharNum: %d, digitCharNum: %d, specialCharNum: %d, unknownCharNum: %d",
		s, NormalCharNum, DigitCharNum, SpecialCharNum, UnknownCharNum)
}

func TestGetRandomBytes(t *testing.T) {
	asst := assert.New(t)

	b := GetRandomBytes(DefaultNormalCharString+DefaultDigitalCharString+DefaultSpecialCharString, DefaultNormalCharNum+DefaultDigitalNum+DefaultSpecialNum)
	NormalCharNum, DigitCharNum, SpecialCharNum, UnknownCharNum := checkRandomString(string(b))
	asst.NotEmptyf(string(b), "test GetRandomBytes() failed")

	t.Logf("randomBytes: %s, normalCharNum: %d, digitCharNum: %d, specialCharNum: %d, unknownCharNum: %d",
		string(b), NormalCharNum, DigitCharNum, SpecialCharNum, UnknownCharNum)
}
