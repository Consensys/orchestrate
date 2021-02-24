package traefik

import (
	"context"
	"fmt"

	"github.com/ConsenSys/orchestrate/pkg/configwatcher/provider"
	"github.com/ConsenSys/orchestrate/pkg/http/config/dynamic"
	"github.com/ConsenSys/orchestrate/pkg/log"
	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
	traefikprovider "github.com/containous/traefik/v2/pkg/provider"
	"github.com/containous/traefik/v2/pkg/safe"
)

type Provider struct {
	prvdr  traefikprovider.Provider
	pool   *safe.Pool
	logger *log.Logger
}

func New(prvdr traefikprovider.Provider, pool *safe.Pool) *Provider {
	return &Provider{
		prvdr: prvdr,
		pool:  pool,
		logger: log.NewLogger().SetComponent("configwatcher").
			WithField("provider", fmt.Sprintf("%T", prvdr)),
	}
}

func (p *Provider) Provide(ctx context.Context, msgs chan<- provider.Message) error {
	p.logger.Debug("start providing")
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
			p.logger.Infof("stop providing")
			return nil
		case err := <-errors:
			p.logger.WithError(err).Errorf("stop providing")
			return err
		case msg := <-pipedMsgs:
			msgs <- dynamic.FromTraefikMessage(&msg)
		}
	}
}
