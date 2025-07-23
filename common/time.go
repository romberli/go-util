package common

import (
	"strings"
	"time"

	"github.com/pingcap/errors"

	"github.com/romberli/go-util/constant"
)

func ConvertStringToTime(s string) (time.Time, error) {
	if s == constant.EmptyString {
		return time.Time{}, nil
	}

	layout := constant.TimeLayoutSecond
	if strings.Contains(s, constant.DotString) {
		vl := strings.Split(s, constant.DotString)
		f := vl[constant.OneInt]
		if len(f) == constant.ThreeInt {
			layout = constant.TimeLayoutMillisecond
		} else if len(f) == constant.SixInt {
			layout = constant.TimeLayoutMicrosecond
		}
	}

	t, err := time.ParseInLocation(layout, s, time.Local)
	if err != nil {
		return time.Time{}, errors.Trace(err)
	}

	return t, nil
}
