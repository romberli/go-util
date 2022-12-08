package http

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/buger/jsonparser"
	"github.com/pingcap/errors"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/log"
)

const (
	StatusOK                  = http.StatusOK
	StatusInternalServerError = http.StatusInternalServerError

	defaultResponseCodeJSON = "code"

	defaultClientTimeout         = 60 * time.Second
	defaultDialTimeout           = 30 * time.Second
	defaultKeepAlive             = 30 * time.Second
	defaultTLSHandshakeTimeout   = 10 * time.Second
	defaultContentType           = "application/json"
	defaultMaxIdleConns          = 100
	defaultIdleConnTimeout       = 90 * time.Second
	defaultExpectContinueTimeout = 1 * time.Second
	defaultMaxIdleConnsPerHost   = 20

	DefaultMaxWaitTime   = 60 // seconds
	DefaultMaxRetryCount = 3
	DefaultDelay         = 10 // milliseconds

	DefaultUnlimitedWaitTime   = -1 // seconds
	DefaultUnlimitedRetryCount = -1
)

var (
	DefaultTransport = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           (&net.Dialer{Timeout: defaultDialTimeout, KeepAlive: defaultKeepAlive}).DialContext,
		TLSHandshakeTimeout:   defaultTLSHandshakeTimeout,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          defaultMaxIdleConns,
		IdleConnTimeout:       defaultIdleConnTimeout,
		ExpectContinueTimeout: defaultExpectContinueTimeout,
		MaxIdleConnsPerHost:   defaultMaxIdleConnsPerHost,
	}
)

type Client struct {
	client *http.Client
}

func NewClient(client *http.Client) *Client {
	return newClient(client)
}

func NewClientWithDefault() *Client {
	client := &http.Client{
		Transport: DefaultTransport,
		Timeout:   defaultClientTimeout,
	}

	return newClient(client)
}

func newClient(client *http.Client) *Client {
	c := &Client{
		client: client,
	}

	return c
}

func (c *Client) GetClient() *http.Client {
	return c.client
}

func (c *Client) Close() {
	c.client.CloseIdleConnections()
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func (c *Client) SetMaxIdleConns(maxIdleConns int) {
	c.client.Transport.(*http.Transport).MaxIdleConns = maxIdleConns
}

func (c *Client) SetMaxIdleConnsPerHost(maxIdleConnsPerHost int) {
	c.client.Transport.(*http.Transport).MaxIdleConnsPerHost = maxIdleConnsPerHost
}

func (c *Client) Get(url string) (*http.Response, error) {
	return c.client.Get(PrepareURL(url))
}

func (c *Client) GetWithRetry(url string, maxWaitTime, maxRetryCount, delay int) (*http.Response, error) {
	// validate retry options
	err := c.validate(maxWaitTime, maxRetryCount, delay)
	if err != nil {
		return nil, err
	}

	maxWait := maxWaitTime
	if maxWait < constant.ZeroInt {
		maxWait = int(constant.Century.Seconds())
	}
	timeoutChan := time.After(time.Duration(maxWait) * time.Second)

	var (
		i    int
		resp *http.Response
	)

	for {
		resp, err = c.Get(url)
		if err != nil {
			if maxRetryCount >= constant.ZeroInt && i >= maxRetryCount {
				return resp, err
			}

			i++
			// check for timeout
			select {
			case <-timeoutChan:
				return resp, err
			default:
				time.Sleep(time.Duration(delay) * time.Millisecond)
			}
			continue
		}

		return resp, nil
	}
}

func (c *Client) Post(url string, body []byte) (*http.Response, error) {
	return c.postJSON(url, body)
}

func (c *Client) PostWithRetry(url string, body []byte, maxWaitTime, maxRetryCount, delay int) (*http.Response, error) {
	// validate retry options
	err := c.validate(maxWaitTime, maxRetryCount, delay)
	if err != nil {
		return nil, err
	}

	maxWait := maxWaitTime
	if maxWait < constant.ZeroInt {
		maxWait = int(constant.Century.Seconds())
	}
	timeoutChan := time.After(time.Duration(maxWait) * time.Second)

	var (
		i    int
		resp *http.Response
	)

	for {
		resp, err = c.Post(url, body)
		if err != nil {
			if maxRetryCount >= constant.ZeroInt && i >= maxRetryCount {
				return resp, err
			}

			i++
			// check for timeout
			select {
			case <-timeoutChan:
				return resp, err
			default:
				time.Sleep(time.Duration(delay) * time.Millisecond)
			}
			continue
		}

		return resp, nil
	}
}

func (c *Client) postJSON(url string, body []byte) (*http.Response, error) {
	resp, err := c.client.Post(PrepareURL(url), defaultContentType, bytes.NewBuffer(body))
	if err != nil {
		return resp, errors.Errorf("http Client.postJSON(): http post failed. url: %s, body: %s, error:\n%+v", url, string(body), err)
	}

	return resp, nil
}

func (c *Client) PostDAS(url string, body []byte, maxWaitTime, maxRetryCount, delay int) ([]byte, error) {
	resp, err := c.PostWithRetry(url, body, maxWaitTime, maxRetryCount, delay)
	if err != nil {
		return nil, err
	}
	// read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Errorf("http Client.PostAndParse(): http response body failed. error:\n%+v", errors.Trace(err))
		}
	}()
	// check http status code
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("got wrong http response status code. status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	code, err := jsonparser.GetInt(respBody, defaultResponseCodeJSON)
	if err != nil && err != jsonparser.KeyPathNotFoundError {
		return nil, errors.Errorf("got error when getting code field from response body. error:\n%+v", errors.Trace(err))
	}
	if code != constant.ZeroInt {
		return nil, errors.Errorf("code field in response body is not 0. code: %d, body: %s", code, string(respBody))
	}

	return respBody, nil
}

func (c *Client) validate(maxWaitTime, maxRetryCount, delay int) error {
	// validate maxWaitTime
	if maxWaitTime < DefaultUnlimitedWaitTime {
		return errors.New("maximum wait time argument should not be smaller than -1")
	}
	// validate maxRetryCount
	if maxRetryCount < DefaultUnlimitedRetryCount {
		return errors.New("maximum retry count argument should not be smaller than -1")
	}
	// validate delay
	if delay < constant.ZeroInt {
		return errors.New("maximum delay argument should not be smaller than 0")
	}

	return nil
}
