package mysql

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/romberli/log"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware"
)

const (
	DefaultMaxConnections     = 20
	DefaultInitConnections    = 5
	DefaultMaxIdleConnections = 10
	DefaultMaxIdleTime        = 1800 // Seconds
	DefaultKeepAliveInterval  = 300  // Seconds
	DefaultKeepAliveChunkSize = 5
	DefaultSleepTime          = 1 // Seconds
)

var _ middleware.PoolConn = (*PoolConn)(nil)
var _ middleware.Pool = (*Pool)(nil)

type PoolConfig struct {
	Config
	MaxConnections     int
	InitConnections    int
	MaxIdleConnections int
	MaxIdleTime        int
	KeepAliveInterval  int
}

// NewPoolConfig returns a new PoolConfig
func NewPoolConfig(addr, dbName, dbUser, dbPass string,
	maxConnections, initConnections, maxIdleConnections, maxIdleTime, keepAliveInterval int) PoolConfig {
	config := NewMySQLConfig(addr, dbName, dbUser, dbPass)

	return PoolConfig{
		Config:             config,
		MaxConnections:     maxConnections,
		InitConnections:    initConnections,
		MaxIdleConnections: maxIdleConnections,
		MaxIdleTime:        maxIdleTime,
		KeepAliveInterval:  keepAliveInterval,
	}
}

// NewPoolConfigWithConfig returns a new PoolConfig
func NewPoolConfigWithConfig(config Config, maxConnections, initConnections, maxIdleConnections,
	maxIdleTime, keepAliveInterval int) PoolConfig {
	return PoolConfig{
		Config:             config,
		MaxConnections:     maxConnections,
		InitConnections:    initConnections,
		MaxIdleConnections: maxIdleConnections,
		MaxIdleTime:        maxIdleTime,
		KeepAliveInterval:  keepAliveInterval,
	}
}

// Validate validates pool config
func (cfg *PoolConfig) Validate() (bool, error) {
	// validate MaxConnections
	if cfg.MaxConnections <= 0 {
		return false, errors.New("maximum connection argument should larger than 0")
	}
	// validate InitConnections
	if cfg.InitConnections < 0 {
		return false, errors.New("init connection argument should not be smaller than 0")
	}
	// validate MaxIdleConnections
	if cfg.MaxIdleConnections < 0 {
		return false, errors.New("maximum idle connection argument should not be smaller than 0")
	}
	if cfg.MaxIdleConnections > cfg.MaxConnections {
		return false, errors.New("maximum idle connection argument should not be larger than maximum connection argument")
	}
	// validate MaxIdleTime
	if cfg.MaxIdleTime <= 0 {
		return false, errors.New("maximum idle time argument should be larger than 0")
	}
	// validate KeepAliveInterval
	if cfg.KeepAliveInterval <= 0 {
		return false, errors.New("keep alive interval argument should be larger than 0")
	}

	return true, nil
}

type PoolConn struct {
	*Conn
	Pool *Pool
}

// NewPoolConn returns a new *PoolConn
func NewPoolConn(addr, dbName, dbUser, dbPass string) (*PoolConn, error) {
	conn, err := NewMySQLConn(addr, dbName, dbUser, dbPass)
	if err != nil {
		return nil, err
	}

	return &PoolConn{
		Conn: conn,
		Pool: nil,
	}, nil
}

// NewPoolConnWithPool returns a new *PoolConn
func NewPoolConnWithPool(addr, dbName, dbUser, dbPass string, pool *Pool) (*PoolConn, error) {
	pc, err := NewPoolConn(addr, dbName, dbUser, dbPass)
	if err != nil {
		return nil, err
	}

	if pc.IsValid() {
		// set pool
		pc.Pool = pool
		return pc, nil
	}

	_ = pc.DisConnect()

	return nil, errors.New("new created connection is not valid")
}

// Close returns connection back to the pool
func (pc *PoolConn) Close() error {
	if pc.Pool.isClosed == true || pc.Pool == nil {
		return pc.DisConnect()
	}

	pc.Pool.Lock()
	defer pc.Pool.Unlock()

	pc.Pool.put(pc)

	return nil
}

// DisConnect disconnects from mysql, normally when using connection pool,
// there is no need to disconnect manually, consider to use Close() instead.
func (pc *PoolConn) DisConnect() error {
	pc.Pool = nil
	return pc.Conn.Close()
}

// IsValid validates if connection is valid
func (pc *PoolConn) IsValid() bool {
	return pc.CheckInstanceStatus()
}

// Execute executes given sql and placeholders on the mysql server
func (pc *PoolConn) Execute(command string, args ...interface{}) (middleware.Result, error) {
	r, err := pc.Conn.Execute(command, args...)
	if err != nil {
		return nil, err
	}

	return NewResult(r), err

}

type Pool struct {
	sync.Mutex
	PoolConfig
	freeConnChan    chan *PoolConn
	usedConnections int
	expireTime      time.Time
	keepAliveTime   time.Time
	isClosed        bool
}

// NewMySQLPool returns a new *Pool
func NewMySQLPool(addr, dbName, dbUser, dbPass string,
	maxConnections, initConnections, maxIdleConnections, maxIdleTime, keepAliveInterval int) (*Pool, error) {
	cfg := NewPoolConfig(addr, dbName, dbUser, dbPass, maxConnections, initConnections, maxIdleConnections, maxIdleTime, keepAliveInterval)

	return NewMySQLPoolWithPoolConfig(cfg)
}

// NewMySQLPoolWithDefault returns a new *Pool with default configuration
func NewMySQLPoolWithDefault(addr, dbName, dbUser, dbPass string) (*Pool, error) {
	return NewMySQLPool(addr, dbName, dbUser, dbPass,
		DefaultMaxConnections, DefaultInitConnections, DefaultMaxIdleConnections, DefaultMaxIdleTime, DefaultKeepAliveInterval)
}

// NewMySQLPoolWithConfig returns a new *Pool with a Config object
func NewMySQLPoolWithConfig(config Config, maxConnections, initConnections, maxIdleConnections, maxIdleTime, keepAliveInterval int) (*Pool, error) {
	cfg := NewPoolConfigWithConfig(config, maxConnections, initConnections, maxIdleConnections, maxIdleTime, keepAliveInterval)
	return NewMySQLPoolWithPoolConfig(cfg)
}

// NewMySQLPoolWithPoolConfig returns a new *Pool with a PoolConfig object
func NewMySQLPoolWithPoolConfig(config PoolConfig) (*Pool, error) {
	p := &Pool{
		PoolConfig:      config,
		freeConnChan:    make(chan *PoolConn, config.MaxConnections),
		usedConnections: constant.ZeroInt,
		expireTime:      time.Now().Add(time.Duration(config.MaxIdleTime) * time.Second),
		keepAliveTime:   time.Now().Add(time.Duration(config.KeepAliveInterval) * time.Second),
		isClosed:        false,
	}

	err := p.init()
	if err != nil {
		return nil, err
	}

	return p, nil
}

// init initiates pool
func (p *Pool) init() error {
	// add sufficient connections to the pool
	err := p.supply(p.InitConnections)
	if err != nil {
		return err
	}

	// start a new routine to maintain free connection channel
	go p.maintainFreeChan()

	return nil
}

// UsedConnections returns used connection number
func (p *Pool) UsedConnections() int {
	return p.usedConnections
}

// IsClosed returns if pool had been closed
func (p *Pool) IsClosed() bool {
	return p.isClosed
}

// Supply is an exported alias of supply() function with routine safe
func (p *Pool) Supply(num int) error {
	p.Lock()
	defer p.Unlock()

	return p.supply(num)
}

// supply creates given number of connections and add them to free connection channel
func (p *Pool) supply(num int) error {
	if p.isClosed {
		return nil
	}

	merr := &multierror.Error{}

	for i := 0; i < num; i++ {
		if len(p.freeConnChan)+p.usedConnections < p.MaxConnections {
			pc, err := NewPoolConnWithPool(p.Addr, p.DBName, p.DBUser, p.DBPass, p)
			if err != nil {
				merr = multierror.Append(merr, err)
				continue
			}

			p.addToFreeChan(pc)
		}
	}

	return merr.ErrorOrNil()
}

// Close releases each connection in the pool
func (p *Pool) Close() error {
	p.Lock()
	defer p.Unlock()

	p.isClosed = true
	return p.release(len(p.freeConnChan))
}

// put puts given PoolConn back to the pool
func (p *Pool) put(pc *PoolConn) {
	p.addToFreeChan(pc)
	p.usedConnections--
}

// getFromFreeChan gets a *PoolConn from free connection channel
func (p *Pool) getFromFreeChan() (*PoolConn, bool) {
	pc, ok := <-p.freeConnChan
	return pc, ok
}

// addToFreeChan adds given *PoolConn to free connection channel
func (p *Pool) addToFreeChan(pc *PoolConn) {
	p.freeConnChan <- pc
}

// Get is an exported alias of get() function with routine safe
func (p *Pool) Get() (middleware.PoolConn, error) {
	p.Lock()
	defer p.Unlock()

	return p.get()
}

// get gets a connection from the pool and validate it,
// if there is no valid connection in the pool, it will create a new connection
func (p *Pool) get() (*PoolConn, error) {
	if p.isClosed {
		return nil, errors.New("pool had been closed")
	}

	if p.usedConnections >= p.MaxConnections {
		return nil, errors.New(
			fmt.Sprintf("used connection(%d) had reached maximum connection(%d)",
				p.usedConnections, p.MaxConnections))
	}

	freeChanLen := len(p.freeConnChan)
	// try to get connection from free connection channel
	for i := 0; i < freeChanLen; i++ {
		pc, ok := p.getFromFreeChan()
		// check if connection is still valid
		if ok && pc.IsValid() {
			p.usedConnections++
			return pc, nil
		}

		err := pc.DisConnect()
		if err != nil {
			return nil, fmt.Errorf(
				"disconnecting invalid connection failed when getting connection from the pool.%w",
				err)
		}
	}

	// there is no valid connection in the free connection channel, therefore create a new one
	pc, err := NewPoolConnWithPool(p.Addr, p.DBName, p.DBUser, p.DBPass, p)
	if err != nil {
		return nil, err
	}

	p.usedConnections++
	return pc, nil
}

// maintainFreeChan maintains free connection channel, if there are insufficient connection in the free connection channel,
// it will add some connections, otherwise it will release some.
// for saving disk purpose, if there are errors when maintaining free channel,
// it will log with debug level
func (p *Pool) maintainFreeChan() {
	var (
		err error
		num int
	)

	for {
		if p.isClosed {
			return
		}

		p.Lock()
		now := time.Now()
		// keep alive connections
		if now.After(p.keepAliveTime) {
			p.keepAliveTime = now.Add(time.Duration(p.KeepAliveInterval) * time.Second)
			err = p.keepAlive(DefaultKeepAliveChunkSize)
			if err != nil {
				log.Debugf("got error when keeping alive connections of the pool. total: %d, failed: %d. nested error: %s",
					DefaultKeepAliveChunkSize, err.(*multierror.Error).Len(), err.Error())
			}
		}
		// supply enough connections
		num = p.InitConnections - p.usedConnections - len(p.freeConnChan)
		err = p.supply(num)
		if err != nil {
			log.Debugf("got error when supplying connections to the pool. total: %d, failed: %d. nested error: %s",
				num, err.(*multierror.Error).Len(), err.Error())
		}
		// release excessive connections
		if now.After(p.expireTime) {
			p.expireTime = now.Add(time.Duration(p.MaxIdleTime) * time.Second)
			num = len(p.freeConnChan) + p.usedConnections - p.MaxIdleConnections
			err = p.release(num)
			if err != nil {
				log.Debugf("got error when releasing connections of the pool. total: %d, failed: %d. nested error: %s",
					num, err.(*multierror.Error).Len(), err.Error())
			}
		}

		p.Unlock()
		time.Sleep(DefaultSleepTime * time.Second)
	}
}

// keepAlive keeps alive given number of connections in the pool to avoid disconnected by server side automatically,
// it also checks if connections are still valid
func (p *Pool) keepAlive(num int) error {
	if len(p.freeConnChan) == 0 {
		return nil
	}

	merr := &multierror.Error{}

	for i := 0; i < num; i++ {
		select {
		case pc, ok := <-p.freeConnChan:
			if ok && pc.IsValid() {
				p.addToFreeChan(pc)
				continue
			}

			err := pc.DisConnect()
			if err != nil {
				merr = multierror.Append(merr, err)
			}
		default:
		}
	}

	return merr.ErrorOrNil()
}

// Release is an exported alias of release() function
func (p *Pool) Release(num int) error {
	p.Lock()
	defer p.Unlock()

	return p.release(num)
}

// release releases given number of connections, each connection will disconnect with database
func (p *Pool) release(num int) error {
	merr := &multierror.Error{}

	for i := 0; i < num; i++ {
		if len(p.freeConnChan) == 0 {
			return nil
		}

		// as we didn't lock between get length of free connection channel and get connection from channel,
		// so this is possible to release less than given number of connections, to avoid this,
		// you have to use lock before calling this function.
		// actually, it's not a big deal to release a few less connections,
		// as it will release again at next maintain cycle.
		select {
		case pc, ok := <-p.freeConnChan:
			if ok {
				// disconnect
				err := pc.DisConnect()
				if err != nil {
					merr = multierror.Append(merr, err)
				}
			}
		default:
		}
	}

	return merr.ErrorOrNil()
}
