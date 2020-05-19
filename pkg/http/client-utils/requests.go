package clientutils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

func GetRequest(ctx context.Context, client *http.Client, reqURL string) (*http.Response, error) {
	return request(ctx, client, reqURL, http.MethodGet, nil)
}

func DeleteRequest(ctx context.Context, client *http.Client, reqURL string) (*http.Response, error) {
	return request(ctx, client, reqURL, http.MethodDelete, nil)
}

func PostRequest(ctx context.Context, client *http.Client, reqURL string, postRequest interface{}) (*http.Response, error) {
	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(postRequest)

	return request(ctx, client, reqURL, http.MethodPost, body)
}

func PatchRequest(ctx context.Context, client *http.Client, reqURL string, patchRequest interface{}) (*http.Response, error) {
	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(patchRequest)

	return request(ctx, client, reqURL, http.MethodPatch, body)
}

func PutRequest(ctx context.Context, client *http.Client, reqURL string, putRequest interface{}) (*http.Response, error) {
	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(putRequest)

	return request(ctx, client, reqURL, http.MethodPut, body)
}

func CloseResponse(response *http.Response) {
	if deferErr := response.Body.Close(); deferErr != nil {
		log.WithError(deferErr).Errorf("could not close response body")
	}
}

func request(ctx context.Context, client *http.Client, reqURL, method string, body io.Reader) (*http.Response, error) {
	req, _ := http.NewRequestWithContext(ctx, method, reqURL, body)
	r, err := client.Do(req)
	if err != nil {
		return nil, errors.FromError(err)
	}

	return r, nil
}
