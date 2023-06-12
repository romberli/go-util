package http

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/constant"
)

const (
	testHTTPClientHostInfo   = "http://localhost:6090"
	testHTTPClientGetURL     = testHTTPClientHostInfo + "/status"
	testHTTPClientPostURL    = testHTTPClientHostInfo + "/api/v1/health/ping"
	testHTTPClientDoURL      = testHTTPClientHostInfo + "/api/v1/metadata/mysql-server/is-master/host-info"
	testHTTPClientReqBodyStr = `{"token": "f3171bd9-beec-11ec-acc0-000c291d6734", "host_ip": "192.168.137.11", "port_num": 3306}`

	testHTTPClientSendRequestWithBasicAuthURL = "http://192.168.137.12:8080/api/v2/ob/tenants"
	testHTTPClientBasicAuthUser               = "admin"
	testHTTPClientBasicAuthPass               = "aaAA11.."

	testHTTPClientOCPCreateTenantURL     = "http://192.168.137.12:8080/api/v2/ob/clusters/1/tenants/createTenant"
	testHTTPClientOCPCreateTenantBodyStr = `
		{
			"name": "test_tenant_01",
			"mode": "MYSQL",
			"charset": "utf8mb4",
			"rootPassword": "aaAA11..",
			"saveToCredential": true,
			"zones": [
				{
					"name": "zone1",
					"replicaType": "FULL",
					"resourcePool": {
						"unitSpecName": "OB_X1_V01",
						"unitCount": 1
					}
				},
				{
					"name": "zone2",
					"replicaType": "FULL",
					"resourcePool": {
						"unitSpecName": "OB_X1_V01",
						"unitCount": 1
					}
				},
				{
					"name": "zone3",
					"replicaType": "FULL",
					"resourcePool": {
						"unitSpecName": "OB_X1_V01",
						"unitCount": 1
					}
				}
			]
		}
	`
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

func TestClient_SendRequestWithBasicAuth(t *testing.T) {
	asst := assert.New(t)

	respBody, err := testClient.SendRequestWithBasicAuth(
		http.MethodGet,
		testHTTPClientSendRequestWithBasicAuthURL,
		// []byte(testHTTPClientSendRequestWithBasicAuthBodyStr),
		nil,
		testHTTPClientBasicAuthUser,
		testHTTPClientBasicAuthPass,
	)
	asst.Nil(err, "test SendRequestWithBasicAuth() failed")
	t.Log(string(respBody))
}

func TestClient_GetWithBasicAuth(t *testing.T) {
	asst := assert.New(t)

	respBody, err := testClient.GetWithBasicAuth(
		testHTTPClientSendRequestWithBasicAuthURL,
		nil,
		testHTTPClientBasicAuthUser,
		testHTTPClientBasicAuthPass,
	)
	asst.Nil(err, "test GetWithBasicAuth() failed")
	t.Log(string(respBody))
}

func TestClient_PostWithBasicAuth(t *testing.T) {
	asst := assert.New(t)

	respBody, err := testClient.PostWithBasicAuth(
		testHTTPClientOCPCreateTenantURL,
		[]byte(testHTTPClientOCPCreateTenantBodyStr),
		testHTTPClientBasicAuthUser,
		testHTTPClientBasicAuthPass,
	)
	asst.Nil(err, "test PostWithBasicAuth() failed")
	t.Log(string(respBody))
}
