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
	TestIdentifier_CheckIdentifierWithDefault(t)
}

func TestIdentifier_CheckIdentifier(t *testing.T) {
	asst := assert.New(t)

	err := CheckIdentifier(testIdentifierWithSpecialChar, true, DefaultIdentifierSpecialCharStr)
	asst.NotNil(err, "test CheckIdentifier() failed")
	err = CheckIdentifier(testIdentifierLeadingDigit, true, DefaultIdentifierSpecialCharStr)
	asst.NotNil(err, "test CheckIdentifier() failed")
	err = CheckIdentifier(testIdentifierValid, true, DefaultIdentifierSpecialCharStr)
	asst.Nil(err, "test CheckIdentifier() failed")
}

func TestIdentifier_CheckIdentifierWithDefault(t *testing.T) {
	asst := assert.New(t)

	err := CheckIdentifierWithDefault(testIdentifierWithSpecialChar)
	asst.NotNil(err, "test CheckIdentifier() failed")
	err = CheckIdentifierWithDefault(testIdentifierLeadingDigit)
	asst.NotNil(err, "test CheckIdentifier() failed")
	err = CheckIdentifierWithDefault(testIdentifierValid)
	asst.Nil(err, "test CheckIdentifier() failed")
}
