package result

type Number interface {
	// GetUint returns uint type value of given row and column number
	GetUint(row, column int) (uint, error)
	// GetUintByName returns uint type value of given row number and column name
	GetUintByName(row int, name string) (uint, error)
	// GetUint64 returns uint64 type value of given row and column number
	GetUint64(row, column int) (uint64, error)
	// GetUint64ByName returns uint64 type value of given row number and column name
	GetUint64ByName(row int, name string) (uint64, error)
	// GetInt returns int type value of given row and column number
	GetInt(row, column int) (int, error)
	// GetIntByName returns int type value of given row number and column name
	GetIntByName(row int, name string) (int, error)
	// GetInt64 returns int64 type value of given row and column number
	GetInt64(row, column int) (int64, error)
	// GetInt64ByName returns int64 type value of given row number and column name
	GetInt64ByName(row int, name string) (int64, error)
	// GetFloat returns float64 type value of given row and column number
	GetFloat(row, column int) (float64, error)
	// GetFloatByName returns float64 type value of given row number and column name
	GetFloatByName(row int, name string) (float64, error)
}
