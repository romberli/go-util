package prometheus

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	client "github.com/prometheus/client_golang/api"
	apiv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/romberli/go-util/constant"
)

const (
	defaultHTTPPrefix  = "http://"
	defaultHTTPSPrefix = "https://"
)

type Config struct {
	client.Config
}

// DefaultRoundTripper is used if no RoundTripper is set in Config,
var DefaultRoundTripper http.RoundTripper = &http.Transport{
	Proxy:               http.ProxyFromEnvironment,
	DialContext:         (&net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
	TLSHandshakeTimeout: 10 * time.Second,
}

// NewConfig returns a new client.Config with given address and round tripper
func NewConfig(addr string, rt http.RoundTripper) Config {
	address := strings.ToLower(addr)
	if !strings.HasPrefix(address, defaultHTTPPrefix) && !strings.HasPrefix(address, defaultHTTPSPrefix) {
		addr = defaultHTTPPrefix + addr
	}

	if rt == nil {
		rt = DefaultRoundTripper
	}

	return Config{
		client.Config{
			Address:      addr,
			RoundTripper: rt,
		},
	}
}

// NewConfigWithDefaultRoundTripper returns a new client.Config with given address and default round tripper
func NewConfigWithDefaultRoundTripper(addr string) Config {
	address := strings.ToLower(addr)
	if !strings.HasPrefix(address, defaultHTTPPrefix) && !strings.HasPrefix(address, defaultHTTPSPrefix) {
		addr = defaultHTTPSPrefix + addr
	}

	return Config{
		client.Config{
			Address:      addr,
			RoundTripper: DefaultRoundTripper,
		},
	}
}

// NewConfigWithBasicAuth returns a new client.Config with given address, user and password
func NewConfigWithBasicAuth(addr, user, pass string) Config {
	address := strings.ToLower(addr)
	if !strings.HasPrefix(address, defaultHTTPPrefix) && !strings.HasPrefix(address, defaultHTTPSPrefix) {
		addr = defaultHTTPPrefix + addr
	}

	return Config{
		client.Config{
			Address:      addr,
			RoundTripper: config.NewBasicAuthRoundTripper(user, config.Secret(pass), constant.EmptyString, DefaultRoundTripper),
		},
	}
}

type Conn struct {
	apiv1.API
}

// NewConn returns a new *Conn with given address and round tripper
func NewConn(addr string, rt http.RoundTripper) (*Conn, error) {
	return NewConnWithConfig(NewConfig(addr, rt))
}

// NewConnWithConfig returns a new *Conn with given config
func NewConnWithConfig(config Config) (*Conn, error) {
	cli, err := client.NewClient(config.Config)
	if err != nil {
		return nil, err
	}

	return &Conn{apiv1.NewAPI(cli)}, nil
}

func (conn *Conn) CheckInstanceStatus() bool {
	query := "1"
	result, err := conn.Execute(query)
	if err != nil {
		return false
	}

	status, err := result.GetInt(constant.ZeroInt, constant.ZeroInt)
	if err != nil {
		return false
	}

	return status == 1
}

// Execute executes given command with arguments and return a result,
// note that args should must be either time.Time or apiv1.Range of prometheus golang client package
func (conn *Conn) Execute(command string, args ...interface{}) (*Result, error) {
	return conn.executeContext(context.Background(), command, args...)
}

// ExecuteContext executes given command with arguments and return a result,
// note that args should must be either time.Time or apiv1.Range of prometheus golang client package
func (conn *Conn) ExecuteContext(ctx context.Context, command string, args ...interface{}) (*Result, error) {
	return conn.executeContext(ctx, command, args...)
}

// executeContext executes given command with arguments and return a result,
// note that args should must be either time.Time or apiv1.Range of prometheus golang client package
func (conn *Conn) executeContext(ctx context.Context, command string, args ...interface{}) (*Result, error) {
	var (
		arg      interface{}
		value    model.Value
		warnings apiv1.Warnings
		err      error
	)

	switch len(args) {
	case 0:
		arg = time.Now()
	case 1:
		arg = args[constant.ZeroInt]
	default:
		return nil, errors.New("two many arguments, argument number should be either 0 or 1")
	}

	switch arg.(type) {
	case time.Time:
		value, warnings, err = conn.Query(ctx, command, arg.(time.Time))
		if err != nil {
			return nil, err
		}
	case apiv1.Range:
		value, warnings, err = conn.QueryRange(ctx, command, arg.(apiv1.Range))
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("args must be either time.Time or apiv1.Range of prometheus golang client package")
	}

	return NewResult(value, warnings), nil
}
