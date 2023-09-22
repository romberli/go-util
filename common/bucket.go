package common

import (
	"time"

	"github.com/pingcap/errors"

	"github.com/romberli/go-util/constant"
)

const (
	defaultSleepTime = 10 * time.Millisecond
	minInterval      = 10 * time.Millisecond
)

type Bucket struct {
	capacity int
	num      int
	interval time.Duration
	pause    bool
	ch       chan struct{}
}

// NewBucket returns a new *Bucket
func NewBucket(capacity, num int, interval time.Duration) (*Bucket, error) {
	return newBucket(capacity, num, interval)
}

// newBucket returns a new *Bucket
func newBucket(capacity, num int, interval time.Duration) (*Bucket, error) {
	b := &Bucket{
		capacity: capacity,
		num:      num,
		interval: interval,
		pause:    false,
		ch:       make(chan struct{}, capacity),
	}

	err := b.validate()
	if err != nil {
		return nil, err
	}

	b.put(capacity)

	go b.supply()

	return b, nil
}

// Get gets a token from bucket, if bucket is empty, it will return error immediately
func (b *Bucket) Get() error {
	select {
	case <-b.ch:
		return nil
	default:
		return errors.Errorf("bucket is empty, please try again later")
	}
}

// GetWithTimeout gets a token from bucket, if bucket is empty, it will return error after timeout
// note that each time it failed to get the token, it will sleep for 10ms,
// so beware of the interval value when initializing the bucket
func (b *Bucket) GetWithTimeout(timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-b.ch:
		return nil
	case <-timer.C:
		return errors.Errorf("timeout exceeded, but the bucket is still empty, please try again later")
	}
}

// GetForever gets a token from bucket, if bucket is empty, it will wait forever until it gets a token successfully
func (b *Bucket) GetForever() {
	<-b.ch
}

func (b *Bucket) Pause() {
	b.pause = true
}

func (b *Bucket) Resume() {
	b.pause = false
}

func (b *Bucket) validate() error {
	if b.capacity <= constant.ZeroInt {
		return errors.Errorf("capacity must be greater than 0, %d is not valid", b.capacity)
	}

	if b.num <= constant.ZeroInt {
		return errors.Errorf("num must be greater than 0, %d is not valid", b.num)
	}

	if b.interval < minInterval {
		return errors.Errorf("interval must be greater or equal than 10ms, %d is not valid", b.interval)
	}

	return nil
}

// supply puts the specified number of tokens to the bucket every interval
func (b *Bucket) supply() {
	timer := time.NewTimer(b.interval)
	defer timer.Stop()

	for {
		timer.Reset(b.interval)

		if !b.pause {
			select {
			case <-timer.C:
				b.put(b.num)
			}
		}
	}
}

// put puts the specified number of tokens to the bucket
func (b *Bucket) put(num int) {
	for i := constant.ZeroInt; i < num; i++ {
		select {
		case b.ch <- struct{}{}:
			continue
		default:
			// the bucket is full
			return
		}
	}
}
