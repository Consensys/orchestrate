package tessera

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"

	log "github.com/sirupsen/logrus"
)

const successStatusClass = 200
const serverErrorStatusClass = 500

type HTTPClient struct {
	client *http.Client
	pool   *sync.Pool
}

type StoreRawResponse struct {
	Key string `json:"key"`
}

func NewTesseraClient(newBackOff func() backoff.BackOff, client *http.Client) *HTTPClient {
	return &HTTPClient{
		client: client,
		pool: &sync.Pool{
			New: func() interface{} { return newBackOff() },
		},
	}
}

func (tc *HTTPClient) StoreRaw(ctx context.Context, endpoint string, data []byte, privateFrom string) (string, error) {
	request := map[string]string{
		"payload": base64.StdEncoding.EncodeToString(data),
		"from":    privateFrom,
	}

	storeRawResponse := StoreRawResponse{}
	err := tc.postRequest(ctx, endpoint, "storeraw", request, &storeRawResponse)
	if err != nil {
		return "", errors.HTTPConnectionError("failed to send a request to Tessera enclave: %s", err)
	}
	enclaveKey, err := base64.StdEncoding.DecodeString(storeRawResponse.Key)
	if err != nil {
		return "", errors.HTTPConnectionError("failed to decode base64 encoded string in the 'storeraw' response: %s", err)
	}

	return hexutil.Encode(enclaveKey), nil
}

func (tc *HTTPClient) GetStatus(ctx context.Context, endpoint string) (status string, err error) {
	return tc.getRequest(ctx, endpoint, "upcheck")
}

func (tc *HTTPClient) postRequest(ctx context.Context, endpoint, path string, request, reply interface{}) error {
	requestURL := fmt.Sprintf("%s/%s", endpoint, path)

	log.Debugf("Sending POST request to %s with body %q", requestURL, request)

	jsonValue, err := json.Marshal(request)
	if err != nil {
		return errors.DataCorruptedError("failed to marshal an HTTP request - %q", err)
	}

	resp, err := tc.sendPostRequest(ctx, requestURL, jsonValue)
	if err != nil {
		return err
	}

	return readResponse(requestURL, resp, reply)
}

func (tc *HTTPClient) getRequest(ctx context.Context, endpoint, path string) (string, error) {
	requestURL := fmt.Sprintf("%s/%s", endpoint, path)

	log.Debugf("Sending GET request to %s", requestURL)

	resp, err := tc.sendGetRequest(ctx, requestURL)
	if err != nil {
		return "", errors.HTTPConnectionError("%v", err)
	}

	return readResponseString(requestURL, resp)
}

func readResponseString(requestURL string, resp *http.Response) (string, error) {
	log.Debugf("request to '%s' resulted with %d status code", requestURL, resp.StatusCode)

	reply, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.DataCorruptedError("failed to read body for a request from '%s' endpoint: %s", requestURL, err)
	}

	log.Debugf("received the following reply from '%s' endpoint: %q", requestURL, string(reply))

	return string(reply), nil
}

func readResponse(requestURL string, resp *http.Response, result interface{}) error {
	log.Debugf("request to '%s' resulted with %d status code", requestURL, resp.StatusCode)

	reply, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.DataCorruptedError("failed to read body for a request from '%s' endpoint: %s", requestURL, err)
	}

	log.Debugf("received the following reply from '%s' endpoint: %q", requestURL, string(reply))

	err = json.Unmarshal(reply, &result)
	if err != nil {
		return errors.DataCorruptedError("failed to parse reply from '%s' request: %s", requestURL, err)
	}

	log.Debugf("received the following reply from '%s' endpoint: %q", requestURL, result)

	return nil
}

func (tc *HTTPClient) sendPostRequest(ctx context.Context, requestURL string, jsonValue []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return tc.sendRequest(ctx, req)
}

func (tc *HTTPClient) sendGetRequest(ctx context.Context, requestURL string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	return tc.sendRequest(ctx, req)
}

func (tc *HTTPClient) sendRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	bckoff := backoff.WithContext(tc.pool.Get().(backoff.BackOff), ctx)
	defer tc.pool.Put(bckoff)

	return tc.retryHTTPRequest(
		req.URL.String(),
		func() (*http.Response, error) { return tc.client.Do(req) },
		bckoff,
	)
}

func (tc *HTTPClient) retryHTTPRequest(requestURL string, doRequest func() (*http.Response, error), bckoff backoff.BackOff) (*http.Response, error) {
	var resp *http.Response
	err := backoff.RetryNotify(
		func() error {
			var err error
			resp, err = doRequest()
			if err != nil {
				return errors.HTTPConnectionError("failed to send a request to '%s' - %s", requestURL, err)
			}

			statusClass := getStatusClass(resp)
			if statusClass != successStatusClass {
				err = errors.HTTPConnectionError("request to '%s' failed - %d", requestURL, resp.StatusCode)
				if statusClass != serverErrorStatusClass {
					err = backoff.Permanent(err)
				}
				return err
			}

			return nil
		},
		bckoff,
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
