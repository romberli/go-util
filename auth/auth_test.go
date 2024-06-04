package auth

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/common"
)

const (
	testSecretKey = "test_secret_key"
)

var testAuth *Auth

func init() {
	testAuth = NewAuth(common.StringToBytes(testSecretKey))
}

func TestAuth_Sign(t *testing.T) {
	asst := assert.New(t)

	jwtToken, err := testAuth.Sign()
	asst.Nil(err, common.CombineMessageWithError("test Sign() failed", err))
	t.Log(jwtToken)
}

func TestAuth_SignWithClaims(t *testing.T) {
	asst := assert.New(t)

	claims := jwt.MapClaims{
		"username": "test",
		"password": "test",
	}
	jwtToken, err := testAuth.SignWithMethodAndClaims(DefaultSignMethod, claims, NewGZIPEncodeFunc())
	asst.Nil(err, common.CombineMessageWithError("test SignWithMethodAndClaims() failed", err))
	t.Log(jwtToken)
}
