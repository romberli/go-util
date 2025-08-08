package viper

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pingcap/errors"
	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
)

const (
	testConfileFile = "./test_config.yaml"

	testLogFileKey                = "log.file"
	testLogFileValue              = "log/run.log"
	testLogMaxSizeKey             = "log.maxSize"
	testLogMaxSizeValue           = 20
	testServerPProfEnabledKey     = "server.pprof.enabled"
	testServerSwaggerKey          = "server.swagger"
	testAllowListKey              = "allowList"
	testServerSwaggerAllowListKey = "server.swagger.allowList"
	testServerRouterPathMapKey    = "server.router.pathMap"

	testReaderCount   = 10
	testConcurrentNum = 100
	testMaxModifyNum  = 5
)

var (
	testTempAllowListValue          = []string{"domain-0.localhost"}
	testServerSwaggerAllowListValue = []string{"localhost", "127.0.0.1", "192.168.137.2", "192.168.3.11"}
	testServerSwaggerValue          = map[string][]string{"allowList": testServerSwaggerAllowListValue}
	testServerRouterPathMapValue    = map[string]string{
		"/aaa1": "/aaa2",
		"/bbb1": "/bbb2",
	}
)

func init() {
	testReadConfigFile()
}

func testReadConfigFile() {
	SetConfigFile(testConfileFile)
	SetConfigType(DefaultFileType)
	err := ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func testCreateTempConfig(content string) (*os.File, error) {
	tempFile, err := os.Create("./config_temp.yaml")
	if err != nil {
		return nil, err
	}

	_, err = tempFile.WriteString(content)
	if err != nil {
		return nil, err
	}

	return tempFile, tempFile.Close()
}

func TestSafeViper_All(t *testing.T) {
	TestSafeViper_GetBool(t)
	TestSafeViper_GetInt(t)
	TestSafeViper_GetString(t)
	TestSafeViper_GetStringSlice(t)
	TestSafeViper_GetStringMap(t)
	TestSafeViper_GetStringMapString(t)
	TestSafeViper_GetStringMapStringSlice(t)
}

func TestSafeViper_GetBool(t *testing.T) {
	asst := assert.New(t)

	asst.True(GetBool(testServerPProfEnabledKey), "test GetBool() failed")
}

func TestSafeViper_GetInt(t *testing.T) {
	asst := assert.New(t)

	asst.Equal(testLogMaxSizeValue, GetInt(testLogMaxSizeKey), "test GetInt() failed")
}

func TestSafeViper_GetString(t *testing.T) {
	asst := assert.New(t)

	asst.Equal(testLogFileValue, GetString(testLogFileKey), "test GetString() failed")
}

func TestSafeViper_GetStringSlice(t *testing.T) {
	asst := assert.New(t)

	asst.True(common.ElementEqualInSlice(testServerSwaggerAllowListValue, GetStringSlice(testServerSwaggerAllowListKey)),
		"test GetStringSlice() failed")
}

func TestSafeViper_GetStringMap(t *testing.T) {
	asst := assert.New(t)

	val := GetStringMap(testServerRouterPathMapKey)
	asst.Equal(len(testServerRouterPathMapValue), len(val),
		"test GetStringMap() failed")
	for k, v := range val {
		asst.Equal(testServerRouterPathMapValue[k], v, "test GetStringMap() failed")
	}
}

func TestSafeViper_GetStringMapString(t *testing.T) {
	asst := assert.New(t)

	val := GetStringMapString(testServerRouterPathMapKey)
	asst.Equal(len(testServerRouterPathMapValue), len(val),
		"test GetStringMapString() failed")
	for k, v := range val {
		asst.Equal(testServerRouterPathMapValue[k], v, "test GetStringMapString() failed")
	}
}

func TestSafeViper_GetStringMapStringSlice(t *testing.T) {
	asst := assert.New(t)

	val := GetStringMapStringSlice(testServerSwaggerKey)
	asst.Equal(len(testServerSwaggerValue), len(val),
		"test GetStringMapStringSlice() failed")
	for key, valueList := range val {
		asst.Equal(strings.ToLower(testAllowListKey), key, "test GetStringMapStringSlice() failed")
		asst.True(common.ElementEqualInSlice(valueList, testServerSwaggerAllowListValue),
			"test GetStringMapStringSlice() failed")
	}
}

func TestSafeViper_ReloadConfig(t *testing.T) {
	asst := assert.New(t)

	// 热加载状态跟踪
	var (
		callbackTriggered int            // 回调触发计数器
		callbackMutex     sync.Mutex     // 回调计数器锁
		lastAllowList     []string       // 上次读取的配置值
		readerWg          sync.WaitGroup // 读者协程等待组
		writerWg          sync.WaitGroup // 写者协程等待组
	)

	tempConfigStr :=
		`server:
  swagger:
    allowList: ["domain-0.localhost"]
`

	tempFile, err := testCreateTempConfig(tempConfigStr)
	asst.Nil(err, "test ReloadConfig() failed")
	defer func() {
		err = os.Remove(tempFile.Name())
		if err != nil {
			panic(err)
		}
	}()

	SetConfigFile(tempFile.Name())
	err = ReadInConfig()
	asst.Nil(err, "test ReloadConfig() failed")

	OnConfigChange(func(err error) {
		if err != nil {
			panic(err)
			return
		}

		callbackMutex.Lock()
		defer callbackMutex.Unlock()

		callbackTriggered++
	})

	readerErrors := make(chan error, testReaderCount*testConcurrentNum)

	for i := constant.ZeroInt; i < testReaderCount; i++ {
		readerWg.Add(constant.OneInt)
		go func(id int) {
			defer readerWg.Done()
			for j := constant.ZeroInt; j < testConcurrentNum; j++ {
				allowList := GetStringSlice(testServerSwaggerAllowListKey)

				for _, host := range allowList {
					if !strings.Contains(host, "domain") {
						readerErrors <- errors.Errorf("reader %d: invalid allowList %v", id, allowList)
					}
				}

				if j == testConcurrentNum-constant.OneInt {
					callbackMutex.Lock()
					lastAllowList = allowList
					callbackMutex.Unlock()
				}

				// 短暂休眠模拟实际工作负载
				time.Sleep(time.Millisecond * time.Duration(id%10))
			}
		}(i)
	}

	writerWg.Add(constant.OneInt)
	go func() {
		defer writerWg.Done()

		for i := constant.ZeroInt; i < testMaxModifyNum; i++ {
			// update allowList
			host := fmt.Sprintf("domain%d.local", i+constant.OneInt)

			configContent := fmt.Sprintf(
				`server:
  swagger:
    allowList: ["%s"]
`,
				host,
			)

			err = os.WriteFile(tempFile.Name(), []byte(configContent), 0644)
			if err != nil {
				readerErrors <- fmt.Errorf("writer failed: %w", err)
			}

			time.Sleep(200 * time.Millisecond)
		}
	}()

	readerWg.Wait()
	close(readerErrors)

	writerWg.Wait()

	for err = range readerErrors {
		if err != nil {
			t.Fatalf("Concurrent read error: %v", err)
		}
	}

	callbackMutex.Lock()
	defer callbackMutex.Unlock()
	asst.True(callbackTriggered >= testMaxModifyNum, "Expected at least %d callbacks, got %d",
		testMaxModifyNum, callbackTriggered)
	asst.Equal([]string{"domain5.local"}, lastAllowList, "Should have latest config value")
	t.Logf("lastAllowList: %v", lastAllowList)
}
