package common

import (
	"math/rand"
	"time"

	"github.com/romberli/go-util/constant"
)

const (
	DefaultLetters           = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+-,.?_"
	DefaultNormalCharString  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	DefaultDigitalCharString = "0123456789"
	DefaultSpecialCharString = "+-,.?_"

	DefaultNormalCharNum = 6
	DefaultDigitalNum    = 3
	DefaultSpecialNum    = 2
)

// GetRandomStringWithDefault returns a random string with default length
func GetRandomStringWithDefault() string {
	normalCharBytes := GetRandomBytes(DefaultNormalCharString, DefaultNormalCharNum)
	digitalCharBytes := GetRandomBytes(DefaultDigitalCharString, DefaultDigitalNum)
	specialCharBytes := GetRandomBytes(DefaultSpecialCharString, DefaultSpecialNum)

	s := append(normalCharBytes, digitalCharBytes...)
	s = append(s, specialCharBytes...)

	rand.Shuffle(len(s)/constant.TwoInt, func(i, j int) { s[i], s[j] = s[j], s[i] })

	return string(s)
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
