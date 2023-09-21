package common

import (
	"github.com/buger/jsonparser"
)

// KeyPathExists checks if the key path exists in the json data
func KeyPathExists(data []byte, keys ...string) bool {
	_, _, _, err := jsonparser.Get(data, keys...)
	if err != nil {
		return false
	}

	return true
}
