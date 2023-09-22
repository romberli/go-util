package mysql

import (
	"context"
	"sync"
	"time"

	"github.com/pingcap/errors"
	"github.com/romberli/go-multierror"
	"github.com/romberli/log"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware"
)

const (
	DefaultMaxConnections      = 20
	DefaultInitConnections     = 5
	DefaultMaxIdleConnections  = 10
	DefaultMaxIdleTime         = 1800 // seconds
	DefaultMaxWaitTime         = 1    // seconds
	DefaultMaxRetryCount       = -1
	DefaultKeepAliveInterval   = 300 // seconds
	DefaultKeepAliveChunkSize  = 5
	DefaultSleepTime           = 1 // seconds
	DefaultFreeChanLengthTimes = 2

	DefaultUnlimitedWaitTime   = -1 // seconds
	DefaultUnlimitedRetryCount = -1
	DefaultDelayTime           = 5 // milliseconds
)

var _ middleware.PoolConn = (*PoolConn)(nil)
var _ middleware.Pool = (*Pool)(nil)

type PoolConfig struct {
	Config
	MaxConnections     int
	InitConnections    int
	MaxIdleConnections int
	MaxIdleTime        int
	MaxWaitTime        int
	MaxRetryCount      int
	KeepAliveInterval  int
}

// NewPoolConfig returns a new PoolConfig
func NewPoolConfig(addr, dbName, dbUser, dbPass string,
	maxConnections, initConnections, maxIdleConnections, maxIdleTime, maxWaitTime, maxRetryCount, keepAliveInterval int) PoolConfig {
	config := NewConfig(addr, dbName, dbUser, dbPass)

	return PoolConfig{
		Config:             config,
		MaxConnections:     maxConnections,
		InitConnections:    initConnections,
		MaxIdleConnections: maxIdleConnections,
		MaxIdleTime:        maxIdleTime,
		MaxWaitTime:        maxWaitTime,
		MaxRetryCount:      maxRetryCount,
		KeepAliveInterval:  keepAliveInterval,
	}
}

// NewPoolConfigWithConfig returns a new PoolConfig
func NewPoolConfigWithConfig(config Config, maxConnections, initConnections, maxIdleConnections,
	maxIdleTime, maxWaitTime, maxRetryCount, keepAliveInterval int) PoolConfig {
	return PoolConfig{
		Config:             config,
		MaxConnections:     maxConnections,
		InitConnections:    initConnections,
		MaxIdleConnections: maxIdleConnections,
		MaxIdleTime:        maxIdleTime,
		MaxWaitTime:        maxWaitTime,
		MaxRetryCount:      maxRetryCount,
		KeepAliveInterval:  keepAliveInterval,
	}
}

// Validate validates pool config
func (cfg *PoolConfig) Validate() error {
	// validate MaxConnections
	if cfg.MaxConnections <= constant.ZeroInt {
		return errors.New("maximum connection argument should larger than 0")
	}
	// validate InitConnections
	if cfg.InitConnections < constant.ZeroInt {
		return errors.New("init connection argument should not be smaller than 0")
	}
	if cfg.InitConnections > cfg.MaxConnections {
		return errors.Errorf("init connections should be less or equal than maximum connections. init_connections: %d, max_connections: %d",
			cfg.InitConnections, cfg.MaxConnections)
	}
	// validate MaxIdleConnections
	if cfg.MaxIdleConnections < constant.ZeroInt {
		return errors.New("maximum idle connection argument should not be smaller than 0")
	}
	if cfg.MaxIdleConnections > cfg.MaxConnections {
		return errors.New("maximum idle connection argument should not be larger than maximum connection argument")
	}
	// validate MaxIdleTime
	if cfg.MaxIdleTime <= constant.ZeroInt {
		return errors.New("maximum idle time argument should be larger than 0")
	}
	// validate MaxWaitTime
	if cfg.MaxWaitTime < DefaultUnlimitedWaitTime {
		return errors.New("maximum wait time argument should not be smaller than -1")
	}
	// validate MaxRetryCount
	if cfg.MaxRetryCount < DefaultUnlimitedRetryCount {
		return errors.New("maximum retry count argument should not be smaller than -1")
	}
	// validate KeepAliveInterval
	if cfg.KeepAliveInterval <= constant.ZeroInt {
		return errors.New("keep alive interval argument should be larger than 0")
	}

	return nil
}

type PoolConn struct {
	*Conn
	Pool *Pool
}

// NewPoolConn returns a new *PoolConn
func NewPoolConn(addr, dbName, dbUser, dbPass string) (*PoolConn, error) {
	conn, err := NewConn(addr, dbName, dbUser, dbPass)
	if err != nil {
		return nil, err
	}

	return &PoolConn{
		Conn: conn,
		Pool: nil,
	}, nil
}

// NewPoolConnWithPool returns a new *PoolConn
func NewPoolConnWithPool(pool *Pool, addr, dbName, dbUser, dbPass string) (*PoolConn, error) {
	pc, err := NewPoolConn(addr, dbName, dbUser, dbPass)
	if err != nil {
		return nil, err
	}

	if pc.IsValid() {
		// set pool
		pc.Pool = pool
		return pc, nil
	}

	if pc != nil {
		_ = pc.Disconnect()
	}

	return nil, errors.New("new created connection is not valid")
}

// Close returns connection back to the pool
func (pc *PoolConn) Close() error {
	if pc.Pool.IsClosed() == true || pc.Pool == nil {
		return pc.Disconnect()
	}

	pc.Pool.Lock()
	defer pc.Pool.Unlock()

	err := pc.Rollback()
	if err != nil {
		merr := &multierror.Error{}
		merr = multierror.Append(merr, err)

		disconnectErr := pc.Disconnect()
		if disconnectErr != nil {
			merr = multierror.Append(merr, disconnectErr)
		}

		return merr.ErrorOrNil()
	}

	pc.Pool.put(pc)

	return nil
}

// Disconnect disconnects from mysql, normally when using connection pool,
// there is no need to disconnect manually, consider to use Close() instead.
func (pc *PoolConn) Disconnect() error {
	pc.Pool = nil

	err := pc.Rollback()
	if err != nil {
		return err
	}

	return pc.Conn.Close()
}

// IsValid validates if connection is valid
func (pc *PoolConn) IsValid() bool {
	return pc.CheckInstanceStatus()
}

// Prepare prepares a statement and returns a *Statement
func (pc *PoolConn) Prepare(command string) (middleware.Statement, error) {
	return pc.prepareContext(context.Background(), command)
}

// PrepareContext prepares a statement with context and returns a *Statement,
// note that PrepareContext does not use context, it's only for implementing some interfaces
func (pc *PoolConn) PrepareContext(ctx context.Context, command string) (middleware.Statement, error) {
	return pc.prepareContext(ctx, command)
}

// prepareContext prepares a statement with context and returns a *Statement,
// note that prepareContext does not use context, it's only for implementing some interfaces
func (pc *PoolConn) prepareContext(ctx context.Context, command string) (middleware.Statement, error) {
	return pc.Conn.PrepareContext(ctx, command)
}

// Execute executes given sql and placeholders on the mysql server
func (pc *PoolConn) Execute(command string, args ...interface{}) (middleware.Result, error) {
	return pc.executeContext(context.Background(), command, args...)
}

// ExecuteContext executes given sql and placeholders on the mysql server
// note that ExecuteContext does not use context, it's only for implementing some interfaces
func (pc *PoolConn) ExecuteContext(ctx context.Context, command string, args ...interface{}) (middleware.Result, error) {
	return pc.executeContext(ctx, command, args...)
}

// Execute executes given sql and placeholders on the mysql server
// note that executeContext does not use context, it's only for implementing some interfaces
func (pc *PoolConn) executeContext(ctx context.Context, command string, args ...interface{}) (middleware.Result, error) {
	return pc.Conn.ExecuteContext(ctx, command, args...)
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

// NewPool returns a new *Pool
func NewPool(addr, dbName, dbUser, dbPass string,
	maxConnections, initConnections, maxIdleConnections, maxIdleTime, maxWaitTime, maxRetryCount, keepAliveInterval int) (*Pool, error) {
	cfg := NewPoolConfig(addr, dbName, dbUser, dbPass, maxConnections, initConnections, maxIdleConnections,
		maxIdleTime, maxWaitTime, maxRetryCount, keepAliveInterval)

	return NewPoolWithPoolConfig(cfg)
}

// NewPoolWithDefault returns a new *Pool with default configuration
func NewPoolWithDefault(addr, dbName, dbUser, dbPass string) (*Pool, error) {
	return NewPool(addr, dbName, dbUser, dbPass,
		DefaultMaxConnections, DefaultInitConnections, DefaultMaxIdleConnections,
		DefaultMaxIdleTime, DefaultMaxWaitTime, DefaultMaxRetryCount, DefaultKeepAliveInterval)
}

// NewPoolWithConfig returns a new *Pool with a Config object
func NewPoolWithConfig(config Config, maxConnections, initConnections,
	maxIdleConnections, maxIdleTime, maxWaitTime, maxRetryCount, keepAliveInterval int) (*Pool, error) {
	cfg := NewPoolConfigWithConfig(config, maxConnections, initConnections,
		maxIdleConnections, maxIdleTime, maxWaitTime, maxRetryCount, keepAliveInterval)

	return NewPoolWithPoolConfig(cfg)
}

// NewPoolWithPoolConfig returns a new *Pool with a PoolConfig object
func NewPoolWithPoolConfig(config PoolConfig) (*Pool, error) {
	// validate config
	err := config.Validate()
	if err != nil {
		return nil, err
	}

	p := &Pool{
		PoolConfig:      config,
		freeConnChan:    make(chan *PoolConn, config.MaxConnections*DefaultFreeChanLengthTimes),
		usedConnections: constant.ZeroInt,
		expireTime:      time.Now().Add(time.Duration(config.MaxIdleTime) * time.Second),
		keepAliveTime:   time.Now().Add(time.Duration(config.KeepAliveInterval) * time.Second),
		isClosed:        false,
	}

	err = p.init()
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
	if p.IsClosed() {
		return nil
	}

	if num <= constant.ZeroInt {
		return nil
	}

	merr := &multierror.Error{}

	for i := 0; i < num; i++ {
		if len(p.freeConnChan)+p.usedConnections < p.MaxConnections {
			pc, err := NewPoolConnWithPool(p, p.Addr, p.DBName, p.DBUser, p.DBPass)
			if err != nil {
				merr = multierror.Append(merr, err)
				continue
			}

			p.addToFreeChan(pc)
		}
	}

	return errors.Trace(merr.ErrorOrNil())
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
	return p.getFromPool()
}

// Transaction simply returns *PoolConn, because it had implemented Transaction interface
func (p *Pool) Transaction() (middleware.Transaction, error) {
	return p.getFromPool()
}

// getFromPool gets a connection from the pool
func (p *Pool) getFromPool() (*PoolConn, error) {
	maxWaitTime := p.MaxWaitTime
	if maxWaitTime < constant.ZeroInt {
		maxWaitTime = int(constant.Century.Seconds())
	}

	timer := time.NewTimer(time.Duration(maxWaitTime) * time.Second)
	defer timer.Stop()

	var i int

	for {
		p.Lock()

		pc, err := p.get()
		p.Unlock()

		if err != nil {
			// check retry count
			if p.MaxRetryCount >= constant.ZeroInt && i >= p.MaxRetryCount {
				return nil, err
			}

			// check wait time
			select {
			case <-timer.C:
				return nil, err
			default:
				time.Sleep(time.Duration(DefaultDelayTime) * time.Millisecond)
			}

			i++
			continue
		}

		return pc, nil
	}
}

// get gets a connection from the pool and validate it,
// if there is no valid connection in the pool, it will create a new connection
func (p *Pool) get() (*PoolConn, error) {
	if p.IsClosed() {
		return nil, errors.New("pool had been closed")
	}

	if p.usedConnections >= p.MaxConnections {
		return nil, errors.Errorf("used connections had reached maximum connections. used_connections: %d, max_connections: %d",
			p.usedConnections, p.MaxConnections)
	}

	freeChanLen := len(p.freeConnChan)
	// try to get connection from free connection channel
	for i := constant.ZeroInt; i < freeChanLen; i++ {
		pc, ok := p.getFromFreeChan()
		if ok {
			if pc == nil {
				continue
			}
			// check if connection is still valid
			if pc.IsValid() {
				p.usedConnections++
				return pc, nil
			}

			err := pc.Disconnect()
			if err != nil {
				log.Warnf("disconnecting invalid connection failed when getting connection from the pool. error:\n%+v", err)
			}
		}

		log.Errorf("getting connection from free channel failed, this should not happen cause we only getting from free chan when it has free connection. free channel length: %d", freeChanLen)
	}

	// there is no valid connection in the free connection channel, therefore create a new one
	pc, err := NewPoolConnWithPool(p, p.Addr, p.DBName, p.DBUser, p.DBPass)
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
	for {
		if p.IsClosed() {
			return
		}

		p.Lock()
		now := time.Now()

		// keep alive connections
		if now.After(p.keepAliveTime) {
			p.keepAliveTime = now.Add(time.Duration(p.KeepAliveInterval) * time.Second)
			err := p.keepAlive(DefaultKeepAliveChunkSize)
			if err != nil {
				log.Debugf("got error when keeping alive connections of the pool. total: %d, failed: %d. nested error:\n%+v",
					DefaultKeepAliveChunkSize, err.(*multierror.Error).Len(), err)
			}
		}
		// supply enough connections
		if p.InitConnections+p.usedConnections <= p.MaxConnections {
			num := p.InitConnections - len(p.freeConnChan)
			err := p.supply(num)
			if err != nil {
				log.Debugf("got error when supplying connections to the pool. total: %d, failed: %d. nested error:\n%+v",
					num, err.(*multierror.Error).Len(), err)
			}
		}
		// release excessive connections
		if now.After(p.expireTime) {
			p.expireTime = now.Add(time.Duration(p.MaxIdleTime) * time.Second)
			num := len(p.freeConnChan) - p.MaxIdleConnections
			err := p.release(num)
			if err != nil {
				log.Debugf("got error when releasing connections of the pool. total: %d, failed: %d. nested error:\n%+v",
					num, err.(*multierror.Error).Len(), err)
			}
		}
		p.Unlock()
		time.Sleep(DefaultSleepTime * time.Second)
	}
}

// keepAlive keeps alive given number of connections in the pool to avoid disconnected by server side automatically,
// it also checks if connections are still valid
func (p *Pool) keepAlive(num int) error {
	if len(p.freeConnChan) == constant.ZeroInt {
		return nil
	}

	merr := &multierror.Error{}

	for i := 0; i < num; i++ {
		select {
		case pc, ok := <-p.freeConnChan:
			if ok {
				if pc == nil {
					continue
				}
				// check if connection is still valid
				if pc.IsValid() {
					p.addToFreeChan(pc)
					continue
				}

				err := pc.Disconnect()
				if err != nil {
					merr = multierror.Append(merr, err)
				}
			}
		default:
		}
	}

	return errors.Trace(merr.ErrorOrNil())
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

	for i := constant.ZeroInt; i < num; i++ {
		if len(p.freeConnChan) == constant.ZeroInt {
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
				err := pc.Disconnect()
				if err != nil {
					merr = multierror.Append(merr, err)
				}
			}
		default:
		}
	}

	return errors.Trace(merr.ErrorOrNil())
}
