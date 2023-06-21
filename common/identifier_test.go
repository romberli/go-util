package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testIdentifierWithSpecialChar = "test;"
	testIdentifierLeadingDigit    = "1test"

	testIdentifierValid = "test_01"
)

func TestIdentifier_All(t *testing.T) {
	TestIdentifier_CheckIdentifier(t)
}

func TestIdentifier_CheckIdentifier(t *testing.T) {
	asst := assert.New(t)

	err := CheckIdentifier(testIdentifierWithSpecialChar)
	asst.NotNil(err, "test CheckIdentifier() failed")
	err = CheckIdentifier(testIdentifierLeadingDigit)
	asst.NotNil(err, "test CheckIdentifier() failed")
	err = CheckIdentifier(testIdentifierValid)
	asst.Nil(err, "test CheckIdentifier() failed")
}
