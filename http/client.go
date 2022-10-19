package http

import (
	"bytes"
	errs "errors"
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

	defaultDialTimeout         = 30 * time.Second
	defaultKeepAlive           = 30 * time.Second
	defaultTLSHandshakeTimeout = 10 * time.Second
	defaultContentType         = "application/json"
)

var (
	DefaultTransport = &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		DialContext:         (&net.Dialer{Timeout: defaultDialTimeout, KeepAlive: defaultKeepAlive}).DialContext,
		TLSHandshakeTimeout: defaultTLSHandshakeTimeout,
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
		Timeout:   defaultDialTimeout,
	}

	return newClient(client)
}

func newClient(client *http.Client) *Client {
	return &Client{
		client: client,
	}
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

func (c *Client) Get(url string) (*http.Response, error) {
	return c.client.Get(PrepareURL(url))
}

func (c *Client) Post(url string, body []byte) (*http.Response, error) {
	return c.postJSON(url, body)
}

func (c *Client) postJSON(url string, body []byte) (*http.Response, error) {
	return c.client.Post(PrepareURL(url), defaultContentType, bytes.NewBuffer(body))
}

func (c *Client) PostDAS(url string, body []byte) ([]byte, error) {
	resp, err := c.Post(url, body)
	if err != nil {
		return nil, err
	}
	// read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Errorf("http Client.PostAndParse(): http response body failed. error:\n%+v", err)
		}
	}()
	// check http status code
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("got wrong http response status code. status code: %d, body: %s", resp.StatusCode, respBody)
	}

	code, err := jsonparser.GetInt(respBody, defaultResponseCodeJSON)
	if err != nil && !errs.As(err, &jsonparser.KeyPathNotFoundError) {
		return nil, errors.Errorf("got error when getting code field from response body. error:\n%+v", err)
	}
	if code != constant.ZeroInt {
		return nil, errors.Errorf("code field in response body is not 0. code: %d, body: %s", code, respBody)
	}

	return respBody, nil
}
