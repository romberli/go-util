package constant

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDefaultRandomTime(t *testing.T) {
	asst := assert.New(t)

	asst.Equal(DefaultRandomTime.Format(DefaultTimeLayout), DefaultRandomTimeString)
}
