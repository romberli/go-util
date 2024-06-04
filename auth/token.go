package auth

import (
	"encoding/base64"
	"encoding/json"

	"github.com/golang-jwt/jwt/v5"

	"github.com/romberli/go-util/constant"
)

type EncodeFunc func(*Token, []byte) (string, error)

type Token struct {
	Raw        string
	Method     jwt.SigningMethod
	Header     map[string]interface{}
	Claims     jwt.Claims
	Signature  []byte
	Valid      bool
	EncodeFunc EncodeFunc
}

// NewToken returns a new *Token
func NewToken(method jwt.SigningMethod) *Token {
	return NewTokenWithClaims(method, jwt.MapClaims{}, nil)
}

// NewTokenWithClaims returns a new *Token with the specified signing method and claims
func NewTokenWithClaims(method jwt.SigningMethod, claims jwt.Claims, ef EncodeFunc) *Token {
	return &Token{
		Header: map[string]interface{}{
			"typ": "JWT",
			"alg": method.Alg(),
		},
		Claims:     claims,
		Method:     method,
		EncodeFunc: ef,
	}
}

// SignedString creates and returns a complete, signed JWT.
func (t *Token) SignedString(key interface{}) (string, error) {
	ss, err := t.SigningString()
	if err != nil {
		return constant.EmptyString, err
	}

	sig, err := t.Method.Sign(ss, key)
	if err != nil {
		return constant.EmptyString, err
	}

	return ss + constant.DotString + base64.RawURLEncoding.EncodeToString(sig), nil
}

// SigningString generates the signing string
func (t *Token) SigningString() (string, error) {
	c, err := json.Marshal(t.Claims)
	if err != nil {
		return constant.EmptyString, err
	}

	es, err := t.EncodeSegment(c)
	if err != nil {
		return constant.EmptyString, err
	}

	h, err := json.Marshal(t.Header)
	if err != nil {
		return constant.EmptyString, err
	}

	return base64.RawURLEncoding.EncodeToString(h) + constant.DotString + es, nil
}

// EncodeSegment encodes a JWT specific, this is the place that the EncodeFunc are applied.
func (t *Token) EncodeSegment(seg []byte) (string, error) {
	if t.EncodeFunc != nil {
		return t.EncodeFunc(t, seg)
	}

	return base64.RawURLEncoding.EncodeToString(seg), nil
}
