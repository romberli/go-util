package prometheus

import (
	"time"

	apiv1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

const DefaultStep = time.Minute

type TimeRange struct {
	apiv1.Range
}

// NewTimeRange returns a new TimeRange
func NewTimeRange(start, end time.Time, step time.Duration) TimeRange {
	return TimeRange{
		Range: apiv1.Range{
			Start: start,
			End:   end,
			Step:  step,
		},
	}
}

// GetRange returns apiv1.Range
func (tr TimeRange) GetRange() apiv1.Range {
	return tr.Range
}
