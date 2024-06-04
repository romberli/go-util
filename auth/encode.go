package auth

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"

	"github.com/romberli/go-util/constant"
)

// NewGZIPOption returns a new EncodeFunc with gzip compression, it only compresses the payload
func NewGZIPEncodeFunc() EncodeFunc {
	return func(token *Token, payload []byte) (string, error) {
		delete(token.Header, "typ")
		token.Header["zip"] = "GZIP"

		var buffer bytes.Buffer
		w := gzip.NewWriter(&buffer)
		defer func() {
			_ = w.Close()
		}()

		_, err := w.Write(payload)
		if err != nil {
			return constant.EmptyString, err
		}

		err = w.Flush()
		if err != nil {
			return constant.EmptyString, err
		}

		err = w.Close()
		if err != nil {
			return constant.EmptyString, err
		}

		return base64.RawURLEncoding.EncodeToString(buffer.Bytes()), nil
		// return buffer.String(), nil
	}
}
