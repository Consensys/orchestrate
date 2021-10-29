package rpc

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/consensys/orchestrate/pkg/encoding/json"
	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/ethclient/utils"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	log "github.com/sirupsen/logrus"
)

const successStatusClass = 200
const serverErrorStatusClass = 500

type StoreRawResponse struct {
	Key string `json:"key" validate:"required"`
}

func (ec *Client) SendQuorumRawPrivateTransaction(ctx context.Context, endpoint, signedTxHash string, privateFor,
	mandatoryFor []string, privacyFlag int) (ethcommon.Hash, error) {
	privateForParam := map[string]interface{}{
		"privateFor":  privateFor,
		"privacyFlag": privacyFlag,
	}

	if mandatoryFor != nil {
		privateForParam["mandatoryFor"] = mandatoryFor
	}

	var hash string
	err := ec.Call(ctx, endpoint, utils.ProcessResult(&hash), "eth_sendRawPrivateTransaction",
		signedTxHash, privateForParam)
	if err != nil {
		return ethcommon.Hash{}, errors.FromError(err).ExtendComponent(component)
	}
	return ethcommon.HexToHash(hash), nil
}

func (ec *Client) StoreRaw(ctx context.Context, endpoint string, data []byte, privateFrom string) (string, error) {
	request := map[string]string{
		"payload": base64.StdEncoding.EncodeToString(data),
	}
	if privateFrom != "" {
		request["from"] = privateFrom
	}

	storeRawResponse := &StoreRawResponse{}
	err := ec.postRequest(ctx, endpoint, "storeraw", request, storeRawResponse)
	if err != nil {
		return "", errors.FromError(err).SetMessage("failed to send a request to Tessera enclave: %s", err)
	}

	enclaveKey, err := base64.StdEncoding.DecodeString(storeRawResponse.Key)
	if err != nil {
		return "", errors.DataCorruptedError("failed to decode base64 encoded string in the 'storeraw' response: %s", err)
	}
	return hexutil.Encode(enclaveKey), nil
}

func (ec *Client) GetStatus(ctx context.Context, endpoint string) (status string, err error) {
	return ec.getRequest(ctx, endpoint, "upcheck")
}

func (ec *Client) postRequest(ctx context.Context, endpoint, path string, request, reply interface{}) error {
	requestURL := fmt.Sprintf("%s/%s", endpoint, path)

	log.Debugf("Sending POST request to %s", requestURL)

	jsonValue, err := json.Marshal(request)
	if err != nil {
		return errors.DataCorruptedError("failed to marshal an HTTP request - %q", err)
	}

	resp, err := ec.sendPostRequest(ctx, requestURL, jsonValue)
	if err != nil {
		return err
	}

	return readResponse(requestURL, resp, reply)
}

func (ec *Client) getRequest(ctx context.Context, endpoint, path string) (string, error) {
	requestURL := fmt.Sprintf("%s/%s", endpoint, path)

	log.Debugf("Sending GET request to %s", requestURL)

	resp, err := ec.sendGetRequest(ctx, requestURL)
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

	err := json.UnmarshalBody(resp.Body, result)

	if err != nil {
		return errors.DataCorruptedError("failed to parse reply from '%s' request: %s", requestURL, err)
	}

	log.Debugf("received the following reply from '%s' endpoint: %q", requestURL, result)

	return nil
}

func (ec *Client) sendPostRequest(ctx context.Context, requestURL string, jsonValue []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return ec.sendRequest(ctx, req)
}

func (ec *Client) sendGetRequest(ctx context.Context, requestURL string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	return ec.sendRequest(ctx, req)
}

func (ec *Client) sendRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	bckoff := backoff.WithContext(ec.Pool().Get().(backoff.BackOff), ctx)
	defer ec.Pool().Put(bckoff)

	return ec.retryHTTPRequest(
		req.URL.String(),
		func() (*http.Response, error) { return ec.HTTPClient().Do(req) },
		bckoff,
	)
}

func (ec *Client) retryHTTPRequest(requestURL string, doRequest func() (*http.Response, error), bckoff backoff.BackOff) (*http.Response, error) {
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
				switch resp.StatusCode {
				// Follow official go-quorum docs https://consensys.github.io/tessera/#operation/encryptAndStoreVersion
				case 404:
					err = errors.InvalidParameterError("'from' key in request body not found")
				default:
					err = errors.HTTPConnectionError("request to '%s' failed. code: %d", requestURL, resp.StatusCode)
				}

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
