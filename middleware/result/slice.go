package result

import (
	"errors"
	"fmt"
)

type Slice interface {
	GetSlice(row, column int) ([]interface{}, error)
	GetSliceByName(row int, name string) ([]interface{}, error)
	GetUintSlice(row, column int) ([]uint, error)
	GetUintSliceByName(row int, name string) ([]uint, error)
	GetIntSlice(row, column int) ([]int, error)
	GetIntSliceByName(row int, name string) ([]int, error)
	GetFloatSlice(row, column int) ([]float64, error)
	GetFloatSliceByName(row int, name string) ([]float64, error)
	GetStringSlice(row, column int) ([]string, error)
	GetStringSliceByName(row int, name string) ([]string, error)
}

var _ Slice = (*EmptySlice)(nil)

type EmptySlice struct {
	MiddlewareType string
}

// NewEmptySlice returns *EmptySlice with given middleware type
func NewEmptySlice(middlewareType string) *EmptySlice {
	return &EmptySlice{middlewareType}
}

// GetSlice always returns error, because middleware does not support slice type
// this function is only for implementing the middleware.Result interface
func (esr *EmptySlice) GetSlice(row, column int) ([]interface{}, error) {
	return nil, errors.New(fmt.Sprintf("GetSlice() for %s is not supported, never call this function", esr.MiddlewareType))
}

// GetSliceByName always returns error, because middleware does not support slice type
// this function is only for implementing the middleware.Result interface
func (esr *EmptySlice) GetSliceByName(row int, name string) ([]interface{}, error) {
	return nil, errors.New(fmt.Sprintf("GetSliceByName() for %s is not supported, never call this function", esr.MiddlewareType))
}

// GetUintSlice always returns error, because middleware does not support slice type
// this function is only for implementing the middleware.Result interface
func (esr *EmptySlice) GetUintSlice(row, column int) ([]uint, error) {
	return nil, errors.New(fmt.Sprintf("GetUintSlice() for %s is not supported, never call this function", esr.MiddlewareType))
}

// GetUintSliceByName always returns error, because middleware does not support slice type
// this function is only for implementing the middleware.Result interface
func (esr *EmptySlice) GetUintSliceByName(row int, name string) ([]uint, error) {
	return nil, errors.New(fmt.Sprintf("GetUintSliceByName() for %s is not supported, never call this function", esr.MiddlewareType))
}

// GetIntSlice always returns error, because middleware does not support slice type
// this function is only for implementing the middleware.Result interface
func (esr *EmptySlice) GetIntSlice(row, column int) ([]int, error) {
	return nil, errors.New(fmt.Sprintf("GetIntSlice() for %s is not supported, never call this function", esr.MiddlewareType))
}

// GetIntSliceByName always returns error, because middleware does not support slice type
// this function is only for implementing the middleware.Result interface
func (esr *EmptySlice) GetIntSliceByName(row int, name string) ([]int, error) {
	return nil, errors.New(fmt.Sprintf("GetIntSliceByName() for %s is not supported, never call this function", esr.MiddlewareType))
}

// GetFloatSlice always returns error, because middleware does not support slice type
// this function is only for implementing the middleware.Result interface
func (esr *EmptySlice) GetFloatSlice(row, column int) ([]float64, error) {
	return nil, errors.New(fmt.Sprintf("GetFloatSlice() for %s is not supported, never call this function", esr.MiddlewareType))
}

// GetFloatSliceByName always returns error, because middleware does not support slice type
// this function is only for implementing the middleware.Result interface
func (esr *EmptySlice) GetFloatSliceByName(row int, name string) ([]float64, error) {
	return nil, errors.New(fmt.Sprintf("GetFloatSliceByName() for %s is not supported, never call this function", esr.MiddlewareType))
}

// GetStringSlice always returns error, because middleware does not support slice type
// this function is only for implementing the middleware.Result interface
func (esr *EmptySlice) GetStringSlice(row, column int) ([]string, error) {
	return nil, errors.New(fmt.Sprintf("GetStringSlice() for %s is not supported, never call this function", esr.MiddlewareType))
}

// GetStringSliceByName always returns error, because middleware does not support slice type
// this function is only for implementing the middleware.Result interface
func (esr *EmptySlice) GetStringSliceByName(row int, name string) ([]string, error) {
	return nil, errors.New(fmt.Sprintf("GetStringSliceByName() for %s is not supported, never call this function", esr.MiddlewareType))
}
