package tessera

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cenkalti/backoff"

	log "github.com/sirupsen/logrus"
)

const successStatusClass = 200
const serverErrorStatusClass = 500

type EnclaveHTTPEndpoint struct {
	endpoint       string
	requestBackOff backoff.BackOff
}

func CreateEnclaveHTTPEndpoint(endpoint string) *EnclaveHTTPEndpoint {
	return CreateEnclaveHTTPEndpointWithConfig(endpoint, backoff.NewExponentialBackOff())
}

func CreateEnclaveHTTPEndpointWithConfig(endpoint string, requestBackoff backoff.BackOff) *EnclaveHTTPEndpoint {
	return &EnclaveHTTPEndpoint{
		endpoint:       endpoint,
		requestBackOff: requestBackoff,
	}
}

func (c *EnclaveHTTPEndpoint) PostRequest(path string, request, reply interface{}) error {
	requestURL := fmt.Sprintf("%s/%s", c.endpoint, path)

	log.Debugf("Sending POST request to %s with body %q", requestURL, request)

	jsonValue, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal an HTTP request - %q", err)
	}

	resp, err := c.sendPostRequest(requestURL, jsonValue)
	if err != nil {
		return err
	}

	return readResponse(requestURL, resp, reply)
}

func (c *EnclaveHTTPEndpoint) GetRequest(path string) (string, error) {
	requestURL := fmt.Sprintf("%s/%s", c.endpoint, path)

	log.Debugf("Sending GET request to %s", requestURL)

	resp, err := c.sendGetRequest(requestURL)
	if err != nil {
		return "", err
	}

	return readResponseString(requestURL, resp)
}

func readResponseString(requestURL string, resp *http.Response) (string, error) {
	log.Debugf("request to '%s' resulted with %d status code", requestURL, resp.StatusCode)

	reply, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body for a request from '%s' endpoint: %s", requestURL, err)
	}

	log.Debugf("received the following reply from '%s' endpoint: %q", requestURL, string(reply))

	return string(reply), nil
}

func readResponse(requestURL string, resp *http.Response, result interface{}) error {
	log.Debugf("request to '%s' resulted with %d status code", requestURL, resp.StatusCode)

	reply, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body for a request from '%s' endpoint: %s", requestURL, err)
	}

	log.Debugf("received the following reply from '%s' endpoint: %q", requestURL, string(reply))

	err = json.Unmarshal(reply, &result)
	if err != nil {
		return fmt.Errorf("failed to parse reply from '%s' request: %s", requestURL, err)
	}

	log.Debugf("received the following reply from '%s' endpoint: %q", requestURL, result)

	return nil
}

func (c *EnclaveHTTPEndpoint) sendPostRequest(requestURL string, jsonValue []byte) (*http.Response, error) {
	return c.retryHTTPRequest(requestURL, func() (*http.Response, error) {
		return http.Post(requestURL, "application/json", bytes.NewBuffer(jsonValue))
	})
}

func (c *EnclaveHTTPEndpoint) sendGetRequest(requestURL string) (*http.Response, error) {
	return c.retryHTTPRequest(requestURL, func() (*http.Response, error) {
		return http.Get(requestURL)
	})
}

func (c *EnclaveHTTPEndpoint) retryHTTPRequest(requestURL string, doRequest func() (*http.Response, error)) (*http.Response, error) {

	var resp *http.Response
	err := backoff.RetryNotify(
		func() error {
			var err error
			resp, err = doRequest()
			if err != nil {
				return fmt.Errorf("failed to send a request to '%s' - %s", requestURL, err)
			}

			statusClass := getStatusClass(resp)
			if statusClass != successStatusClass {
				err = fmt.Errorf("request to '%s' failed - %d", requestURL, resp.StatusCode)
				if statusClass != serverErrorStatusClass {
					err = backoff.Permanent(err)
				}
				return err
			}

			return nil
		},
		c.requestBackOff,
		func(err error, duration time.Duration) {
			log.
				WithError(err).
				WithFields(log.Fields{
					"requestURL": requestURL,
				}).Warnf("tessera-http-client: error sending requests to Tessera, retrying in %v...", duration)
		},
	)

	return resp, err
}

func getStatusClass(resp *http.Response) int {
	return resp.StatusCode - resp.StatusCode%100
}
