package http

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/buger/jsonparser"
	"github.com/pingcap/errors"
	"github.com/romberli/log"

	"github.com/romberli/go-util/constant"
)

const (
	MethodGet  = http.MethodGet
	MethodPost = http.MethodPost

	DefaultContentTypeKey   = "Content-Type"
	DefaultContentTypeValue = "application/json"

	StatusOK                  = http.StatusOK
	StatusInternalServerError = http.StatusInternalServerError

	defaultResponseCodeJSON = "code"

	defaultClientTimeout         = 60 * time.Second
	defaultDialTimeout           = 30 * time.Second
	defaultKeepAlive             = 30 * time.Second
	defaultTLSHandshakeTimeout   = 10 * time.Second
	defaultMaxIdleConns          = 100
	defaultIdleConnTimeout       = 90 * time.Second
	defaultExpectContinueTimeout = 1 * time.Second
	defaultMaxIdleConnsPerHost   = 10

	DefaultMaxWaitTime   = 60 // seconds
	DefaultMaxRetryCount = 3
	DefaultDelayTime     = 10 // milliseconds

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
	client        *http.Client
	maxWaitTime   int
	maxRetryCount int
	delayTime     int
}

func NewClient(client *http.Client, maxWaitTime, maxRetryCount, delayTime int) (*Client, error) {
	return newClient(client, maxWaitTime, maxRetryCount, delayTime)
}

func NewClientWithDefault() (*Client, error) {
	client := &http.Client{
		Transport: DefaultTransport,
		Timeout:   defaultClientTimeout,
	}

	return newClient(client, constant.ZeroInt, constant.ZeroInt, constant.ZeroInt)
}

func NewClientWithRetry() (*Client, error) {
	client := &http.Client{
		Transport: DefaultTransport,
		Timeout:   defaultClientTimeout,
	}

	return newClient(client, DefaultMaxWaitTime, DefaultMaxRetryCount, DefaultDelayTime)
}

func newClient(client *http.Client, maxWaitTime, maxRetryCount, delayTime int) (*Client, error) {
	c := &Client{
		client:        client,
		maxWaitTime:   maxWaitTime,
		maxRetryCount: maxRetryCount,
		delayTime:     delayTime,
	}

	err := c.validate()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) validate() error {
	// validate maxWaitTime
	if c.maxWaitTime < DefaultUnlimitedWaitTime {
		return errors.New("maximum wait time argument should not be smaller than -1")
	}
	// validate maxRetryCount
	if c.maxRetryCount < DefaultUnlimitedRetryCount {
		return errors.New("maximum retry count argument should not be smaller than -1")
	}
	// validate delay
	if c.delayTime < constant.ZeroInt {
		return errors.New("delay time argument should not be smaller than 0")
	}

	return nil
}

func (c *Client) GetClient() *http.Client {
	return c.client
}

func (c *Client) Close() {
	c.client.CloseIdleConnections()
}

func (c *Client) SetMaxIdleConns(maxIdleConns int) {
	c.client.Transport.(*http.Transport).MaxIdleConns = maxIdleConns
}

func (c *Client) SetMaxIdleConnsPerHost(maxIdleConnsPerHost int) {
	c.client.Transport.(*http.Transport).MaxIdleConnsPerHost = maxIdleConnsPerHost
}

func (c *Client) SetRetryOption(maxWaitTime, maxRetryCount, delay int) {
	c.maxWaitTime = maxWaitTime
	c.maxRetryCount = maxRetryCount
	c.delayTime = delay
}

// PrepareURL prepares the url
func (c *Client) PrepareURL(scheme, addr, path string, params map[string]string) string {
	query := url.Values{}
	for key, value := range params {
		query.Add(key, value)
	}

	u := &url.URL{
		Scheme:   scheme,
		Host:     addr,
		Path:     path,
		RawQuery: query.Encode(),
	}

	return u.String()
}

func (c *Client) Get(url string) (*http.Response, error) {
	maxWait := c.maxWaitTime
	if maxWait < constant.ZeroInt {
		maxWait = int(constant.Century.Seconds())
	}
	timeoutChan := time.After(time.Duration(maxWait) * time.Second)

	var i int

	for {
		resp, err := c.client.Get(PrepareURL(url))
		if err != nil {
			// check retry count
			if c.maxRetryCount >= constant.ZeroInt && i >= c.maxRetryCount {
				return resp, errors.Trace(err)
			}
			// check wait time
			select {
			case <-timeoutChan:
				return resp, errors.Trace(err)
			default:
				time.Sleep(time.Duration(c.delayTime) * time.Millisecond)
			}

			i++
			continue
		}

		return resp, nil
	}
}

func (c *Client) Post(url string, body []byte) (*http.Response, error) {
	maxWait := c.maxWaitTime
	if maxWait < constant.ZeroInt {
		maxWait = int(constant.Century.Seconds())
	}
	timeoutChan := time.After(time.Duration(maxWait) * time.Second)

	var i int

	for {
		resp, err := c.client.Post(PrepareURL(url), DefaultContentTypeValue, bytes.NewBuffer(body))
		if err != nil {
			// check retry count
			if c.maxRetryCount >= constant.ZeroInt && i >= c.maxRetryCount {
				return resp, errors.Trace(err)
			}
			// check wait timeout
			select {
			case <-timeoutChan:
				return resp, errors.Trace(err)
			default:
				time.Sleep(time.Duration(c.delayTime) * time.Millisecond)
			}

			i++
			continue
		}

		return resp, nil
	}
}

func (c *Client) PostDAS(url string, body []byte) ([]byte, error) {
	resp, err := c.Post(url, body)
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

func (c *Client) SendRequestWithBasicAuth(method, url string, body []byte, user, pass string) ([]byte, error) {
	maxWait := c.maxWaitTime
	if maxWait < constant.ZeroInt {
		maxWait = int(constant.Century.Seconds())
	}
	timeoutChan := time.After(time.Duration(maxWait) * time.Second)

	var i int

	for {
		resp, err := c.sendRequestWithBasicAuth(method, url, body, user, pass)
		if err != nil {
			// check retry count
			if c.maxRetryCount >= constant.ZeroInt && i >= c.maxRetryCount {
				return resp, errors.Trace(err)
			}
			// check wait timeout
			select {
			case <-timeoutChan:
				return resp, errors.Trace(err)
			default:
				time.Sleep(time.Duration(c.delayTime) * time.Millisecond)
			}

			i++
			continue
		}

		return resp, nil
	}
}

func (c *Client) SendRequestWithHeaderAndBody(method, url string, header map[string]string, body []byte) ([]byte, error) {
	maxWait := c.maxWaitTime
	if maxWait < constant.ZeroInt {
		maxWait = int(constant.Century.Seconds())
	}
	timeoutChan := time.After(time.Duration(maxWait) * time.Second)

	var i int

	for {
		resp, err := c.sendRequestWithHeaderAndBody(method, url, header, body)
		if err != nil {
			// check retry count
			if c.maxRetryCount >= constant.ZeroInt && i >= c.maxRetryCount {
				return resp, errors.Trace(err)
			}
			// check wait timeout
			select {
			case <-timeoutChan:
				return resp, errors.Trace(err)
			default:
				time.Sleep(time.Duration(c.delayTime) * time.Millisecond)
			}

			i++
			continue
		}

		return resp, nil
	}
}

func (c *Client) GetWithBasicAuth(url string, body []byte, user, pass string) ([]byte, error) {
	return c.SendRequestWithBasicAuth(http.MethodGet, url, body, user, pass)
}

func (c *Client) PostWithBasicAuth(url string, body []byte, user, pass string) ([]byte, error) {
	return c.SendRequestWithBasicAuth(http.MethodPost, url, body, user, pass)
}

func (c *Client) sendRequestWithBasicAuth(method, url string, body []byte, user, pass string) ([]byte, error) {
	req, err := http.NewRequest(method, PrepareURL(url), bytes.NewReader(body))
	if err != nil {
		return nil, errors.Trace(err)
	}

	req.Header.Set(DefaultContentTypeKey, DefaultContentTypeValue)
	req.SetBasicAuth(user, pass)

	resp, err := c.GetClient().Do(req)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("metadata Repository.SendOCPRequest(): request failed. statusCode: %d, respBody: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (c *Client) sendRequestWithHeaderAndBody(method, url string, header map[string]string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, PrepareURL(url), bytes.NewReader(body))
	if err != nil {
		return nil, errors.Trace(err)
	}

	req.Header.Set(DefaultContentTypeKey, DefaultContentTypeValue)
	for k, v := range header {
		req.Header.Set(k, v)
	}

	resp, err := c.GetClient().Do(req)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("metadata Repository.SendOCPRequest(): request failed. statusCode: %d, header: %v, respBody: %s", resp.StatusCode, header, string(respBody))
	}

	return respBody, nil
}
