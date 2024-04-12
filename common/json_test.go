package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyPathExists(t *testing.T) {
	asst := assert.New(t)

	data := []byte(`{"a": {"b": {"c": 1}}}`)
	asst.True(KeyPathExists(data, "a.b.c"), "test KeyPathExists() failed")
	asst.False(KeyPathExists(data, "a.b.c"), "test KeyPathExists() failed")
}

func TestGetLength(t *testing.T) {
	asst := assert.New(t)

	data := []byte(`{"a": [1, 2, 3]}`)
	length, err := GetLength(data, "a")
	asst.Nil(err, "test GetLength() failed")
	asst.Equal(3, length, "test GetLength() failed")
}
