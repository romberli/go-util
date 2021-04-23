package result

type Number interface {
	// GetUint returns uint type value of given row and column number
	GetUint(row, column int) (uint, error)
	// GetUintByName returns uint type value of given row number and column name
	GetUintByName(row int, name string) (uint, error)
	// GetInt returns int type value of given row and column number
	GetInt(row, column int) (int, error)
	// GetIntByName returns int type value of given row number and column name
	GetIntByName(row int, name string) (int, error)
	// GetFloat returns float64 type value of given row and column number
	GetFloat(row, column int) (float64, error)
	// GetFloatByName returns float64 type value of given row number and column name
	GetFloatByName(row int, name string) (float64, error)
}
