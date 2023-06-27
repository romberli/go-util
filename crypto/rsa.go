package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"

	hrsa "github.com/hnlq715/rsa"

	"github.com/romberli/go-util/constant"
)

const (
	DefaultRSAKeySize = 512
)

type RSA struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewRSA returns a new *RSA
func NewRSA() (*RSA, error) {
	return newRSAWithKeySize(DefaultRSAKeySize)
}

// NewRSAWithKeySize returns a new *RSA with given key size
func NewRSAWithKeySize(size int) (*RSA, error) {
	return newRSAWithKeySize(size)
}

// NewRSAWithPrivateKey returns a new *RSA with given private key
func NewRSAWithPrivateKey(privateKey *rsa.PrivateKey) *RSA {
	return newRSAWithPrivateKey(privateKey)
}

// NewRSAWithPrivateKeyString returns a new *RSA with given private key base64 string
func NewRSAWithPrivateKeyString(privateKeyStr string) (*RSA, error) {
	privateKey, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		return nil, err
	}

	priv, err := x509.ParsePKCS1PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	return newRSAWithPrivateKey(priv), nil
}

// NewEmptyRSA returns a new *RSA
func NewEmptyRSA() *RSA {
	return &RSA{}
}

// newRSAWithKeySize returns a new *RSA
func newRSAWithKeySize(size int) (*RSA, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return nil, err
	}

	return newRSAWithPrivateKey(privateKey), nil
}

// newRSAWithPrivateKey returns a new *RSA with given private key
func newRSAWithPrivateKey(privateKey *rsa.PrivateKey) *RSA {
	return &RSA{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}
}

// GetPrivateKey gets the private key base64 string
func (r *RSA) GetPrivateKey() (string, error) {
	b := x509.MarshalPKCS1PrivateKey(r.privateKey)

	return base64.StdEncoding.EncodeToString(b), nil
}

// GetPublicKey gets the public key base64 string
func (r *RSA) GetPublicKey() (string, error) {
	b := x509.MarshalPKCS1PublicKey(r.publicKey)

	return base64.StdEncoding.EncodeToString(b), nil
}

// SetPrivateKey sets the private key with base64 string
func (r *RSA) SetPrivateKey(privateKeyStr string) error {
	b, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		return err
	}

	priv, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		return err
	}

	r.privateKey = priv
	r.publicKey = &priv.PublicKey

	return nil
}

// SetPublicKey sets the public key with base64 string
func (r *RSA) SetPublicKey(publicKeyStr string) error {
	b, err := base64.StdEncoding.DecodeString(publicKeyStr)
	if err != nil {
		return err
	}

	pub, err := x509.ParsePKCS1PublicKey(b)
	if err != nil {
		return err
	}

	r.publicKey = pub

	return nil
}

// Encrypt encrypts the string and returns the base64 string
func (r *RSA) Encrypt(message string) (string, error) {
	return r.EncryptWithPublicKey(message)
}

// Decrypt decrypts the base64 string
func (r *RSA) Decrypt(cipher string) (string, error) {
	return r.DecryptWithPrivateKey(cipher)
}

// EncryptWithPublicKey encrypts the message with public key
func (r *RSA) EncryptWithPublicKey(message string) (string, error) {
	cipher, err := rsa.EncryptPKCS1v15(rand.Reader, r.publicKey, []byte(message))
	if err != nil {
		return constant.EmptyString, err
	}

	return base64.StdEncoding.EncodeToString(cipher), nil
}

// DecryptWithPrivateKey decrypts the cipher with private key
func (r *RSA) DecryptWithPrivateKey(cipher string) (string, error) {
	c, err := base64.StdEncoding.DecodeString(cipher)
	if err != nil {
		return constant.EmptyString, err
	}

	message, err := rsa.DecryptPKCS1v15(rand.Reader, r.privateKey, c)
	if err != nil {
		return constant.EmptyString, err
	}

	return string(message), nil
}

// EncryptWithPrivateKey encrypts the message with private key
func (r *RSA) EncryptWithPrivateKey(message string) (string, error) {
	b, err := hrsa.PriKeyByte(r.privateKey, []byte(message), true)
	if err != nil {
		return constant.EmptyString, err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

// DecryptWithPublicKey decrypts the cipher with public key
func (r *RSA) DecryptWithPublicKey(cipher string) (string, error) {
	c, err := base64.StdEncoding.DecodeString(cipher)
	if err != nil {
		return constant.EmptyString, err
	}

	b, err := hrsa.PubKeyByte(r.publicKey, c, false)
	if err != nil {
		return constant.EmptyString, err
	}

	return string(b), nil
}
