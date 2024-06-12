package auth

import (
	"encoding/base64"
	"encoding/json"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pingcap/errors"

	"github.com/romberli/go-util/constant"
)

const (
	tokenTypeHeader      = "typ"
	tokenAlgorithmHeader = "alg"
	tokenZIPHeader       = "zip"
	tokenGZipType        = "GZIP"
	tokenJWTType         = "JWT"
)

type EncodeFunc func(*Token, []byte) (string, error)

type Token struct {
	Raw       string
	Method    jwt.SigningMethod
	Header    map[string]string
	Claims    jwt.MapClaims
	Signature []byte
	Valid     bool
}

// NewToken returns a new *Token
func NewToken(method jwt.SigningMethod) *Token {
	return NewTokenWithClaims(method, jwt.MapClaims{})
}

// NewTokenWithRawString returns a new *Token with the specified raw string
func NewTokenWithRawString(raw string) *Token {
	return &Token{
		Raw:    raw,
		Claims: jwt.MapClaims{},
	}
}

// NewTokenWithClaims returns a new *Token with the specified signing method and claims
func NewTokenWithClaims(method jwt.SigningMethod, claims jwt.MapClaims) *Token {
	return &Token{
		Header: map[string]string{
			tokenTypeHeader:      tokenJWTType,
			tokenAlgorithmHeader: method.Alg(),
		},
		Claims: claims,
		Method: method,
	}
}

// SignedString creates and returns a complete, signed JWT.
func (t *Token) SignedString(key interface{}, ef EncodeFunc) (string, error) {
	ss, err := t.SigningString(ef)
	if err != nil {
		return constant.EmptyString, err
	}

	sig, err := t.Method.Sign(ss, key)
	if err != nil {
		return constant.EmptyString, errors.Trace(err)
	}

	return ss + constant.DotString + base64.RawURLEncoding.EncodeToString(sig), nil
}

// SigningString generates the signing string
func (t *Token) SigningString(ef EncodeFunc) (string, error) {
	c, err := json.Marshal(t.Claims)
	if err != nil {
		return constant.EmptyString, errors.Trace(err)
	}

	es, err := t.EncodeSegment(c, ef)
	if err != nil {
		return constant.EmptyString, errors.Trace(err)
	}

	h, err := json.Marshal(t.Header)
	if err != nil {
		return constant.EmptyString, errors.Trace(err)
	}

	return base64.RawURLEncoding.EncodeToString(h) + constant.DotString + es, nil
}

// EncodeSegment encodes a JWT segment, this is the place that the EncodeFunc are applied.
func (t *Token) EncodeSegment(seg []byte, ef EncodeFunc) (string, error) {
	if ef != nil {
		encoded, err := ef(t, seg)
		if err != nil {
			return constant.EmptyString, errors.Trace(err)
		}

		return encoded, nil
	}

	return base64.RawURLEncoding.EncodeToString(seg), nil
}
