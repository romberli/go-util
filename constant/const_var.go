package constant

import (
	"time"
)

var DefaultRandomTime = GetDefaultRandomTime()

func GetDefaultRandomTime() time.Time {
	t, _ := time.ParseInLocation(DefaultTimeLayout, DefaultRandomTimeString, time.Local)

	return t
}
