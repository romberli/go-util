package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func createAES() (*AES, error) {
	return NewAES()
}

func TestAES_All(t *testing.T) {
	TestAES_Encrypt(t)
	TestAES_Decrypt(t)
}

func TestAES_Encrypt(t *testing.T) {
	asst := assert.New(t)

	a, err := createAES()
	asst.Nil(err, "test AES.Encrypt() failed")

	cipherText, err := a.Encrypt(defaultMessage)
	asst.Nil(err, "test AES.Encrypt() failed")
	message, err := a.Decrypt(cipherText)
	asst.Nil(err, "test AES.Encrypt() failed")
	asst.Equal(defaultMessage, message, "test AES.Encrypt() failed")
}

func TestAES_Decrypt(t *testing.T) {
	asst := assert.New(t)

	a, err := createAES()
	asst.Nil(err, "test AES.Decrypt() failed")

	cipherText, err := a.Encrypt(defaultMessage)
	asst.Nil(err, "test AES.Decrypt() failed")
	message, err := a.Decrypt(cipherText)
	asst.Nil(err, "test AES.Decrypt() failed")
	asst.Equal(defaultMessage, message, "test AES.Decrypt() failed")
}

func TestAES_Temp(t *testing.T) {
	asst := assert.New(t)

	key := "abcdefghijklmnopqrstuvwxyz123456"
	message := defaultMessage
	a := newAESWithKey(key)

	cipherText, err := a.Encrypt(message)
	asst.Nil(err, "test AES.Temp() failed")
	plainText, err := a.Decrypt(cipherText)
	asst.Nil(err, "test AES.Temp() failed")
	asst.Equal(message, plainText, "test AES.Temp() failed")
}
