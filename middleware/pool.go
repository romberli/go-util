package middleware

type PoolConn interface {
	Close() error
	DisConnect() error
	IsValid() bool
	Execute(command string, args ...interface{}) (interface{}, error)
}

type Pool interface {
	Close() error
	IsClosed() bool
	Get() (PoolConn, error)
	Supply(num int) error
	Release(num int) error
}
