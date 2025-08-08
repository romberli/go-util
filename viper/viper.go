package viper

import (
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type SafeViper struct {
	mutex sync.RWMutex
	viper *viper.Viper
}

// NewSafeViper returns a new *SafeViper
func NewSafeViper() *SafeViper {
	return &SafeViper{
		viper: viper.New(),
	}
}

// SetConfigFile sets the config file
func (sv *SafeViper) SetConfigFile(file string) {
	sv.mutex.Lock()
	defer sv.mutex.Unlock()

	sv.viper.SetConfigFile(file)
}

// SetConfigType sets the type of the config file
func (sv *SafeViper) SetConfigType(fileType string) {
	sv.mutex.Lock()
	defer sv.mutex.Unlock()

	sv.viper.SetConfigType(fileType)
}

// SetConfigName sets the name of the config file
func (sv *SafeViper) ReadInConfig() error {
	sv.mutex.Lock()
	defer sv.mutex.Unlock()

	return sv.viper.ReadInConfig()
}

// WatchConfig sets the config file
func (sv *SafeViper) WatchConfig() {
	sv.mutex.Lock()
	defer sv.mutex.Unlock()

	sv.viper.WatchConfig()
}

// OnConfigChange sets the config file while the config file changes
func (sv *SafeViper) OnConfigChange(handler func(error)) {
	sv.viper.OnConfigChange(func(event fsnotify.Event) {
		sv.mutex.Lock()
		defer sv.mutex.Unlock()

		// 重新加载配置
		err := sv.viper.ReadInConfig()

		handler(err)
	})

	sv.WatchConfig()
}

func (sv *SafeViper) SetDefault(key string, value interface{}) {
	sv.mutex.Lock()
	defer sv.mutex.Unlock()

	sv.viper.SetDefault(key, value)
}

func (sv *SafeViper) Set(key string, value interface{}) {
	sv.mutex.Lock()
	defer sv.mutex.Unlock()

	sv.viper.Set(key, value)
}

func (sv *SafeViper) Get(key string) interface{} {
	sv.mutex.RLock()
	defer sv.mutex.RUnlock()

	return sv.viper.Get(key)
}

func (sv *SafeViper) GetBool(key string) bool {
	sv.mutex.RLock()
	defer sv.mutex.RUnlock()

	return sv.viper.GetBool(key)
}

func (sv *SafeViper) GetInt(key string) int {
	sv.mutex.RLock()
	defer sv.mutex.RUnlock()

	return sv.viper.GetInt(key)
}

func (sv *SafeViper) GetString(key string) string {
	sv.mutex.RLock()
	defer sv.mutex.RUnlock()

	return sv.viper.GetString(key)
}

func (sv *SafeViper) GetFloat64(key string) float64 {
	sv.mutex.RLock()
	defer sv.mutex.RUnlock()

	return sv.viper.GetFloat64(key)
}

func (sv *SafeViper) GetTime(key string) time.Time {
	sv.mutex.RLock()
	defer sv.mutex.RUnlock()

	return sv.viper.GetTime(key)
}

func (sv *SafeViper) GetDuration(key string) time.Duration {
	sv.mutex.RLock()
	defer sv.mutex.RUnlock()

	return sv.viper.GetDuration(key)
}

func (sv *SafeViper) GetIntSlice(key string) []int {
	sv.mutex.RLock()
	defer sv.mutex.RUnlock()

	return sv.viper.GetIntSlice(key)
}

func (sv *SafeViper) GetStringSlice(key string) []string {
	sv.mutex.RLock()
	defer sv.mutex.RUnlock()

	return sv.viper.GetStringSlice(key)
}

func (sv *SafeViper) GetStringMap(key string) map[string]interface{} {
	sv.mutex.RLock()
	defer sv.mutex.RUnlock()

	return sv.viper.GetStringMap(key)
}

func (sv *SafeViper) GetStringMapString(key string) map[string]string {
	sv.mutex.RLock()
	defer sv.mutex.RUnlock()

	return sv.viper.GetStringMapString(key)
}

func (sv *SafeViper) GetStringMapStringSlice(key string) map[string][]string {
	sv.mutex.RLock()
	defer sv.mutex.RUnlock()

	return sv.viper.GetStringMapStringSlice(key)
}
