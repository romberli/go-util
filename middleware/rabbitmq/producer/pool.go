package producer

import (
	"fmt"
	"sync"
	"time"

	"github.com/pingcap/errors"
	"github.com/romberli/go-multierror"
	"github.com/romberli/log"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware/rabbitmq"
	"github.com/romberli/go-util/uid"
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

type PoolProducer struct {
	*Producer
	Pool *Pool
}

// NewPoolProducer returns a new *PoolProducer
func NewPoolProducer(pool *Pool) (*PoolProducer, error) {
	cfg := pool.PoolConfig.Config
	if cfg.Tag != constant.EmptyString {
		cfg = pool.PoolConfig.Config.Clone()
		cfg.Tag = pool.GetFullTag()
	}

	p, err := NewProducerWithConfig(cfg)
	if err != nil {
		return nil, err
	}

	pp := &PoolProducer{
		Producer: p,
		Pool:     pool,
	}

	if pp.IsValid() {
		return pp, nil
	}

	err = pp.Disconnect()
	if err != nil {
		log.Warnf("disconnecting invalid connection failed when creating new pool producer. error:\n%+v", err)
	}

	return nil, errors.New("new created pool producer is not valid")
}

// GetTag returns the tag of the pool producer, should use this method to get the tag instead of using the tag property of the pool config,
// because the tag of the pool config is the prefix of the tag of the pool producer
func (pp *PoolProducer) GetTag() string {
	return pp.Producer.Conn.Config.Tag
}

// Close closes the channel and returns the producer back to the pool
func (pp *PoolProducer) Close() error {
	if pp.Pool.IsClosed() == true || pp.Pool == nil {
		return pp.Disconnect()
	}

	pp.Pool.Lock()
	defer pp.Pool.Unlock()

	err := pp.Producer.Close()
	if err != nil {
		return err
	}
	pp.Pool.put(pp)

	return nil
}

// Disconnect disconnects from rabbitmq, normally when using connection pool,
// there is no need to disconnect manually, consider to use Close() instead.
func (pp *PoolProducer) Disconnect() error {
	pp.Pool = nil

	return pp.Producer.Disconnect()
}

// IsValid validates if connection is valid
func (pp *PoolProducer) IsValid() bool {
	if pp.Conn != nil && !pp.Conn.IsClosed() {
		return true
	}

	return false
}

type Pool struct {
	sync.Mutex
	*rabbitmq.PoolConfig
	uidNode          *uid.Node
	freeProducerChan chan *PoolProducer
	usedConnections  int
	expireTime       time.Time
	keepAliveTime    time.Time
	isClosed         bool
}

// NewPool returns a new *Pool
func NewPool(addr, user, host, vhost, tag, exchange, queue, key string,
	maxConnections, initConnections, maxIdleConnections, maxIdleTime, maxWaitTime, maxRetryCount, keepAliveInterval int) (*Pool, error) {
	cfg := rabbitmq.NewPoolConfig(addr, user, host, vhost, tag, exchange, queue, key,
		maxConnections, initConnections, maxIdleConnections, maxIdleTime, maxWaitTime, maxRetryCount, keepAliveInterval)

	return NewPoolWithPoolConfig(cfg)
}

// NewPoolWithDefault returns a new *Pool with default configuration
func NewPoolWithDefault(addr, user, host, vhost, tag, exchange, queue, key string) (*Pool, error) {
	return NewPool(addr, user, host, vhost, tag, exchange, queue, key,
		DefaultMaxConnections, DefaultInitConnections, DefaultMaxIdleConnections,
		DefaultMaxIdleTime, DefaultMaxWaitTime, DefaultMaxRetryCount, DefaultKeepAliveInterval)
}

// NewPoolWithConfig returns a new *Pool with a Config object
func NewPoolWithConfig(config *rabbitmq.Config, maxConnections, initConnections,
	maxIdleConnections, maxIdleTime, maxWaitTime, maxRetryCount, keepAliveInterval int) (*Pool, error) {
	cfg := rabbitmq.NewPoolConfigWithConfig(config, maxConnections, initConnections,
		maxIdleConnections, maxIdleTime, maxWaitTime, maxRetryCount, keepAliveInterval)

	return NewPoolWithPoolConfig(cfg)
}

// NewPoolWithPoolConfig returns a new *Pool with a PoolConfig object
func NewPoolWithPoolConfig(config *rabbitmq.PoolConfig) (*Pool, error) {
	// validate config
	err := config.Validate()
	if err != nil {
		return nil, err
	}

	var uidNode *uid.Node

	if config.Tag != constant.EmptyString {
		uidNode, err = uid.NewNodeWithDefault()
		if err != nil {
			return nil, err
		}
	}

	p := &Pool{
		PoolConfig:       config,
		uidNode:          uidNode,
		freeProducerChan: make(chan *PoolProducer, config.MaxConnections*DefaultFreeChanLengthTimes),
		usedConnections:  constant.ZeroInt,
		expireTime:       time.Now().Add(time.Duration(config.MaxIdleTime) * time.Second),
		keepAliveTime:    time.Now().Add(time.Duration(config.KeepAliveInterval) * time.Second),
		isClosed:         false,
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

// GetFullTag gets the full tag
func (p *Pool) GetFullTag() string {
	var ft string

	if p.PoolConfig.Config.Tag != constant.EmptyString {
		ft = fmt.Sprintf("%s-%s", p.PoolConfig.Config.Tag, p.uidNode.Generate().Base36())
	}

	return ft
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

	for i := constant.ZeroInt; i < num; i++ {
		if len(p.freeProducerChan)+p.usedConnections < p.MaxConnections {
			pc, err := NewPoolProducer(p)
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
	return p.release(len(p.freeProducerChan))
}

// put puts given PoolProducer back to the pool
func (p *Pool) put(pc *PoolProducer) {
	p.addToFreeChan(pc)
	p.usedConnections--
}

// getFromFreeChan gets a *PoolProducer from free connection channel
func (p *Pool) getFromFreeChan() (*PoolProducer, bool) {
	pc, ok := <-p.freeProducerChan

	return pc, ok
}

// addToFreeChan adds given *PoolProducer to free connection channel
func (p *Pool) addToFreeChan(pc *PoolProducer) {
	p.freeProducerChan <- pc
}

// Get is an exported alias of get() function with routine safe
func (p *Pool) Get() (*PoolProducer, error) {
	return p.getFromPool()
}

// getFromPool gets a connection from the pool
func (p *Pool) getFromPool() (*PoolProducer, error) {
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
func (p *Pool) get() (*PoolProducer, error) {
	if p.IsClosed() {
		return nil, errors.New("pool had been closed")
	}

	if p.usedConnections >= p.MaxConnections {
		return nil, errors.Errorf("used connections had reached maximum connections. used_connections: %d, max_connections: %d",
			p.usedConnections, p.MaxConnections)
	}

	freeChanLen := len(p.freeProducerChan)
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
	pc, err := NewPoolProducer(p)
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
			num := p.InitConnections - len(p.freeProducerChan)
			err := p.supply(num)
			if err != nil {
				log.Debugf("got error when supplying connections to the pool. total: %d, failed: %d. nested error:\n%+v",
					num, err.(*multierror.Error).Len(), err)
			}
		}
		// release excessive connections
		if now.After(p.expireTime) {
			p.expireTime = now.Add(time.Duration(p.MaxIdleTime) * time.Second)
			num := len(p.freeProducerChan) - p.MaxIdleConnections
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
	if len(p.freeProducerChan) == constant.ZeroInt {
		return nil
	}

	merr := &multierror.Error{}

	for i := 0; i < num; i++ {
		select {
		case pc, ok := <-p.freeProducerChan:
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
		if len(p.freeProducerChan) == constant.ZeroInt {
			return nil
		}

		// as we didn't lock between get length of free connection channel and get connection from channel,
		// so this is possible to release less than given number of connections, to avoid this,
		// you have to use lock before calling this function.
		// actually, it's not a big deal to release a few less connections,
		// as it will release again at next maintain cycle.
		select {
		case pc, ok := <-p.freeProducerChan:
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
