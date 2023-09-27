package common

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/constant"
)

const (
	testWeakPassNotLongEnough = "123456"
	testWeakPassTooLong       = "Hello,World.123Hello,World.123"
	testWeakPassNoUpperCase   = "hello,world.123"
	testWeakPassNoLowerCase   = "HELLO,WORLD.123"
	testWeakPassNoDigit       = "Hello,World."
	testWeakPassNoSpecial     = "HelloWorld123"
	testWeakPassUnsupported   = "Hello,World中文"

	testStrongPass = "Hello,World.123"
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
	asst.Equal(DefaultNormalCharNum+constant.OneInt, NormalCharNum, "test GetRandomStringWithDefault() failed")
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

func TestCheckPasswordStrength(t *testing.T) {
	asst := assert.New(t)

	err := CheckPasswordStrength(testWeakPassNotLongEnough, DefaultMinPasswordLength, DefaultMaxPasswordLength, DefaultMinUpperCaseCharNum,
		DefaultMinLowerCaseCharNum, DefaultMinDigitalCharNum, DefaultMinSpecialCharNum, DefaultSpecialCharString)
	asst.NotNil(err, "test CheckPasswordStrength() failed")
	t.Logf("check error: %s", err.Error())
	err = CheckPasswordStrength(testWeakPassTooLong, DefaultMinPasswordLength, DefaultMaxPasswordLength, DefaultMinUpperCaseCharNum,
		DefaultMinLowerCaseCharNum, DefaultMinDigitalCharNum, DefaultMinSpecialCharNum, DefaultSpecialCharString)
	asst.NotNil(err, "test CheckPasswordStrength() failed")
	t.Logf("check error: %s", err.Error())
	err = CheckPasswordStrength(testWeakPassNoUpperCase, DefaultMinPasswordLength, DefaultMaxPasswordLength, DefaultMinUpperCaseCharNum,
		DefaultMinLowerCaseCharNum, DefaultMinDigitalCharNum, DefaultMinSpecialCharNum, DefaultSpecialCharString)
	asst.NotNil(err, "test CheckPasswordStrength() failed")
	t.Logf("check error: %s", err.Error())
	err = CheckPasswordStrength(testWeakPassNoLowerCase, DefaultMinPasswordLength, DefaultMaxPasswordLength, DefaultMinUpperCaseCharNum,
		DefaultMinLowerCaseCharNum, DefaultMinDigitalCharNum, DefaultMinSpecialCharNum, DefaultSpecialCharString)
	asst.NotNil(err, "test CheckPasswordStrength() failed")
	t.Logf("check error: %s", err.Error())
	err = CheckPasswordStrength(testWeakPassNoDigit, DefaultMinPasswordLength, DefaultMaxPasswordLength, DefaultMinUpperCaseCharNum,
		DefaultMinLowerCaseCharNum, DefaultMinDigitalCharNum, DefaultMinSpecialCharNum, DefaultSpecialCharString)
	asst.NotNil(err, "test CheckPasswordStrength() failed")
	t.Logf("check error: %s", err.Error())
	err = CheckPasswordStrength(testWeakPassNoSpecial, DefaultMinPasswordLength, DefaultMaxPasswordLength, DefaultMinUpperCaseCharNum,
		DefaultMinLowerCaseCharNum, DefaultMinDigitalCharNum, DefaultMinSpecialCharNum, DefaultSpecialCharString)
	asst.NotNil(err, "test CheckPasswordStrength() failed")
	t.Logf("check error: %s", err.Error())
	err = CheckPasswordStrength(testWeakPassUnsupported, DefaultMinPasswordLength, DefaultMaxPasswordLength, DefaultMinUpperCaseCharNum,
		DefaultMinLowerCaseCharNum, DefaultMinDigitalCharNum, DefaultMinSpecialCharNum, DefaultSpecialCharString)
	asst.NotNil(err, "test CheckPasswordStrength() failed")
	t.Logf("check error: %s", err.Error())

	err = CheckPasswordStrength(testStrongPass, DefaultMinPasswordLength, DefaultMaxPasswordLength, DefaultMinUpperCaseCharNum,
		DefaultMinLowerCaseCharNum, DefaultMinDigitalCharNum, DefaultMinSpecialCharNum, DefaultSpecialCharString)
	asst.Nil(err, "test CheckPasswordStrength() failed")
}

func TestCheckPasswordStrengthWithDefault(t *testing.T) {
	asst := assert.New(t)

	err := CheckPasswordStrengthWithDefault(testWeakPassNotLongEnough)
	asst.NotNil(err, "test CheckPasswordStrengthWithDefault() failed")
	t.Logf("check error: %s", err.Error())
	err = CheckPasswordStrengthWithDefault(testWeakPassTooLong)
	asst.NotNil(err, "test CheckPasswordStrengthWithDefault() failed")
	t.Logf("check error: %s", err.Error())
	err = CheckPasswordStrengthWithDefault(testWeakPassNoUpperCase)
	asst.NotNil(err, "test CheckPasswordStrengthWithDefault() failed")
	t.Logf("check error: %s", err.Error())
	err = CheckPasswordStrengthWithDefault(testWeakPassNoLowerCase)
	asst.NotNil(err, "test CheckPasswordStrengthWithDefault() failed")
	t.Logf("check error: %s", err.Error())
	err = CheckPasswordStrengthWithDefault(testWeakPassNoDigit)
	asst.NotNil(err, "test CheckPasswordStrengthWithDefault() failed")
	t.Logf("check error: %s", err.Error())
	err = CheckPasswordStrengthWithDefault(testWeakPassNoSpecial)
	asst.NotNil(err, "test CheckPasswordStrengthWithDefault() failed")
	t.Logf("check error: %s", err.Error())
	err = CheckPasswordStrengthWithDefault(testWeakPassUnsupported)
	asst.NotNil(err, "test CheckPasswordStrengthWithDefault() failed")
	t.Logf("check error: %s", err.Error())

	err = CheckPasswordStrengthWithDefault(testStrongPass)
	asst.Nil(err, "test CheckPasswordStrengthWithDefault() failed")
}
