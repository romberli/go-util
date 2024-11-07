package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
)

const (
	DefaultAESKeySize = AESKeySize32

	AESKeySize16 = 16
	AESKeySize24 = 24
	AESKeySize32 = 32
)

type AES struct {
	key string
}

// NewAES returns a new *AES
func NewAES() (*AES, error) {
	return newAESWithKeySize(DefaultAESKeySize)
}

// NewAESWithKeySize returns a new *AES with given key size
func NewAESWithKeySize(size int) (*AES, error) {
	return newAESWithKeySize(size)
}

// NewAESWithKey returns a new *AES with given key
func NewAESWithKey(key string) *AES {
	return newAESWithKey(key)
}

// newAESWithKeySize returns a new *AES
func newAESWithKeySize(size int) (*AES, error) {
	key := make([]byte, size)

	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	return newAESWithKey(padKey(base64.StdEncoding.EncodeToString(key))), nil
}

// newAESWithKey returns a new *AES with given key
func newAESWithKey(key string) *AES {
	return &AES{
		key: padKey(key),
	}
}

func (a *AES) Encrypt(message string) (string, error) {
	// create block
	block, err := aes.NewCipher(common.StringToBytes(a.key))
	if err != nil {
		return constant.EmptyString, err
	}

	plaintext := common.StringToBytes(message)
	// pad message
	blockSize := block.BlockSize()
	padText := a.applyPKCS5Padding(plaintext, blockSize)

	// generate a random initialization vector
	ciphertext := make([]byte, blockSize+len(padText))
	iv := ciphertext[:blockSize]
	if _, err := rand.Read(iv); err != nil {
		return constant.EmptyString, err
	}

	// create a new CBC encryptor
	stream := cipher.NewCBCEncrypter(block, iv)
	stream.CryptBlocks(ciphertext[blockSize:], padText)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (a *AES) Decrypt(ciphertextBase64 string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return constant.EmptyString, err
	}

	// create block
	block, err := aes.NewCipher(common.StringToBytes(a.key))
	if err != nil {
		return "", err
	}

	// get IV and ciphertext
	blockSize := block.BlockSize()
	iv := ciphertext[:blockSize]
	ciphertext = ciphertext[blockSize:]

	// create a new CBC decrypter
	stream := cipher.NewCBCDecrypter(block, iv)
	// decrypt cipher
	stream.CryptBlocks(ciphertext, ciphertext)

	// remove padding
	padding := int(ciphertext[len(ciphertext)-1])
	plaintext := ciphertext[:len(ciphertext)-padding]

	return string(plaintext), nil
}

func (a *AES) applyPKCS5Padding(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(plaintext, padText...)
}
