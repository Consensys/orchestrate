package traefik

import (
	"context"
	"fmt"

	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/log"
	traefikprovider "github.com/containous/traefik/v2/pkg/provider"
	"github.com/containous/traefik/v2/pkg/safe"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/configwatcher/provider"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/config/dynamic"
)

type Provider struct {
	prvdr traefikprovider.Provider
	pool  *safe.Pool
}

func New(prvdr traefikprovider.Provider, pool *safe.Pool) *Provider {
	return &Provider{
		prvdr: prvdr,
		pool:  pool,
	}
}

func (p *Provider) Provide(ctx context.Context, msgs chan<- provider.Message) error {
	logger := log.FromContext(ctx).WithField("provider", fmt.Sprintf("%T", p.prvdr))
	logger.Infof("start providing")
	// We can not close pipedMsgs due to Traefik implementation that does not allow to
	// determine when to close pipedMsgs without risking to have Traefik provider
	// sending messages to channel after closing
	pipedMsgs := make(chan traefikdynamic.Message)

	errors := make(chan error)
	defer close(errors)

	go func() {
		err := p.prvdr.Provide(pipedMsgs, p.pool)
		if err != nil {
			errors <- err
		}
	}()

	for {
		select {
		case <-ctx.Done():
			logger.Infof("stop providing")
			return nil
		case err := <-errors:
			logger.WithError(err).Errorf("stop providing")
			return err
		case msg := <-pipedMsgs:
			msgs <- dynamic.FromTraefikMessage(&msg)
		}
	}
}
