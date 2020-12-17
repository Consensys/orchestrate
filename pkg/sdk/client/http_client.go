package client

import (
	"context"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
)

type HTTPClient struct {
	client *http.Client
	config *Config
}

func NewHTTPClient(h *http.Client, c *Config) OrchestrateClient {
	return &HTTPClient{
		client: h,
		config: c,
	}
}

func callWithBackOff(ctx context.Context, backOff backoff.BackOff, requestCall func() error) error {
	return backoff.RetryNotify(
		func() error {
			err := requestCall()
			// If not errors, it does not retry
			if err == nil {
				return nil
			}

			// Retry on following errors
			if errors.IsInvalidStateError(err) || errors.IsServiceConnectionError(err) {
				return err
			}

			// Otherwise, stop retrying
			return backoff.Permanent(err)
		}, backoff.WithContext(backOff, ctx),
		func(e error, duration time.Duration) {
			log.FromContext(ctx).
				WithError(e).
				Warnf("transaction-scheduler: http call failed, retrying in %v...", duration)
		},
	)
}
