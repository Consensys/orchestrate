package client

import (
	"context"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	backoff2 "github.com/consensys/orchestrate/pkg/backoff"
	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
)

type HTTPClient struct {
	client *http.Client
	config *Config
}

var _ OrchestrateClient = &HTTPClient{}

func NewHTTPClient(h *http.Client, c *Config) *HTTPClient {
	return &HTTPClient{
		client: h,
		config: c,
	}
}

func callWithBackOff(ctx context.Context, backOff backoff2.BackOff, requestCall func() error) error {
	return backoff.RetryNotify(
		func() error {
			err := requestCall()
			// If not errors, it does not retry
			if err == nil {
				return nil
			}

			if err == context.Canceled || err == context.DeadlineExceeded {
				return backoff.Permanent(err)
			}

			if ctx.Err() != nil {
				return backoff.Permanent(ctx.Err())
			}

			// Retry on following errors
			if errors.IsInvalidStateError(err) || errors.IsServiceConnectionError(err) || errors.IsDependencyFailureError(err) {
				return err
			}

			// Otherwise, stop retrying
			return backoff.Permanent(err)
		}, backoff.WithContext(backOff.NewBackOff(), ctx),
		func(e error, duration time.Duration) {
			log.FromContext(ctx).
				WithError(e).
				Warnf("http call has failed, retrying in %v...", duration)
		},
	)
}
