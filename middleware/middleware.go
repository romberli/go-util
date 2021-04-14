package middleware

import (
	"errors"
	"fmt"
)

type SliceResult interface {
	GetSlice(row, column int) ([]interface{}, error)
	GetSliceByName(row int, name string) ([]interface{}, error)
	GetUintSlice(row, column int) ([]uint64, error)
	GetUintSliceByName(row int, name string) ([]uint64, error)
	GetIntSlice(row, column int) ([]int64, error)
	GetIntSliceByName(row int, name string) ([]int64, error)
	GetFloatSlice(row, column int) ([]float64, error)
	GetFloatSliceByName(row int, name string) ([]float64, error)
	GetStringSlice(row, column int) ([]string, error)
	GetStringSliceByName(row int, name string) ([]string, error)
}

var _ SliceResult = (*EmptySliceResult)(nil)

type EmptySliceResult struct {
	MiddlewareType string
}

// NewEmptySliceResult returns *EmptySliceResult with given middleware type
func NewEmptySliceResult(middlewareType string) *EmptySliceResult {
	return &EmptySliceResult{middlewareType}
}

// GetSlice always returns error, because middleware does not support slice type
// this function is only for implementing the middleware.Result interface
func (esr *EmptySliceResult) GetSlice(row, column int) ([]interface{}, error) {
	return nil, errors.New(fmt.Sprintf("GetSlice() for %s is not supported, never call this function", esr.MiddlewareType))
}

// GetSliceByName always returns error, because middleware does not support slice type
// this function is only for implementing the middleware.Result interface
func (esr *EmptySliceResult) GetSliceByName(row int, name string) ([]interface{}, error) {
	return nil, errors.New(fmt.Sprintf("GetSliceByName() for %s is not supported, never call this function", esr.MiddlewareType))
}

// GetUintSlice always returns error, because middleware does not support slice type
// this function is only for implementing the middleware.Result interface
func (esr *EmptySliceResult) GetUintSlice(row, column int) ([]uint64, error) {
	return nil, errors.New(fmt.Sprintf("GetUintSlice() for %s is not supported, never call this function", esr.MiddlewareType))
}

// GetUintSliceByName always returns error, because middleware does not support slice type
// this function is only for implementing the middleware.Result interface
func (esr *EmptySliceResult) GetUintSliceByName(row int, name string) ([]uint64, error) {
	return nil, errors.New(fmt.Sprintf("GetUintSliceByName() for %s is not supported, never call this function", esr.MiddlewareType))
}

// GetIntSlice always returns error, because middleware does not support slice type
// this function is only for implementing the middleware.Result interface
func (esr *EmptySliceResult) GetIntSlice(row, column int) ([]int64, error) {
	return nil, errors.New(fmt.Sprintf("GetIntSlice() for %s is not supported, never call this function", esr.MiddlewareType))
}

// GetIntSliceByName always returns error, because middleware does not support slice type
// this function is only for implementing the middleware.Result interface
func (esr *EmptySliceResult) GetIntSliceByName(row int, name string) ([]int64, error) {
	return nil, errors.New(fmt.Sprintf("GetIntSliceByName() for %s is not supported, never call this function", esr.MiddlewareType))
}

// GetFloatSlice always returns error, because middleware does not support slice type
// this function is only for implementing the middleware.Result interface
func (esr *EmptySliceResult) GetFloatSlice(row, column int) ([]float64, error) {
	return nil, errors.New(fmt.Sprintf("GetFloatSlice() for %s is not supported, never call this function", esr.MiddlewareType))
}

// GetFloatSliceByName always returns error, because middleware does not support slice type
// this function is only for implementing the middleware.Result interface
func (esr *EmptySliceResult) GetFloatSliceByName(row int, name string) ([]float64, error) {
	return nil, errors.New(fmt.Sprintf("GetFloatSliceByName() for %s is not supported, never call this function", esr.MiddlewareType))
}

// GetStringSlice always returns error, because middleware does not support slice type
// this function is only for implementing the middleware.Result interface
func (esr *EmptySliceResult) GetStringSlice(row, column int) ([]string, error) {
	return nil, errors.New(fmt.Sprintf("GetStringSlice() for %s is not supported, never call this function", esr.MiddlewareType))
}

// GetStringSliceByName always returns error, because middleware does not support slice type
// this function is only for implementing the middleware.Result interface
func (esr *EmptySliceResult) GetStringSliceByName(row int, name string) ([]string, error) {
	return nil, errors.New(fmt.Sprintf("GetStringSliceByName() for %s is not supported, never call this function", esr.MiddlewareType))
}

type MapResult interface {
	GetMap(row, column int) (map[string]interface{}, error)
	GetMapByName(row int, name string) (map[string]interface{}, error)
	GetMapUint(row, column int) (map[string]uint64, error)
	GetMapUintByName(row int, name string) (map[string]uint64, error)
	GetMapInt(row, column int) (map[string]int64, error)
	GetMapIntByName(row int, name string) (map[string]int64, error)
	GetMapFloat(row, column int) (map[string]float64, error)
	GetMapFloatByName(row int, name string) (map[string]float64, error)
	GetMapString(row, column int) (map[string]string, error)
	GetMapStringByName(row int, name string) (map[string]string, error)
}

var _ MapResult = (*EmptyMapResult)(nil)

type EmptyMapResult struct {
	MiddlewareType string
}

// NewEmptyMapResult returns *EmptyMapResult with given middleware type
func NewEmptyMapResult(middlewareType string) *EmptyMapResult {
	return &EmptyMapResult{middlewareType}
}

// GetMap always returns error, because middleware does not support map data type,
// this function is only for implementing the middleware.Result interface
func (emr *EmptyMapResult) GetMap(row, column int) (map[string]interface{}, error) {
	return nil, errors.New(fmt.Sprintf("map data type for %s is not supported, never call this function", emr.MiddlewareType))
}

// GetMapByName always returns error, because middleware does not support map data type,
// this function is only for implementing the middleware.Result interface
func (emr *EmptyMapResult) GetMapByName(row int, name string) (map[string]interface{}, error) {
	return nil, errors.New(fmt.Sprintf("map data type for %s is not supported, never call this function", emr.MiddlewareType))
}

// GetMapUint always returns error, because middleware does not support map data type,
// this function is only for implementing the middleware.Result interface
func (emr *EmptyMapResult) GetMapUint(row, column int) (map[string]uint64, error) {
	return nil, errors.New(fmt.Sprintf("map data type for %s is not supported, never call this function", emr.MiddlewareType))
}

// GetMapUintByName always returns error, because middleware does not support map data type,
// this function is only for implementing the middleware.Result interface
func (emr *EmptyMapResult) GetMapUintByName(row int, name string) (map[string]uint64, error) {
	return nil, errors.New(fmt.Sprintf("map data type for %s is not supported, never call this function", emr.MiddlewareType))
}

// GetMapInt always returns error, because middleware does not support map data type,
// this function is only for implementing the middleware.Result interface
func (emr *EmptyMapResult) GetMapInt(row, column int) (map[string]int64, error) {
	return nil, errors.New(fmt.Sprintf("map data type for %s is not supported, never call this function", emr.MiddlewareType))
}

// GetMapIntByName always returns error, because middleware does not support map data type,
// this function is only for implementing the middleware.Result interface
func (emr *EmptyMapResult) GetMapIntByName(row int, name string) (map[string]int64, error) {
	return nil, errors.New(fmt.Sprintf("map data type for %s is not supported, never call this function", emr.MiddlewareType))
}

// GetMapFloat always returns error, because middleware does not support map data type,
// this function is only for implementing the middleware.Result interface
func (emr *EmptyMapResult) GetMapFloat(row, column int) (map[string]float64, error) {
	return nil, errors.New(fmt.Sprintf("map data type for %s is not supported, never call this function", emr.MiddlewareType))
}

// GetMapFloatByName always returns error, because middleware does not support map data type,
// this function is only for implementing the middleware.Result interface
func (emr *EmptyMapResult) GetMapFloatByName(row int, name string) (map[string]float64, error) {
	return nil, errors.New(fmt.Sprintf("map data type for %s is not supported, never call this function", emr.MiddlewareType))
}

// GetMapString always returns error, because middleware does not support map data type,
// this function is only for implementing the middleware.Result interface
func (emr *EmptyMapResult) GetMapString(row, column int) (map[string]string, error) {
	return nil, errors.New(fmt.Sprintf("map data type for %s is not supported, never call this function", emr.MiddlewareType))
}

// GetMapStringByName always returns error, because middleware does not support map data type,
// this function is only for implementing the middleware.Result interface
func (emr *EmptyMapResult) GetMapStringByName(row int, name string) (map[string]string, error) {
	return nil, errors.New(fmt.Sprintf("map data type for %s is not supported, never call this function", emr.MiddlewareType))
}

type Result interface {
	// LastInsertID returns the database's auto-generated ID
	// after, for example, an INSERT into a table with primary key.
	LastInsertID() (uint64, error)
	// RowsAffected returns the number of rows affected by the query.
	RowsAffected() (uint64, error)
	// RowNumber returns how many rows in the result
	RowNumber() int
	// ColumnNumber return how many columns in the result
	ColumnNumber() int
	// GetValue returns interface{} type value of given row and column number
	GetValue(row, column int) (interface{}, error)
	// ColumnExists check if column exists in the result
	ColumnExists(name string) bool
	// NameIndex returns number of given column
	NameIndex(name string) (int, error)
	// GetValueByName returns interface{} type value of given row number and column name
	GetValueByName(row int, name string) (interface{}, error)
	// IsNull checks if value of given row and column number is nil
	IsNull(row, column int) (bool, error)
	// IsNullByName checks if value of given row number and column name is nil
	IsNullByName(row int, name string) (bool, error)
	// GetUint returns uint64 type value of given row and column number
	GetUint(row, column int) (uint, error)
	// GetUintByName returns uint64 type value of given row number and column name
	GetUintByName(row int, name string) (uint, error)
	// GetInt returns int64 type value of given row and column number
	GetInt(row, column int) (int, error)
	// GetIntByName returns int64 type value of given row number and column name
	GetIntByName(row int, name string) (int, error)
	// GetFloat returns float64 type value of given row and column number
	GetFloat(row, column int) (float64, error)
	// GetFloatByName returns float64 type value of given row number and column name
	GetFloatByName(row int, name string) (float64, error)
	// GetString returns string type value of given row and column number
	GetString(row, column int) (string, error)
	// GetStringByName returns string type value of given row number and column name
	GetStringByName(row int, name string) (string, error)
	// MapToStructSlice maps each row to a struct of the first argument,
	// first argument must be a slice of pointers to structs,
	// each row in the result maps to a struct in the slice,
	// each column in the row maps to a field of the struct,
	// tag argument is the tag of the field, it represents the column name,
	// if there is no such tag in the field, this field will be ignored,
	// so set tag to each field that need to be mapped,
	// using "middleware" as the tag is recommended.
	MapToStructSlice(in interface{}, tag string) error
	// MapToStructByRowIndex maps row of given index result to the struct
	// first argument must be a pointer to struct,
	// each column in the row maps to a field of the struct,
	// tag argument is the tag of the field, it represents the column name,
	// if there is no such tag in the field, this field will be ignored,
	// so set tag to each field that need to be mapped,
	// using "middleware" as the tag is recommended.
	MapToStructByRowIndex(in interface{}, row int, tag string) error
	SliceResult
	MapResult
}

type PoolConn interface {
	// Close returns connection back to the pool
	Close() error
	// Disconnect disconnects from the middleware, normally when using connection pool
	Disconnect() error
	// IsValid validates if connection is valid
	IsValid() bool
	// Execute executes given command and placeholders on the middleware
	Execute(command string, args ...interface{}) (Result, error)
}

type Transaction interface {
	PoolConn
	// Begin begins a transaction
	Begin() error
	// Commit commits current transaction
	Commit() error
	// Rollback rollbacks current transaction
	Rollback() error
}

type Pool interface {
	// Close releases each connection in the pool
	Close() error
	// IsClosed returns if pool had been closed
	IsClosed() bool
	// Get gets a connection from the pool
	Get() (PoolConn, error)
	// Transaction returns a connection that could run multiple statements in the same transaction
	Transaction() (Transaction, error)
	// Supply creates given number of connections and add them to the pool
	Supply(num int) error
	// Release releases given number of connections, each connection will disconnect with the middleware
	Release(num int) error
}
