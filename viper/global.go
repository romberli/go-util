package viper

import "time"

const (
	DefaultFileType = "yaml"
)

var (
	sv *SafeViper
)

func init() {
	sv = NewSafeViper()
}

func SetConfigFile(file string) {
	sv.SetConfigFile(file)
}

func SetConfigType(fileType string) {
	sv.SetConfigType(fileType)
}

func ReadInConfig() error {
	return sv.ReadInConfig()
}

func WatchConfig() {
	sv.WatchConfig()
}

func OnConfigChange(handler func(err error)) {
	sv.OnConfigChange(handler)
}

func SetDefault(key string, value interface{}) {
	sv.SetDefault(key, value)
}

func Set(key string, value interface{}) {
	sv.Set(key, value)
}

func Get(key string) interface{} {
	return sv.Get(key)
}

func GetBool(key string) bool {
	return sv.GetBool(key)
}

func GetInt(key string) int {
	return sv.GetInt(key)
}

func GetString(key string) string {
	return sv.GetString(key)
}

func GetFloat64(key string) float64 {
	return sv.GetFloat64(key)
}

func GetTime(key string) time.Time {
	return sv.GetTime(key)
}

func GetDuration(key string) time.Duration {
	return sv.GetDuration(key)
}

func GetIntSlice(key string) []int {
	return sv.GetIntSlice(key)
}

func GetStringSlice(key string) []string {
	return sv.GetStringSlice(key)
}

func GetStringMap(key string) map[string]interface{} {
	return sv.GetStringMap(key)
}

func GetStringMapString(key string) map[string]string {
	return sv.GetStringMapString(key)
}

func GetStringMapStringSlice(key string) map[string][]string {
	return sv.GetStringMapStringSlice(key)
}
