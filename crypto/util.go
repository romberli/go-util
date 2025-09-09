package crypto

import (
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/pingcap/errors"
	"github.com/tjfoc/gmsm/sm2"

	"github.com/romberli/go-util/constant"
)

const (
	SM2PublicKeyPrefix                      = "04"
	validHexSM2PublicKeyLengthWithOutPrefix = 128
	validHexSM2PublicKeyLengthWithPrefix    = 130
)

func padAESKey(key string) string {
	keyLengths := []int{AESKeySize16, AESKeySize24, AESKeySize32}

	if len(key) > AESKeySize32 {
		return key[:AESKeySize32]
	}

	for _, length := range keyLengths {
		if len(key) < length {
			return key + strings.Repeat(constant.EqualString, length-len(key))
		} else if len(key) == length {
			return key
		}
	}

	return key
}

func ConvertHexToSM2PrivateKey(hexPrivateKeyStr string) (*sm2.PrivateKey, error) {
	keyBytes, err := hex.DecodeString(hexPrivateKeyStr)
	if err != nil {
		return nil, errors.Trace(err)
	}

	d := new(big.Int).SetBytes(keyBytes)
	curve := sm2.P256Sm2()
	n := curve.Params().N
	if d.Cmp(big.NewInt(constant.ZeroInt)) <= constant.ZeroInt || d.Cmp(n) >= constant.ZeroInt {
		return nil, errors.Errorf("invalid private key, d should be in the range of [1, %x], %x is not valid", n, d)
	}

	x, y := curve.ScalarBaseMult(d.Bytes())

	if !curve.IsOnCurve(x, y) {
		return nil, errors.New("invalid private key, (x, y) is not on the curve")
	}

	return &sm2.PrivateKey{
		PublicKey: sm2.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
		D: d,
	}, nil
}

func ConvertHexToSM2PublicKey(hexPublicKeyStr string) (*sm2.PublicKey, error) {
	length := len(hexPublicKeyStr)
	if length != validHexSM2PublicKeyLengthWithPrefix &&
		length != validHexSM2PublicKeyLengthWithOutPrefix {
		return nil, errors.Errorf("public key length should be 128 or 130, %d is not valid", length)
	}
	if length == validHexSM2PublicKeyLengthWithOutPrefix &&
		!strings.HasSuffix(hexPublicKeyStr, SM2PublicKeyPrefix) {
		return nil, errors.Errorf("public key with prefix should start with %s", SM2PublicKeyPrefix)
	}

	if length == validHexSM2PublicKeyLengthWithPrefix {
		hexPublicKeyStr = hexPublicKeyStr[2:]
	}

	xHex := hexPublicKeyStr[0:64]
	yHex := hexPublicKeyStr[64:128]

	xBytes, err := hex.DecodeString(xHex)
	if err != nil {
		return nil, err
	}

	yBytes, err := hex.DecodeString(yHex)
	if err != nil {
		return nil, err
	}

	curve := sm2.P256Sm2()
	x := new(big.Int).SetBytes(xBytes)
	y := new(big.Int).SetBytes(yBytes)

	if !curve.IsOnCurve(x, y) {
		return nil, errors.New("invalid public key, (x, y) is not on the curve")
	}

	return &sm2.PublicKey{
		Curve: sm2.P256Sm2(),
		X:     x,
		Y:     y,
	}, nil
}

func ConvertSM2PrivateKeyToHex(privateKey *sm2.PrivateKey) string {

	dBytes := padSM2Key(privateKey.D)

	return strings.ToUpper(hex.EncodeToString(dBytes))
}

func ConvertSM2PublicKeyToHex(publicKey *sm2.PublicKey) string {
	xBytes := padSM2Key(publicKey.X)
	yBytes := padSM2Key(publicKey.Y)

	fullBytes := append([]byte{0x04}, xBytes...)
	fullBytes = append(fullBytes, yBytes...)

	return strings.ToUpper(hex.EncodeToString(fullBytes))
}

func padSM2Key(value *big.Int) []byte {
	buf := make([]byte, 32)
	value.FillBytes(buf)

	return buf
}
