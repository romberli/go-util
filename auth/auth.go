package auth

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pingcap/errors"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
)

var (
	DefaultSignMethod = jwt.SigningMethodRS256
)

type Auth struct {
	secretKey []byte
}

// NewAuth returns a new *Auth
func NewAuth(secretKey []byte) *Auth {
	return &Auth{
		secretKey: secretKey,
	}
}

// NewAuthWithDefault returns a new *Auth with empty secret key
func NewAuthWithDefault() *Auth {
	return NewAuth([]byte{})
}

// Sign signs with the default method and claims
func (a *Auth) Sign() (string, error) {
	return a.SignWithMethodAndClaims(DefaultSignMethod, jwt.MapClaims{}, nil)
}

// SignWithMethodAndClaims signs with the given method and claims
func (a *Auth) SignWithMethodAndClaims(method jwt.SigningMethod, claims jwt.MapClaims, ef EncodeFunc) (string, error) {
	token := NewTokenWithClaims(method, claims)

	var key interface{}

	switch method {
	case jwt.SigningMethodRS256, jwt.SigningMethodRS384, jwt.SigningMethodRS512:
		secretKey, err := base64.StdEncoding.DecodeString(common.BytesToString(a.secretKey))
		if err != nil {
			return constant.EmptyString, errors.Trace(err)
		}

		key, err = x509.ParsePKCS1PrivateKey(secretKey)
		if err != nil {
			return constant.EmptyString, errors.Trace(err)
		}
	default:
		key = a.secretKey
	}

	return token.SignedString(key, ef)
}

// Parse parses the payload from the token, it verifies the signature
func (a *Auth) Parse(tokenString string, in interface{}) error {
	parser := NewParserWithDefault()
	token, err := parser.Parse(tokenString, a.secretKey)
	if err != nil {
		return err
	}

	return a.unmarshal(token, in)
}

// ParseUnverified parses the payload from the token, it does not verify the signature
func (a *Auth) ParseUnverified(tokenString string, in interface{}) error {
	parser := NewParserWithDefault()
	token, err := parser.ParseUnverified(tokenString)
	if err != nil {
		return err
	}

	return a.unmarshal(token, in)
}

func (a *Auth) unmarshal(token *Token, in interface{}) error {
	bytes, err := json.Marshal(token.Claims)
	if err != nil {
		return errors.Trace(err)
	}

	err = json.Unmarshal(bytes, in)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}
