package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	backoffmock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/backoff/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/httputil"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
)

var jobUUID = "jobUUID"

func testServer(responses ...interface{}) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			var res interface{}
			res, responses = responses[0], responses[1:]
			r, _ := json.Marshal(res)
			if _, ok := res.(httputil.ErrorResponse); ok {
				rw.WriteHeader(500)
			}
			_, _ = rw.Write(r)
		}),
	)
}

func TestClientUpdate_DefaultSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	expectedRes := types.JobResponse{UUID: jobUUID}
	req := &types.UpdateJobRequest{Status: "PENDING"}

	server := testServer(expectedRes)
	client = NewHTTPClient(
		server.Client(),
		NewConfig(server.URL, nil),
	)

	res, err := client.UpdateJob(ctx, jobUUID, req)
	assert.NoError(t, err)
	assert.Equal(t, &expectedRes, res)
}

func TestClientUpdate_DoesNotRetryOnSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	bckoff := &backoffmock.MockIntervalBackoff{}
	expectedRes := types.JobResponse{UUID: jobUUID}
	req := &types.UpdateJobRequest{Status: "PENDING"}

	server := testServer(expectedRes)
	client = NewHTTPClient(
		server.Client(),
		NewConfig(server.URL, bckoff),
	)

	res, err := client.UpdateJob(ctx, jobUUID, req)
	assert.NoError(t, err)
	assert.Equal(t, &expectedRes, res)
	assert.False(t, bckoff.HasRetried())
}

func TestClientUpdate_RetryOnInvalidStateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	bckoff := &backoffmock.MockIntervalBackoff{}
	expectedRes := types.JobResponse{UUID: jobUUID}
	req := &types.UpdateJobRequest{Status: "PENDING"}

	server := testServer(httputil.ErrorResponse{
		Code:    errors.InvalidState,
		Message: "err",
	}, expectedRes)
	client = NewHTTPClient(
		server.Client(),
		NewConfig(server.URL, bckoff),
	)

	res, err := client.UpdateJob(ctx, jobUUID, req)
	assert.NoError(t, err)
	assert.Equal(t, &expectedRes, res)
	assert.True(t, bckoff.HasRetried())
}

func TestClientUpdate_NotRetryOnNotInvalidStateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	bckoff := &backoffmock.MockIntervalBackoff{}

	req := &types.UpdateJobRequest{Status: "PENDING"}

	server := testServer(httputil.ErrorResponse{
		Code:    errors.InvalidParameter,
		Message: "err",
	})

	client = NewHTTPClient(
		server.Client(),
		NewConfig(server.URL, bckoff),
	)

	_, err := client.UpdateJob(ctx, jobUUID, req)
	assert.Error(t, err)
	assert.True(t, errors.IsInvalidParameterError(err))
	assert.False(t, bckoff.HasRetried())
}
