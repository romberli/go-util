package common

import (
	"os"
	"testing"
	"time"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"
)

const (
	testBucketCapacity = 5
	testBucketNum      = 5
	testBucketInterval = 100 * time.Millisecond
)

var testBucket *Bucket

func init() {
	testBucket = testInitBucket()
}

func testInitBucket() *Bucket {
	b, err := NewBucket(testBucketCapacity, testBucketNum, testBucketInterval)
	if err != nil {
		log.Errorf("init bucket failed. error:\n%+v", err)
		os.Exit(constant.DefaultAbnormalExitCode)
	}

	return b
}

func TestBucket_All(t *testing.T) {
	TestBucket_Get(t)
	TestBucket_GetWithTimeout(t)
	TestBucket_GetForever(t)
	TestBucket_Pause(t)
	TestBucket_Resume(t)
	TestBucket_put(t)
	TestBucket_supply(t)
}

func TestBucket_Get(t *testing.T) {
	asst := assert.New(t)

	err := testBucket.Get()
	asst.Nil(err, "test Get() failed. error:\n%+v", err)
}

func TestBucket_GetWithTimeout(t *testing.T) {
	asst := assert.New(t)

	testBucket.Resume()
	asst.False(testBucket.pause, "test GetWithTimeout() failed")
	time.Sleep(testBucket.interval * 2)

	err := testBucket.GetWithTimeout(testBucketInterval)
	asst.Nil(err, "test GetWithTimeout() failed. error:\n%+v", err)
	time.Sleep(testBucket.interval)
	testBucket.Pause()
	time.Sleep(testBucket.interval)

	for i := constant.ZeroInt; i < testBucket.capacity; i++ {
		err = testBucket.Get()
		asst.Nil(err, "test GetWithTimeout() failed. error:\n%+v", err)
	}

	err = testBucket.GetWithTimeout(testBucketInterval)
	asst.Zero(len(testBucket.ch), "test GetWithTimeout() failed")
	asst.NotNil(err, "test GetWithTimeout() failed")
}

func TestBucket_GetForever(t *testing.T) {
	asst := assert.New(t)

	testBucket.Resume()
	asst.False(testBucket.pause, "test GetWithTimeout() failed")
	time.Sleep(testBucket.interval * 2)

	testBucket.Pause()
	asst.True(testBucket.pause, "test GetForever() failed")
	time.Sleep(testBucket.interval)

	for i := constant.ZeroInt; i < testBucket.capacity; i++ {
		err := testBucket.Get()
		asst.Nil(err, "test GetForever() failed. error:\n%+v", err)
	}
	asst.Zero(len(testBucket.ch), "test GetForever() failed")

	startTime := time.Now()

	go func() {
		time.Sleep(time.Second)
		testBucket.Resume()
	}()

	testBucket.GetForever()
	waitTIme := time.Now().Sub(startTime)
	asst.Equal(int(waitTIme.Seconds()), int(time.Second.Seconds()), "test GetForever() failed")
}

func TestBucket_Pause(t *testing.T) {
	asst := assert.New(t)

	testBucket.Resume()
	asst.False(testBucket.pause, "test GetWithTimeout() failed")
	time.Sleep(testBucket.interval * 2)

	testBucket.Pause()
	asst.True(testBucket.pause, "test Pause() failed")
}

func TestBucket_Resume(t *testing.T) {
	asst := assert.New(t)

	testBucket.Resume()
	asst.False(testBucket.pause, "test GetWithTimeout() failed")
	time.Sleep(testBucket.interval * 2)

	testBucket.Pause()
	asst.True(testBucket.pause, "test Resume() failed")
	time.Sleep(testBucket.interval)

	testBucket.Resume()
	asst.False(testBucket.pause, "test Resume() failed")
}

func TestBucket_put(t *testing.T) {
	asst := assert.New(t)

	testBucket.Resume()
	asst.False(testBucket.pause, "test GetWithTimeout() failed")
	time.Sleep(testBucket.interval * 2)

	for i := constant.ZeroInt; i < testBucket.capacity; i++ {
		err := testBucket.Get()
		asst.Nil(err, "test put() failed. error:\n%+v", err)
	}
	asst.Zero(len(testBucket.ch), "test put() failed")
	testBucket.put(testBucketNum)
	asst.Equal(testBucketNum, len(testBucket.ch), "test put() failed")
}

func TestBucket_supply(t *testing.T) {
	asst := assert.New(t)

	testBucket.Resume()
	asst.False(testBucket.pause, "test GetWithTimeout() failed")
	time.Sleep(testBucket.interval * 2)

	testBucket.Pause()
	asst.True(testBucket.pause, "test supply() failed")
	time.Sleep(testBucket.interval)

	for i := constant.ZeroInt; i < testBucket.capacity; i++ {
		err := testBucket.Get()
		asst.Nil(err, "test supply() failed. error:\n%+v", err)
	}
	asst.Zero(len(testBucket.ch), "test supply() failed")
	testBucket.Resume()
	time.Sleep(testBucket.interval * 2)
	asst.Equal(testBucket.capacity, len(testBucket.ch), "test supply() failed")
}
