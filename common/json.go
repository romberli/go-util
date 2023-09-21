package common

import (
	"github.com/buger/jsonparser"
)

// KeyExists checks if the key exists in the json
func KeyExists(json []byte, keys ...string) bool {
	_, _, _, err := jsonparser.Get(json, keys...)
	if err != nil {
		return false
	}

	return true
}
