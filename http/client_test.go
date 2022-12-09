package http

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/log"
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
	c, err := NewClientWithRetry()
	if err != nil {
		log.Errorf("testInitClient() failed. err:\n%+v", err)
		os.Exit(constant.DefaultAbnormalExitCode)
	}

	return c
}

func TestClient_All(t *testing.T) {
	TestClient_Get(t)
	TestClient_Post(t)
	TestClient_PostDAS(t)
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

func TestClient_PostDAS(t *testing.T) {
	asst := assert.New(t)

	resp, err := testClient.PostDAS(testHTTPClientPostURL, nil)
	asst.Nil(err, "test PostDAS() failed")
	t.Log(string(resp))
}
