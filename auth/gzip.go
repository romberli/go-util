package auth

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"

	"github.com/pingcap/errors"

	"github.com/romberli/go-util/constant"
)

// NewGZIPEncodeFunc returns a new EncodeFunc with gzip compression, it only compresses the payload
func NewGZIPEncodeFunc() EncodeFunc {
	return func(token *Token, payload []byte) (string, error) {
		delete(token.Header, tokenTypeHeader)
		token.Header[tokenZIPHeader] = tokenGZipType

		var buffer bytes.Buffer
		w := gzip.NewWriter(&buffer)
		defer func() {
			_ = w.Close()
		}()

		_, err := w.Write(payload)
		if err != nil {
			return constant.EmptyString, errors.Trace(err)
		}

		err = w.Flush()
		if err != nil {
			return constant.EmptyString, errors.Trace(err)
		}

		err = w.Close()
		if err != nil {
			return constant.EmptyString, errors.Trace(err)
		}

		return base64.RawURLEncoding.EncodeToString(buffer.Bytes()), nil
	}
}
