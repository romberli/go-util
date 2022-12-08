package http

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testHTTPClientHostInfo   = "http://localhost:6090"
	testHTTPClientGetURL     = testHTTPClientHostInfo + "/status"
	testHTTPClientPostURL    = testHTTPClientHostInfo + "/api/v1/health/ping"
	testHTTPClientDoURL      = testHTTPClientHostInfo + "/api/v1/metadata/mysql-server/is-master/host-info"
	testHTTPClientReqBodyStr = `{"token":"f3171bd9-beec-11ec-acc0-000c291d6734", "host_ip": "192.168.137.11", "port_num": 3306}`
)

var testClient *Client

func init() {
	testClient = testInitClient()
}

func testInitClient() *Client {
	return NewClientWithDefault()
}

func TestClient_All(t *testing.T) {
	TestClient_Close(t)
	TestClient_Do(t)
	TestClient_Get(t)
	TestClient_Post(t)
}

func TestClient_Close(t *testing.T) {
	testClient.Close()
}

func TestClient_Do(t *testing.T) {
	asst := assert.New(t)
	body := []byte(testHTTPClientReqBodyStr)
	req, err := http.NewRequest(http.MethodPost, testHTTPClientDoURL, bytes.NewBuffer(body))
	asst.Nil(err, "test Do() failed")
	resp, err := testClient.Do(req)
	asst.Nil(err, "test Do() failed")
	// read response body
	respBody, err := io.ReadAll(resp.Body)
	defer func() { _ = resp.Body.Close() }()

	asst.Nil(err, "test Do() failed")
	asst.Equal(http.StatusOK, resp.StatusCode, "test Do() failed. statusCode: %s", resp.StatusCode)
	t.Log(string(respBody))
}

func TestClient_Get(t *testing.T) {
	asst := assert.New(t)

	resp, err := testClient.Get(testHTTPClientGetURL)
	asst.Nil(err, "test Get() failed")
	// read response body
	respBody, err := io.ReadAll(resp.Body)
	defer func() { _ = resp.Body.Close() }()

	asst.Nil(err, "test Get() failed")
	asst.Equal(http.StatusOK, resp.StatusCode, "test Get() failed. statusCode: %s", resp.StatusCode)
	t.Log(string(respBody))
}

func TestClient_GetWithRetry(t *testing.T) {
	asst := assert.New(t)

	resp, err := testClient.GetWithRetry(testHTTPClientGetURL, DefaultMaxWaitTime, DefaultMaxRetryCount, DefaultDelay)
	asst.Nil(err, "test GetWithRetry() failed")
	// read response body
	respBody, err := io.ReadAll(resp.Body)
	defer func() { _ = resp.Body.Close() }()

	asst.Nil(err, "test GetWithRetry() failed")
	asst.Equal(http.StatusOK, resp.StatusCode, "test GetWithRetry() failed. statusCode: %s", resp.StatusCode)
	t.Log(string(respBody))
}

func TestClient_Post(t *testing.T) {
	asst := assert.New(t)

	resp, err := testClient.Post(testHTTPClientPostURL, nil)
	asst.Nil(err, "test Post() failed")
	// read response body
	respBody, err := io.ReadAll(resp.Body)
	defer func() { _ = resp.Body.Close() }()

	asst.Nil(err, "test Post() failed")
	asst.Equal(http.StatusOK, resp.StatusCode, "test Post() failed. statusCode: %s", resp.StatusCode)
	t.Log(string(respBody))
}

func TestClient_PostWithRetry(t *testing.T) {
	asst := assert.New(t)

	resp, err := testClient.PostWithRetry(testHTTPClientPostURL, nil, DefaultMaxWaitTime, DefaultMaxRetryCount, DefaultDelay)
	asst.Nil(err, "test GetWithRetry() failed")
	// read response body
	respBody, err := io.ReadAll(resp.Body)
	defer func() { _ = resp.Body.Close() }()

	asst.Nil(err, "test GetWithRetry() failed")
	asst.Equal(http.StatusOK, resp.StatusCode, "test GetWithRetry() failed. statusCode: %d", resp.StatusCode)
	t.Log(string(respBody))
}

func TestClient_PostDAS(t *testing.T) {
	asst := assert.New(t)

	resp, err := testClient.PostDAS(testHTTPClientPostURL, nil, DefaultMaxWaitTime, DefaultMaxRetryCount, DefaultDelay)
	asst.Nil(err, "test PostDAS() failed")
	t.Log(string(resp))
}
