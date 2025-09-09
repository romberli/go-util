package crypto

import (
	"github.com/romberli/go-util/constant"
)

// EncryptWithPublicKeyString encrypts the data with public key string
func EncryptWithPublicKeyString(publicKeyStr, message string) (string, error) {
	rsa := NewEmptyRSA()
	err := rsa.SetPublicKey(publicKeyStr)
	if err != nil {
		return constant.EmptyString, err
	}

	return rsa.EncryptWithPublicKey(message)
}

// DecryptWithPublicKeyString decrypts the data with public key string
func DecryptWithPublicKeyString(publicKeyStr, cipher string) (string, error) {
	rsa := NewEmptyRSA()
	err := rsa.SetPublicKey(publicKeyStr)
	if err != nil {
		return constant.EmptyString, err
	}

	return rsa.DecryptWithPublicKey(cipher)
}

// EncryptWithPrivateKeyString encrypts the data with private key string
func EncryptWithPrivateKeyString(privateKeyStr, message string) (string, error) {
	rsa, err := NewRSAWithPrivateKeyString(privateKeyStr)
	if err != nil {
		return constant.EmptyString, err
	}

	return rsa.EncryptWithPrivateKey(message)
}

// DecryptWithPrivateKeyString decrypts the data with private key string
func DecryptWithPrivateKeyString(privateKeyStr, cipher string) (string, error) {
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
