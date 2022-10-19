package http

import (
	"bytes"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	defaultHTTPScheme         = "http://"
	defaultHTTPSScheme        = "https://"
	StatusOk                  = http.StatusOK
	StatusInternalServerError = http.StatusInternalServerError

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
	return c.client.Get(c.getURL(url))
}

func (c *Client) Post(url string, body []byte) (*http.Response, error) {
	return c.postJSON(url, body)
}

func (c *Client) postJSON(url string, body []byte) (*http.Response, error) {
	return c.client.Post(c.getURL(url), defaultContentType, bytes.NewBuffer(body))
}

func (c *Client) getURL(url string) string {
	if strings.HasPrefix(url, defaultHTTPScheme) || strings.HasPrefix(url, defaultHTTPSScheme) {
		return url
	}

	return defaultHTTPScheme + url
}
