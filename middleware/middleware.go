package middleware

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
	GetUint(row, column int) (uint64, error)
	// GetUintByName returns uint64 type value of given row number and column name
	GetUintByName(row int, name string) (uint64, error)
	// GetInt returns int64 type value of given row and column number
	GetInt(row, column int) (int64, error)
	// GetIntByName returns int64 type value of given row number and column name
	GetIntByName(row int, name string) (int64, error)
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
}

type PoolConn interface {
	// Close returns connection back to the pool
	Close() error
	// DisConnect disconnects from the middleware, normally when using connection pool
	DisConnect() error
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
