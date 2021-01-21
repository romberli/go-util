package middleware

type Result interface {
	RowNumber() int
	ColumnNumber() int
	GetValue(row, column int) (interface{}, error)
	NameIndex(name string) (int, error)
	GetValueByName(row int, name string) (interface{}, error)
	IsNull(row, column int) (bool, error)
	IsNullByName(row int, name string) (bool, error)
	GetUint(row, column int) (uint64, error)
	GetUintByName(row int, name string) (uint64, error)
	GetInt(row, column int) (int64, error)
	GetIntByName(row int, name string) (int64, error)
	GetFloat(row, column int) (float64, error)
	GetFloatByName(row int, name string) (float64, error)
	GetString(row, column int) (string, error)
	GetStringByName(row int, name string) (string, error)
}

type PoolConn interface {
	Close() error
	DisConnect() error
	IsValid() bool
	Execute(command string, args ...interface{}) (Result, error)
}

type Pool interface {
	Close() error
	IsClosed() bool
	Get() (PoolConn, error)
	Supply(num int) error
	Release(num int) error
}
