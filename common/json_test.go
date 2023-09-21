package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyPathExists(t *testing.T) {
	asst := assert.New(t)

	data := []byte(`{"a": {"b": {"c": 1}}}`)
	asst.True(KeyPathExists(data, "a", "b", "c"), "test KeyPathExists() failed")
	asst.False(KeyPathExists(data, "a", "b", "d"), "test KeyPathExists() failed")
}
