package crypto

import (
	"github.com/romberli/go-util/constant"
)

// EncryptWithRSAPublicKeyString encrypts the data with public key string
func EncryptWithRSAPublicKeyString(publicKeyStr, message string) (string, error) {
	rsa := NewEmptyRSA()
	err := rsa.SetPublicKey(publicKeyStr)
	if err != nil {
		return constant.EmptyString, err
	}

	return rsa.EncryptWithPublicKey(message)
}

// DecryptWithRSAPublicKeyString decrypts the data with public key string
func DecryptWithRSAPublicKeyString(publicKeyStr, cipher string) (string, error) {
	rsa := NewEmptyRSA()
	err := rsa.SetPublicKey(publicKeyStr)
	if err != nil {
		return constant.EmptyString, err
	}

	return rsa.DecryptWithPublicKey(cipher)
}

// EncryptWithRSAPrivateKeyString encrypts the data with private key string
func EncryptWithRSAPrivateKeyString(privateKeyStr, message string) (string, error) {
	rsa, err := NewRSAWithPrivateKeyString(privateKeyStr)
	if err != nil {
		return constant.EmptyString, err
	}

	return rsa.EncryptWithPrivateKey(message)
}

// DecryptWithRSAPrivateKeyString decrypts the data with private key string
func DecryptWithRSAPrivateKeyString(privateKeyStr, cipher string) (string, error) {
	rsa, err := NewRSAWithPrivateKeyString(privateKeyStr)
	if err != nil {
		return constant.EmptyString, err
	}

	return rsa.DecryptWithPrivateKey(cipher)
}

// EncryptWithSM2PublicKeyString encrypts the data with public key string
func EncryptWithSM2PublicKeyString(publicKeyStr, message string) (string, error) {
	sm2, err := NewSM2WithPublicKeyHexString(publicKeyStr)
	if err != nil {
		return constant.EmptyString, err
	}

	return sm2.Encrypt(message)
}

// DecryptWithSM2PrivateKeyString decrypts the data with public key string
func DecryptWithSM2PrivateKeyString(privateKeyStr, cipher string) (string, error) {
	sm2, err := NewSM2WithPrivateKeyHexString(privateKeyStr)
	if err != nil {
		return constant.EmptyString, err
	}

	return sm2.Decrypt(cipher)
}
