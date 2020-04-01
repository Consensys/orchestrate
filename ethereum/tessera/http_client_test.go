// +build unit

package tessera

import (
	"net/http"
	"testing"

	"github.com/cenkalti/backoff"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

type TestStruct struct {
	Field string `json:"field"`
}

var expectedReply = TestStruct{
	Field: "bar",
}

var backOff = backoff.WithMaxRetries(backoff.NewConstantBackOff(0), 3)
var enclaveEndpoint = CreateEnclaveHTTPEndpointWithConfig("http://test-endpoint", backOff)
var response = TestStruct{}

func TestSuccessfulPostRequest(t *testing.T) {
	defer gock.Off()

	gock.New("http://test-endpoint").
		Post("bar").
		Reply(200).
		JSON(map[string]string{"field": "bar"})

	err := enclaveEndpoint.PostRequest("bar", map[string]interface{}{}, &response)

	assert.NoError(t, err, "HTTP request should not fail")
	assert.Equal(t, expectedReply, response, "Unexpected response")
}

func TestParseStringResponse(t *testing.T) {
	defer gock.Off()

	gock.New("http://test-endpoint").
		Get("bar").
		Reply(200).
		JSON("reply")

	reply, err := enclaveEndpoint.GetRequest("bar")

	assert.NoError(t, err, "HTTP request should not fail")
	assert.Equal(t, "reply", reply, "Unexpected response")
}

func TestGetRequestWithHTTPErrorCode(t *testing.T) {
	defer gock.Off()

	gock.New("http://test-endpoint").
		Get("bar").
		Reply(400)

	reply, err := enclaveEndpoint.GetRequest("bar")

	assert.EqualError(t, err, "request to 'http://test-endpoint/bar' failed - 400")
	assert.Equal(t, "", reply, "Unexpected response")
}

func TestRequestWithInvalidBody(t *testing.T) {
	defer gock.Off()

	gock.New("http://test-endpoint").
		Post("bar").
		Reply(200).
		BodyString("{")

	err := enclaveEndpoint.PostRequest("bar", map[string]interface{}{}, &response)

	assert.EqualError(t, err, "failed to parse reply from 'http://test-endpoint/bar' request: unexpected end of JSON input")
}

func TestRequestWithHTTPErrorCode(t *testing.T) {
	defer gock.Off()

	gock.New("http://test-endpoint").
		Post("bar").
		Reply(400)

	err := enclaveEndpoint.PostRequest("bar", map[string]interface{}{}, &response)

	assert.EqualError(t, err, "request to 'http://test-endpoint/bar' failed - 400")
}

func TestRequestWithRetries(t *testing.T) {
	defer gock.Off()
	requestsNum := 0
	gock.Observe(func(request *http.Request, mock gock.Mock) {
		requestsNum++
	})

	gock.New("http://test-endpoint").
		Post("bar").
		Reply(500)

	err := enclaveEndpoint.PostRequest("bar", map[string]interface{}{}, &response)

	assert.EqualError(t, err, "failed to send a request to 'http://test-endpoint/bar' - Post http://test-endpoint/bar: gock: cannot match any request")
	assert.Equal(t, 4, requestsNum, "should retry")
}
