package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tjfoc/gmsm/sm2"

	"github.com/romberli/go-util/constant"
)

const (
	sm2PrivateKeyStr = "DA666C31B8CC743EF2D444172C03E3E002E5849802D309D871D3B2923DB2386A"
	sm2PublicKeyStr  = "040D159681373C89237E472E68B2AFD3F53505DDB5866DF70D70C9B9952658AC37AFAC8D62F81D6A4FB8D49078F1F08E60CEFE3C6664B26EFE79679BF17AD82ABF"
)

func verifyKeyMatch(privateKey *sm2.PrivateKey, publicKey *sm2.PublicKey) bool {
	calculatedX, calculatedY := privateKey.Curve.ScalarBaseMult(privateKey.D.Bytes())

	isMatch := calculatedX.Cmp(publicKey.X) == constant.ZeroInt && calculatedY.Cmp(publicKey.Y) == constant.ZeroInt
	isOnCurve := publicKey.Curve.IsOnCurve(publicKey.X, publicKey.Y)

	return isMatch && isOnCurve
}

func TestSM2_All(t *testing.T) {
	TestSM2_GetPrivateKeyHexString(t)
	TestSM2_Encrypt(t)
	TestSM2_Decrypt(t)
}

func TestSM2_GetPrivateKeyHexString(t *testing.T) {
	asst := assert.New(t)
	s, err := NewSM2()
	asst.Nil(err, "test SM2.GetPrivateKeyHexString() failed")
	privateKeyStr := s.GetPrivateKeyHexString()
	asst.Nil(err, "test SM2.GetPrivateKeyHexString() failed")
	asst.NotEmpty(privateKeyStr, "test SM2.GetPrivateKeyHexString() failed")
	t.Logf("private key: %s", privateKeyStr)
}

func TestSM2_Encrypt(t *testing.T) {
	asst := assert.New(t)
	s, err := NewSM2()
	asst.Nil(err, "test SM2.Encrypt() failed")
	cipher, err := s.Encrypt(defaultMessage)
	asst.Nil(err, "test SM2.Encrypt() failed")
	message, err := s.Decrypt(cipher)
	asst.Nil(err, "test SM2.Encrypt() failed")
	asst.Equal(defaultMessage, string(message), "test SM2.Encrypt() failed")
}

func TestSM2_Decrypt(t *testing.T) {
	asst := assert.New(t)
	s, err := NewSM2()
	asst.Nil(err, "test SM2.Decrypt() failed")

	privateKeyStr := s.GetPrivateKeyHexString()
	t.Logf("private key: %s", privateKeyStr)
	publicKeyStr := s.GetPublicKeyHexString()
	t.Logf("public key: %s", publicKeyStr)

	cipher, err := s.Encrypt(defaultMessage)
	asst.Nil(err, "test SM2.Decrypt() failed")
	message, err := s.Decrypt(cipher)
	asst.Nil(err, "test SM2.Decrypt() failed")
	asst.Equal(defaultMessage, string(message), "test SM2.Decrypt() failed")
}

func TestSM2_VerifyKey(t *testing.T) {
	asst := assert.New(t)
	privateKey, err := ConvertHexToSM2PrivateKey(sm2PrivateKeyStr)
	asst.Nil(err, "test SM2.VerifyKey() failed")
	publicKey, err := ConvertHexToSM2PublicKey(sm2PublicKeyStr)
	asst.Nil(err, "test SM2.VerifyKey() failed")
	isMatch := verifyKeyMatch(privateKey, publicKey)
	asst.True(isMatch, "test SM2.VerifyKey() failed")
}
