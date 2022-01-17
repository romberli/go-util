package result

import (
	"fmt"

	"github.com/pingcap/errors"
)

type Map interface {
	GetMap(row, column int) (map[string]interface{}, error)
	GetMapByName(row int, name string) (map[string]interface{}, error)
	GetMapUint(row, column int) (map[string]uint, error)
	GetMapUintByName(row int, name string) (map[string]uint, error)
	GetMapInt(row, column int) (map[string]int, error)
	GetMapIntByName(row int, name string) (map[string]int, error)
	GetMapFloat(row, column int) (map[string]float64, error)
	GetMapFloatByName(row int, name string) (map[string]float64, error)
	GetMapString(row, column int) (map[string]string, error)
	GetMapStringByName(row int, name string) (map[string]string, error)
}

var _ Map = (*EmptyMap)(nil)

type EmptyMap struct {
	MiddlewareType string
}

// NewEmptyMap returns *EmptyMap with given middleware type
func NewEmptyMap(middlewareType string) *EmptyMap {
	return &EmptyMap{middlewareType}
}

// GetMap always returns error, because middleware does not support map data type,
// this function is only for implementing the middleware.Result interface
func (emr *EmptyMap) GetMap(row, column int) (map[string]interface{}, error) {
	return nil, errors.New(fmt.Sprintf("map data type for %s is not supported, never call this function", emr.MiddlewareType))
}

// GetMapByName always returns error, because middleware does not support map data type,
// this function is only for implementing the middleware.Result interface
func (emr *EmptyMap) GetMapByName(row int, name string) (map[string]interface{}, error) {
	return nil, errors.New(fmt.Sprintf("map data type for %s is not supported, never call this function", emr.MiddlewareType))
}

// GetMapUint always returns error, because middleware does not support map data type,
// this function is only for implementing the middleware.Result interface
func (emr *EmptyMap) GetMapUint(row, column int) (map[string]uint, error) {
	return nil, errors.New(fmt.Sprintf("map data type for %s is not supported, never call this function", emr.MiddlewareType))
}

// GetMapUintByName always returns error, because middleware does not support map data type,
// this function is only for implementing the middleware.Result interface
func (emr *EmptyMap) GetMapUintByName(row int, name string) (map[string]uint, error) {
	return nil, errors.New(fmt.Sprintf("map data type for %s is not supported, never call this function", emr.MiddlewareType))
}

// GetMapInt always returns error, because middleware does not support map data type,
// this function is only for implementing the middleware.Result interface
func (emr *EmptyMap) GetMapInt(row, column int) (map[string]int, error) {
	return nil, errors.New(fmt.Sprintf("map data type for %s is not supported, never call this function", emr.MiddlewareType))
}

// GetMapIntByName always returns error, because middleware does not support map data type,
// this function is only for implementing the middleware.Result interface
func (emr *EmptyMap) GetMapIntByName(row int, name string) (map[string]int, error) {
	return nil, errors.New(fmt.Sprintf("map data type for %s is not supported, never call this function", emr.MiddlewareType))
}

// GetMapFloat always returns error, because middleware does not support map data type,
// this function is only for implementing the middleware.Result interface
func (emr *EmptyMap) GetMapFloat(row, column int) (map[string]float64, error) {
	return nil, errors.New(fmt.Sprintf("map data type for %s is not supported, never call this function", emr.MiddlewareType))
}

// GetMapFloatByName always returns error, because middleware does not support map data type,
// this function is only for implementing the middleware.Result interface
func (emr *EmptyMap) GetMapFloatByName(row int, name string) (map[string]float64, error) {
	return nil, errors.New(fmt.Sprintf("map data type for %s is not supported, never call this function", emr.MiddlewareType))
}

// GetMapString always returns error, because middleware does not support map data type,
// this function is only for implementing the middleware.Result interface
func (emr *EmptyMap) GetMapString(row, column int) (map[string]string, error) {
	return nil, errors.New(fmt.Sprintf("map data type for %s is not supported, never call this function", emr.MiddlewareType))
}

// GetMapStringByName always returns error, because middleware does not support map data type,
// this function is only for implementing the middleware.Result interface
func (emr *EmptyMap) GetMapStringByName(row int, name string) (map[string]string, error) {
	return nil, errors.New(fmt.Sprintf("map data type for %s is not supported, never call this function", emr.MiddlewareType))
}
