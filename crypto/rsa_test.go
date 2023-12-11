package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	defaultMessage = "hello, world!12345678"
)

func createRSA() (*RSA, error) {
	return NewRSA()
}

func TestCrypt_All(t *testing.T) {
	TestRSA_Encrypt(t)
	TestRSA_Decrypt(t)
	TestRSA_EncryptWithPublicKey(t)
	TestRSA_DecryptWithPrivateKey(t)
	TestRSA_EncryptWithPrivateKey(t)
	TestRSA_DecryptWithPublicKey(t)
}

func TestRSA_Encrypt(t *testing.T) {
	asst := assert.New(t)

	r, err := createRSA()
	asst.Nil(err, "test RSA.Encrypt() failed")
	privateKeyString, err := r.GetPrivateKey()
	asst.Nil(err, "test RSA.Encrypt() failed")
	publicKeyString, err := r.GetPublicKey()
	asst.Nil(err, "test RSA.Encrypt() failed")

	t.Logf("private key: %s", privateKeyString)
	t.Logf("public key: %s", publicKeyString)

	cipherText, err := r.EncryptWithPublicKey(defaultMessage)
	asst.Nil(err, "test RSA.Encrypt() failed")
	message, err := r.DecryptWithPrivateKey(cipherText)
	asst.Nil(err, "test RSA.Encrypt() failed")
	asst.Equal(defaultMessage, message, "test RSA.Encrypt() failed")

	r, err = NewRSAWithPrivateKeyString(privateKeyString)
	asst.Nil(err, "test RSA.Encrypt() failed")
	message, err = r.DecryptWithPrivateKey(cipherText)
	asst.Nil(err, "test RSA.Encrypt() failed")
	asst.Equal(defaultMessage, message, "test RSA.Encrypt() failed")
}

func TestRSA_Decrypt(t *testing.T) {
	asst := assert.New(t)

	r, err := createRSA()
	asst.Nil(err, "test RSA.Decrypt() failed")
	privateKeyString, err := r.GetPrivateKey()
	asst.Nil(err, "test RSA.Decrypt() failed")
	publicKeyString, err := r.GetPublicKey()
	asst.Nil(err, "test RSA.Decrypt() failed")

	t.Logf("private key: %s", privateKeyString)
	t.Logf("public key: %s", publicKeyString)

	cipherText, err := r.EncryptWithPublicKey(defaultMessage)
	asst.Nil(err, "test RSA.Decrypt() failed")
	message, err := r.DecryptWithPrivateKey(cipherText)
	asst.Nil(err, "test RSA.Decrypt() failed")
	asst.Equal(defaultMessage, message, "test RSA.Decrypt() failed")

	r, err = NewRSAWithPrivateKeyString(privateKeyString)
	asst.Nil(err, "test RSA.Decrypt() failed")
	message, err = r.DecryptWithPrivateKey(cipherText)
	asst.Nil(err, "test RSA.Decrypt() failed")
	asst.Equal(defaultMessage, message, "test RSA.Decrypt() failed")
}

func TestRSA_EncryptWithPublicKey(t *testing.T) {
	asst := assert.New(t)

	r, err := createRSA()
	asst.Nil(err, "test RSA.EncryptWithPublicKey() failed")
	privateKeyString, err := r.GetPrivateKey()
	asst.Nil(err, "test RSA.EncryptWithPublicKey() failed")
	publicKeyString, err := r.GetPublicKey()
	asst.Nil(err, "test RSA.EncryptWithPublicKey() failed")

	t.Logf("private key: %s", privateKeyString)
	t.Logf("public key: %s", publicKeyString)

	cipherText, err := r.EncryptWithPublicKey(defaultMessage)
	asst.Nil(err, "test RSA.EncryptWithPublicKey() failed")
	message, err := r.DecryptWithPrivateKey(cipherText)
	asst.Nil(err, "test RSA.EncryptWithPublicKey() failed")
	asst.Equal(defaultMessage, message, "test RSA.EncryptWithPublicKey() failed")

	r, err = NewRSAWithPrivateKeyString(privateKeyString)
	asst.Nil(err, "test RSA.EncryptWithPublicKey() failed")
	message, err = r.DecryptWithPrivateKey(cipherText)
	asst.Nil(err, "test RSA.EncryptWithPublicKey() failed")
	asst.Equal(defaultMessage, message, "test RSA.EncryptWithPublicKey() failed")
}

func TestRSA_DecryptWithPrivateKey(t *testing.T) {
	asst := assert.New(t)

	r, err := createRSA()
	asst.Nil(err, "test RSA.DecryptWithPrivateKey() failed")
	privateKeyString, err := r.GetPrivateKey()
	asst.Nil(err, "test RSA.DecryptWithPrivateKey() failed")
	publicKeyString, err := r.GetPublicKey()
	asst.Nil(err, "test RSA.DecryptWithPrivateKey() failed")

	t.Logf("private key: %s", privateKeyString)
	t.Logf("public key: %s", publicKeyString)

	cipherText, err := r.EncryptWithPublicKey(defaultMessage)
	asst.Nil(err, "test RSA.DecryptWithPrivateKey() failed")
	message, err := r.DecryptWithPrivateKey(cipherText)
	asst.Nil(err, "test RSA.DecryptWithPrivateKey() failed")
	asst.Equal(defaultMessage, message, "test RSA.DecryptWithPrivateKey() failed")

	r, err = NewRSAWithPrivateKeyString(privateKeyString)
	asst.Nil(err, "test RSA.DecryptWithPrivateKey() failed")
	message, err = r.DecryptWithPrivateKey(cipherText)
	asst.Nil(err, "test RSA.DecryptWithPrivateKey() failed")
	asst.Equal(defaultMessage, message, "test RSA.DecryptWithPrivateKey() failed")
}

func TestRSA_EncryptWithPrivateKey(t *testing.T) {
	asst := assert.New(t)

	r, err := createRSA()
	asst.Nil(err, "test RSA.EncryptWithPrivateKey() failed")
	privateKeyString, err := r.GetPrivateKey()
	asst.Nil(err, "test RSA.EncryptWithPrivateKey() failed")
	publicKeyString, err := r.GetPublicKey()
	asst.Nil(err, "test RSA.EncryptWithPrivateKey() failed")

	t.Logf("private key: %s", privateKeyString)
	t.Logf("public key: %s", publicKeyString)

	cipherText, err := r.EncryptWithPrivateKey(defaultMessage)
	asst.Nil(err, "test RSA.EncryptWithPrivateKey() failed")
	message, err := r.DecryptWithPublicKey(cipherText)
	asst.Nil(err, "test RSA.EncryptWithPrivateKey() failed")
	asst.Equal(defaultMessage, message, "test RSA.EncryptWithPrivateKey() failed")

	r, err = NewRSAWithPrivateKeyString(privateKeyString)
	asst.Nil(err, "test RSA.EncryptWithPrivateKey() failed")
	message, err = r.DecryptWithPublicKey(cipherText)
	asst.Nil(err, "test RSA.EncryptWithPrivateKey() failed")
	asst.Equal(defaultMessage, message, "test RSA.EncryptWithPrivateKey() failed")
}

func TestRSA_DecryptWithPublicKey(t *testing.T) {
	asst := assert.New(t)

	r, err := createRSA()
	asst.Nil(err, "test RSA.DecryptWithPublicKey() failed")
	privateKeyString, err := r.GetPrivateKey()
	asst.Nil(err, "test RSA.DecryptWithPublicKey() failed")
	publicKeyString, err := r.GetPublicKey()
	asst.Nil(err, "test RSA.DecryptWithPublicKey() failed")

	t.Logf("private key: %s", privateKeyString)
	t.Logf("public key: %s", publicKeyString)

	cipherText, err := r.EncryptWithPrivateKey(defaultMessage)
	asst.Nil(err, "test RSA.DecryptWithPublicKey() failed")
	message, err := r.DecryptWithPublicKey(cipherText)
	asst.Nil(err, "test RSA.DecryptWithPublicKey() failed")
	asst.Equal(defaultMessage, message, "test RSA.DecryptWithPublicKey() failed")

	r, err = NewRSAWithPrivateKeyString(privateKeyString)
	asst.Nil(err, "test RSA.DecryptWithPublicKey() failed")
	message, err = r.DecryptWithPublicKey(cipherText)
	asst.Nil(err, "test RSA.DecryptWithPublicKey() failed")
	asst.Equal(defaultMessage, message, "test RSA.DecryptWithPublicKey() failed")
}

func Test_Temp(t *testing.T) {
	asst := assert.New(t)

	r, err := newRSAWithKeySize(256)
	asst.Nil(err, "test Temp failed")
	privateKeyString, err := r.GetPrivateKey()
	asst.Nil(err, "test Temp failed")
	publicKeyString, err := r.GetPublicKey()
	asst.Nil(err, "test Temp failed")

	t.Logf("private key: %s", privateKeyString)
	t.Logf("public key: %s", publicKeyString)
}
