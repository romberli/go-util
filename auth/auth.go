package auth

import (
	"encoding/json"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pingcap/errors"
)

var (
	DefaultSignMethod = jwt.SigningMethodHS512
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

// Sign signs with the default method and claims
func (a *Auth) Sign() (string, error) {
	return a.SignWithMethodAndClaims(DefaultSignMethod, jwt.MapClaims{}, nil)
}

// SignWithMethodAndClaims signs with the given method and claims
func (a *Auth) SignWithMethodAndClaims(method jwt.SigningMethod, claims jwt.MapClaims, ef EncodeFunc) (string, error) {
	token := NewTokenWithClaims(method, claims)

	return token.SignedString(a.secretKey, ef)
}

// Parse parses the payload from the token, it verifies the signature
func (a *Auth) Parse(tokenString string, in interface{}) error {
	parser := NewParserWithDefault()
	token, err := parser.Parse(tokenString, a.secretKey)
	if err != nil {
		return err
	}

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

// ParseUnverified parses the payload from the token, it does not verify the signature
func (a *Auth) ParseUnverified(tokenString string, in interface{}) error {
	parser := NewParserWithDefault()
	token, err := parser.ParseUnverified(tokenString)
	if err != nil {
		return err
	}

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
