package auth

import (
	"github.com/golang-jwt/jwt/v5"
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
func (a *Auth) SignWithMethodAndClaims(method jwt.SigningMethod, claims jwt.Claims, ef EncodeFunc) (string, error) {
	token := NewTokenWithClaims(method, claims, ef)

	return token.SignedString(a.secretKey)
}
