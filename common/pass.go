package common

import (
	"math/rand"
	"time"

	"github.com/pingcap/errors"

	"github.com/romberli/go-util/constant"
)

const (
	DefaultLetters             = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+-,.?_"
	DefaultNormalCharString    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	DefaultUpperCaseCharString = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	DefaultLowerCaseCharString = "abcdefghijklmnopqrstuvwxyz"
	DefaultDigitalCharString   = "0123456789"
	DefaultSpecialCharString   = "@#$%-,._"

	DefaultNormalCharNum    = 6
	DefaultUpperCaseCharNum = 2
	DefaultLowerCaseCharNum = 4
	DefaultDigitalNum       = 3
	DefaultSpecialNum       = 2

	DefaultMinPasswordLength   = 8
	DefaultMaxPasswordLength   = 20
	DefaultMinUpperCaseCharNum = 1
	DefaultMinLowerCaseCharNum = 1
	DefaultMinDigitalCharNum   = 1
	DefaultMinSpecialCharNum   = 1
)

// GetRandomStringWithDefault returns a random string with default length
func GetRandomStringWithDefault() string {
	first := GetRandomString(DefaultNormalCharString, constant.OneInt)
	upperCaseCharBytes := GetRandomBytes(DefaultUpperCaseCharString, DefaultUpperCaseCharNum)
	lowerCaseCharBytes := GetRandomBytes(DefaultLowerCaseCharString, DefaultLowerCaseCharNum)
	digitalCharBytes := GetRandomBytes(DefaultDigitalCharString, DefaultDigitalNum)
	specialCharBytes := GetRandomBytes(DefaultSpecialCharString, DefaultSpecialNum)

	s := append(upperCaseCharBytes, lowerCaseCharBytes...)
	s = append(s, digitalCharBytes...)
	s = append(s, specialCharBytes...)

	rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })

	return first + string(s)
}

// GetRandomNormalCharString returns a random string with given length
func GetRandomNormalCharString(num int) string {
	return GetRandomString(DefaultNormalCharString, num)
}

// GetRandomDigitalString returns a random string with given length
func GetRandomDigitalString(num int) string {
	return GetRandomString(DefaultDigitalCharString, num)
}

// GetRandomSpecialString returns a random string with given length
func GetRandomSpecialString(num int) string {
	return GetRandomString(DefaultSpecialCharString, num)
}

// GetRandomString returns a random string from given string with given length
func GetRandomString(s string, num int) string {
	return string(GetRandomBytes(s, num))
}

func GetRandomBytes(s string, num int) []byte {
	source := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, num)
	for i := range b {
		b[i] = s[source.Int63()%int64(len(s))]
	}

	return b
}

// CheckPasswordStrength checks the strength of password
func CheckPasswordStrength(password string, minLen, maxLen, minUpperCaseCharNum, minLowerCaseCharNum, minDigitCharNum,
	minSpecialCharNum int, supportedSpecialCharString string) error {
	length := len(password)
	if length < minLen || length > maxLen {
		return errors.Errorf("password length must be between %d and %d, %d is not valid", minLen, maxLen, length)
	}

	var (
		upperCaseCharNum   int
		lowerCaseCharNum   int
		digitCharNum       int
		specialCharNum     int
		unsupportedCharNum int
	)

	for _, c := range password {
		var hasSpecialChar bool

		switch {
		case c >= 'A' && c <= 'Z':
			upperCaseCharNum++
		case c >= 'a' && c <= 'z':
			lowerCaseCharNum++
		case c >= '0' && c <= '9':
			digitCharNum++
		default:
			for _, sc := range supportedSpecialCharString {
				if c == sc {
					specialCharNum++
					hasSpecialChar = true
				}
			}

			if !hasSpecialChar {
				unsupportedCharNum++
			}
		}
	}

	if upperCaseCharNum < minUpperCaseCharNum {
		return errors.Errorf("password must contain at least %d upper case characters, %d is not valid", minUpperCaseCharNum, upperCaseCharNum)
	}
	if lowerCaseCharNum < minLowerCaseCharNum {
		return errors.Errorf("password must contain at least %d lower case characters, %d is not valid", minLowerCaseCharNum, lowerCaseCharNum)
	}
	if digitCharNum < minDigitCharNum {
		return errors.Errorf("password must contain at least %d digit characters, %d is not valid", minDigitCharNum, digitCharNum)
	}
	if specialCharNum < minSpecialCharNum {
		return errors.Errorf("password must contain at least %d special characters, %d is not valid", minSpecialCharNum, specialCharNum)
	}
	if unsupportedCharNum > constant.ZeroInt {
		return errors.Errorf("special character must be one of [%s], the password contains unsupported characters", supportedSpecialCharString)
	}

	return nil
}

// CheckPasswordStrengthWithDefault checks the strength of password with default parameters
func CheckPasswordStrengthWithDefault(password string) error {
	return CheckPasswordStrength(password, DefaultMinPasswordLength, DefaultMaxPasswordLength, DefaultMinUpperCaseCharNum,
		DefaultMinLowerCaseCharNum, DefaultMinDigitalCharNum, DefaultMinSpecialCharNum, DefaultSpecialCharString)
}
