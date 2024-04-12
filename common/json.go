package common

import (
	"strings"

	"github.com/buger/jsonparser"
	"github.com/pingcap/errors"
	"github.com/tidwall/gjson"

	"github.com/romberli/go-util/constant"
)

// KeyPathExists checks if the key path exists in the json data
func KeyPathExists(data []byte, path string) bool {
	_, _, _, err := jsonparser.Get(data, strings.Split(path, constant.DotString)...)
	if err != nil {
		return false
	}

	return true
}

// GetLength returns the length of the json array, the value of the path should be a json array
func GetLength(data []byte, path string) (int, error) {
	if !gjson.ValidBytes(data) {
		return constant.ZeroInt, errors.Errorf("invalid json data. data: %s", BytesToString(data))
	}

	result := gjson.GetBytes(data, path)
	if !result.Exists() {
		return constant.ZeroInt, errors.Errorf("path not found. path: %s", path)
	}
	if !result.IsArray() {
		return constant.ZeroInt, errors.Errorf("path is not an array. path: %s", path)
	}

	return len(result.Array()), nil
}
