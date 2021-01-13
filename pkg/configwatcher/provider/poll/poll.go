package poll

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/containous/traefik/v2/pkg/job"
	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/configwatcher/provider"
)

type Provider struct {
	poll    func(ctx context.Context) (provider.Message, error)
	refresh time.Duration
}

func New(poll func(ctx context.Context) (provider.Message, error), refresh time.Duration) *Provider {
	return &Provider{
		poll:    poll,
		refresh: refresh,
	}
}

func (p *Provider) Provide(ctx context.Context, msgs chan<- provider.Message) error {
	logger := log.FromContext(ctx).WithField("provider", fmt.Sprintf("%T", p))
	logger.Info("start providing")

	operation := func() error {
		ticker := time.NewTicker(p.refresh)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				msg, err := p.poll(ctx)
				if err != nil {
					return err
				}

				msgs <- msg
			case <-ctx.Done():
				logger.Infof("stopped providing")
				return nil
			}
		}
	}

	notify := func(err error, time time.Duration) {
		logger.WithError(err).Warnf("error while providing (retrying in %s)", time)
	}

	return backoff.RetryNotify(operation, backoff.WithContext(job.NewBackOff(backoff.NewExponentialBackOff()), ctx), notify)
}
