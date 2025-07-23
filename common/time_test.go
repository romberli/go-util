package common

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/constant"
)

func TestTime_All(t *testing.T) {
	TestTime_ConvertStringToTime(t)
}

func TestTime_ConvertStringToTime(t *testing.T) {
	asst := assert.New(t)

	timeString := "2025-09-25 10:02:05"
	time, err := ConvertStringToTime(timeString)
	asst.Nil(err, "test ConvertStringToTime() failed")
	t.Logf("time: %s", time.Format(constant.TimeLayoutMicrosecond))

	timeString = "2025-12-28 10:02:05.000"
	time, err = ConvertStringToTime(timeString)
	asst.Nil(err, "test ConvertStringToTime() failed")
	t.Logf("time: %s", time.Format(constant.TimeLayoutMicrosecond))

	timeString = "2025-07-09 10:02:05.000000"
	time, err = ConvertStringToTime(timeString)
	asst.Nil(err, "test ConvertStringToTime() failed")
	t.Logf("time: %s", time.Format(constant.TimeLayoutMicrosecond))

	timeString = "2025-12-28 10:02:05.012"
	time, err = ConvertStringToTime(timeString)
	asst.Nil(err, "test ConvertStringToTime() failed")
	t.Logf("time: %s", time.Format(constant.TimeLayoutMicrosecond))

	timeString = "2025-07-09 10:02:05.123456"
	time, err = ConvertStringToTime(timeString)
	asst.Nil(err, "test ConvertStringToTime() failed")
	t.Logf("time: %s", time.Format(constant.TimeLayoutMicrosecond))
}
