// +build unit

package tessera

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

var testConfig = &utils.Config{
	Retry: &utils.RetryConfig{
		InitialInterval:     time.Millisecond,
		RandomizationFactor: 0.5,
		Multiplier:          1.5,
		MaxInterval:         time.Millisecond,
		MaxElapsedTime:      time.Millisecond,
	},
}

type mockRoundTripper struct{}
type testCtxKey string

func (rt mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if err, ok := req.Context().Value(testCtxKey("resp.error")).(error); ok {
		return nil, err
	}

	resp := &http.Response{}
	if statusCode, ok := req.Context().Value(testCtxKey("resp.statusCode")).(int); ok {
		resp.StatusCode = statusCode
		resp.Status = http.StatusText(statusCode)
	}

	if body, ok := req.Context().Value(testCtxKey("resp.body")).(io.ReadCloser); ok {
		resp.Body = body
	}

	return resp, nil
}

func newClient() *HTTPClient {
	newBackOff := func() backoff.BackOff { return utils.NewBackOff(testConfig) }
	return NewTesseraClient(newBackOff, &http.Client{
		Transport: mockRoundTripper{},
	})
}

func newContext(err error, statusCode int, body io.ReadCloser) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, testCtxKey("resp.error"), err)
	ctx = context.WithValue(ctx, testCtxKey("resp.statusCode"), statusCode)
	ctx = context.WithValue(ctx, testCtxKey("resp.body"), body)
	return ctx
}

func TestStoreRawTransaction(t *testing.T) {
	testSet := []struct {
		name                   string
		httpBodyResponse       interface{}
		httpStatusCodeResponse int
		httpErrorResponse      error
		expectedEnclaveKey     string
		expectedError          error
	}{
		{
			"success storeraw",
			StoreRawResponse{Key: "test"},
			200,
			nil,
			"0xb5eb2d",
			nil,
		},
		{
			"fail storeraw",
			StoreRawResponse{Key: "test"},
			400,
			fmt.Errorf("test"),
			"",
			errors.HTTPConnectionError("failed to send a request to Tessera enclave: 08200@: failed to send a request to 'test/storeraw' - Post \"test/storeraw\": test"),
		},
	}
	ec := newClient()

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			b, _ := json.Marshal(test.httpBodyResponse)

			ctx := newContext(test.httpErrorResponse, test.httpStatusCodeResponse, ioutil.NopCloser(bytes.NewReader(b)))

			enclaveKey, err := ec.StoreRaw(ctx, "test", []byte{}, "testPrivateFrom")
			assert.Equal(t, err, test.expectedError)
			assert.Equal(t, enclaveKey, test.expectedEnclaveKey)
		})
	}
}

func TestGetStatus(t *testing.T) {
	testSet := []struct {
		name                   string
		httpBodyResponse       interface{}
		httpStatusCodeResponse int
		httpErrorResponse      error
		expectedStatus         string
		expectedError          error
	}{
		{
			"success get status",
			"testStatus",
			200,
			nil,
			"\"testStatus\"",
			nil,
		},
	}
	ec := newClient()

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			b, _ := json.Marshal(test.httpBodyResponse)

			ctx := newContext(test.httpErrorResponse, test.httpStatusCodeResponse, ioutil.NopCloser(bytes.NewReader(b)))

			status, err := ec.GetStatus(ctx, "testEndpoint")
			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedStatus, status)
		})
	}
}
