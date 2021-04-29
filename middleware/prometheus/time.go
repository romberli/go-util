package prometheus

import (
	"time"

	apiv1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

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

// NewTimeRangeWithRange returns a new TimeRange with given apiv1.Range
func NewTimeRangeWithRange(r apiv1.Range) TimeRange {
	return TimeRange{r}
}

// GetRange returns apiv1.Range
func (tr TimeRange) GetRange() apiv1.Range {
	return tr.Range
}
