package testutils

import (
	"context"
	"io"
	"net/http"
	"time"

	pkgUtils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

var TestConfig = &pkgUtils.Config{
	Retry: &pkgUtils.RetryConfig{
		InitialInterval:     time.Millisecond,
		RandomizationFactor: 0.5,
		Multiplier:          1.5,
		MaxInterval:         time.Millisecond,
		MaxElapsedTime:      time.Millisecond,
	},
}

type MockRoundTripper struct{}

var skipPreCallRoundTrip bool

func (rt MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	if preCtx, ok := ctx.Value(TestCtxKey("pre_call")).(context.Context); ok && !skipPreCallRoundTrip {
		skipPreCallRoundTrip = true
		ctx = preCtx
	}

	if err, ok := ctx.Value(TestCtxKey("resp.error")).(error); ok {
		return nil, err
	}

	resp := &http.Response{}
	if statusCode, ok := ctx.Value(TestCtxKey("resp.statusCode")).(int); ok {
		resp.StatusCode = statusCode
		resp.Status = http.StatusText(statusCode)
	}

	if body, ok := ctx.Value(TestCtxKey("resp.body")).(io.ReadCloser); ok {
		resp.Body = body
	}

	return resp, nil
}
