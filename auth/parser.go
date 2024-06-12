package auth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pingcap/errors"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
)

type Parser struct {
	useJSONNumber bool
}

// NewParser returns a new *Parser
func NewParser(useJSONNumber bool) *Parser {
	return &Parser{
		useJSONNumber: useJSONNumber,
	}
}

// NewParserWithDefault returns a new *Parser with default options
func NewParserWithDefault() *Parser {
	return NewParser(true)
}

// Parse parses the token string and verifies the signature
func (p *Parser) Parse(tokenString, key string) (*Token, error) {
	token, err := p.ParseUnverified(tokenString)
	if err != nil {
		return nil, errors.Trace(err)
	}

	decoded, err := base64.RawURLEncoding.DecodeString(common.BytesToString(token.Signature))
	if err != nil {
		return nil, errors.Trace(err)
	}
	token.Signature = decoded

	parts := strings.Split(tokenString, constant.DotString)
	if len(parts) != constant.ThreeInt {
		return nil, errors.Errorf("token contains an invalid number of segments, expected 3, actual: %d", len(parts))
	}

	err = token.Method.Verify(strings.Join(parts[:constant.TwoInt], constant.DotString), token.Signature, key)
	if err != nil {
		return token, errors.Trace(err)
	}

	token.Valid = true

	return token, nil
}

// ParseUnverified parses the token string without verifying the signature
func (p *Parser) ParseUnverified(tokenString string) (*Token, error) {
	parts := strings.Split(tokenString, constant.DotString)
	if len(parts) != constant.ThreeInt {
		return nil, errors.Errorf("token contains an invalid number of segments, expected 3, actual: %d", len(parts))
	}
	// init token
	token := NewTokenWithRawString(tokenString)
	token.Signature = common.StringToBytes(parts[constant.TwoInt])
	// header
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[constant.ZeroInt])
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(headerBytes, &token.Header)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if token.Header == nil {
		return nil, errors.New("could not decode token header")
	}
	// payload
	var payloadBytes []byte
	jwtType, ok := token.Header[tokenTypeHeader]
	if ok {
		if jwtType != tokenJWTType {
			return nil, errors.Errorf("unsupported token type. key: %s, value: %s", tokenTypeHeader, jwtType)
		}

		payloadBytes, err = base64.RawURLEncoding.DecodeString(parts[constant.OneInt])
		if err != nil {
			return nil, err
		}
	} else {
		gzipType, ok := token.Header[tokenZIPHeader]
		if ok {
			if gzipType != tokenGZipType {
				return nil, errors.Errorf("unsupported token type. key: %s, value: %s", tokenZIPHeader, gzipType)
			}

			payloadBytes, err = NewGZIPDecodeFunc()(token, parts[constant.OneInt])
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.Errorf("unsupported token type. header: %s", token.Header)
		}
	}

	if p.useJSONNumber {
		decoder := json.NewDecoder(bytes.NewBuffer(payloadBytes))
		decoder.UseNumber()

		err = decoder.Decode(&token.Claims)
		if err != nil {
			return nil, errors.Trace(err)
		}
	} else {
		err = json.Unmarshal(payloadBytes, &token.Claims)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	// signature
	alg, ok := token.Header[tokenAlgorithmHeader]
	if !ok {
		return nil, errors.Errorf("token algorithm header is missing. header: %s", token.Header)
	}
	token.Method = jwt.GetSigningMethod(alg)
	if token.Method == nil {
		return token, errors.Errorf("signing method (alg) is unavailable")
	}

	return token, nil
}
