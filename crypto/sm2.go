package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/pingcap/errors"
	"github.com/tjfoc/gmsm/sm2"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
)

type SM2 struct {
	privateKey *sm2.PrivateKey
	publicKey  *sm2.PublicKey
}

// NewSM2 returns a new *SM2
func NewSM2() (*SM2, error) {
	privateKey, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return newSM2WithPrivateKey(privateKey), nil
}

// NewSM2WithPrivateKeyString returns a new *SM2 with given private key hex string
func NewSM2WithPrivateKeyHexString(privateKeyStr string) (*SM2, error) {
	privateKey, err := ConvertHexToSM2PrivateKey(privateKeyStr)
	if err != nil {
		return nil, err
	}

	return newSM2WithPrivateKey(privateKey), nil
}

// NewSM2WithPrivateKey returns a new *SM2 with given private key
func NewSM2WithPrivateKey(privateKey *sm2.PrivateKey) *SM2 {
	return newSM2WithPrivateKey(privateKey)
}

// NewSM2WithPublicKey returns a new *SM2 with given public key
func NewSM2WithPublicKey(publicKey *sm2.PublicKey) *SM2 {
	return newSM2WithPublicKey(publicKey)
}

// NewSM2WithPublicKeyString returns a new *SM2 with given public key hex string
func NewSM2WithPublicKeyHexString(publicKeyStr string) (*SM2, error) {
	publicKey, err := ConvertHexToSM2PublicKey(publicKeyStr)
	if err != nil {
		return nil, err
	}

	return newSM2WithPublicKey(publicKey), nil
}

// NewEmptyRSA returns a new *SM2
func NewEmptySM2() *SM2 {
	return &SM2{}
}

// newSM2WithPrivateKey returns a new *SM2
func newSM2WithPrivateKey(privateKey *sm2.PrivateKey) *SM2 {
	return &SM2{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}
}

// newSM2WithPublicKey returns a new *SM2
func newSM2WithPublicKey(publicKey *sm2.PublicKey) *SM2 {
	return &SM2{
		publicKey: publicKey,
	}
}

// GetPrivateKey returns the private key
func (s *SM2) GetPrivateKey() *sm2.PrivateKey {
	return s.privateKey
}

// GetPublicKey returns the public key
func (s *SM2) GetPublicKey() *sm2.PublicKey {
	return s.publicKey
}

// GetPrivateKeyHexString returns the private key hex string
func (s *SM2) GetPrivateKeyHexString() string {
	return ConvertSM2PrivateKeyToHex(s.privateKey)
}

// GetPublicKeyHexString returns the public key hex string
func (s *SM2) GetPublicKeyHexString() string {
	return ConvertSM2PublicKeyToHex(s.publicKey)
}

// SetPrivateKey sets the private key
func (s *SM2) SetPrivateKey(privateKey *sm2.PrivateKey) {
	s.privateKey = privateKey
	s.publicKey = &privateKey.PublicKey
}

// SetPublicKey sets the public key
func (s *SM2) SetPublicKey(publicKey *sm2.PublicKey) {
	s.publicKey = publicKey
}

// Encrypt  encrypts the data
func (s *SM2) Encrypt(message string) (string, error) {
	return s.EncryptWithPublicKey(message, sm2.C1C3C2)
}

// Decrypt decrypts the data
func (s *SM2) Decrypt(cipher string) (string, error) {
	return s.DecryptWithPrivateKey(cipher, sm2.C1C3C2)
}

// EncryptWithPublicKey encrypts the data with public key
func (s *SM2) EncryptWithPublicKey(message string, mode int) (string, error) {
	cipher, err := sm2.Encrypt(s.publicKey, common.StringToBytes(message), rand.Reader, mode)
	if err != nil {
		return constant.EmptyString, errors.Trace(err)
	}

	return strings.ToUpper(hex.EncodeToString(cipher)), nil
}

// DecryptWithPrivateKey decrypts the data with private key
func (s *SM2) DecryptWithPrivateKey(cipher string, mode int) (string, error) {
	cipherBytes, err := hex.DecodeString(cipher)
	if err != nil {
		return constant.EmptyString, errors.Trace(err)
	}
	message, err := sm2.Decrypt(s.privateKey, cipherBytes, mode)
	if err != nil {
		return constant.EmptyString, errors.Trace(err)
	}

	return common.BytesToString(message), nil
}
