package common

import (
	"sync"
	"time"
)

type Timer struct {
	mutex         *sync.Mutex
	nextStartTime time.Time
}

// NewTimer returns a new *Timer
func NewTimer(mutex *sync.Mutex, nextStartTime time.Time) *Timer {
	return newTimer(mutex, nextStartTime)
}

// NewTimerWithDefault returns a new *Timer
func NewTimerWithDefault() *Timer {
	return newTimer(&sync.Mutex{}, time.Now())
}

// newTimer returns a new *Timer
func newTimer(mutex *sync.Mutex, nextStartTime time.Time) *Timer {
	return &Timer{
		mutex:         mutex,
		nextStartTime: nextStartTime,
	}
}

// GetMutex returns the mutex
func (t *Timer) GetMutex() *sync.Mutex {
	return t.mutex
}

// GetNextStartTime returns the next start time
func (t *Timer) GetNextStartTime() time.Time {
	t.GetMutex().Lock()
	defer t.GetMutex().Unlock()

	return t.nextStartTime
}

func (t *Timer) SetNextStartTime(nextStartTime time.Time) {
	t.GetMutex().Lock()
	defer t.GetMutex().Unlock()

	t.nextStartTime = nextStartTime
}
