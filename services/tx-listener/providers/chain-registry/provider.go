package chainregistry

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/containous/traefik/v2/pkg/job"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/containous/traefik/v2/pkg/safe"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
)

type Provider struct {
	Client chainregistry.ChainRegistryClient
	conf   *Config
}

func (p *Provider) Run(ctx context.Context, configInput chan<- *dynamic.Message) error {
	return p.runWithRetry(
		ctx,
		configInput,
		backoff.WithContext(job.NewBackOff(backoff.NewExponentialBackOff()), ctx),
	)
}

func (p *Provider) runWithRetry(ctx context.Context, configInput chan<- *dynamic.Message, bckff backoff.BackOff) error {
	return backoff.RetryNotify(
		safe.OperationWithRecover(func() error {
			return p.run(ctx, configInput)
		}),
		bckff,
		func(err error, d time.Duration) {
			log.FromContext(ctx).WithError(err).Warnf("provider restarting in... %v", d)
		},
	)
}

func (p *Provider) run(ctx context.Context, configInput chan<- *dynamic.Message) (err error) {
	logCtx := log.With(ctx, log.Str("provider", "chain-registry"))
	if p.Client == nil {
		return errors.InternalError("client not initialized")
	}

	ticker := time.NewTicker(p.conf.RefreshInterval)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-ticker.C:
			var chains []*models.Chain
			chains, err = p.Client.GetChains(ctx)
			if err != nil {
				log.FromContext(logCtx).WithError(err).Errorf("failed to fetch chains from chain registry")
				break loop
			}
			configInput <- p.buildConfiguration(ctx, chains)
		case <-logCtx.Done():
		}
	}

	return
}

func (p *Provider) buildConfiguration(ctx context.Context, chains []*models.Chain) *dynamic.Message {
	msg := &dynamic.Message{
		Provider: "chain-registry",
		Configuration: &dynamic.Configuration{
			Chains: make(map[string]*dynamic.Chain),
		},
	}

	for _, chain := range chains {
		duration, err := time.ParseDuration(*chain.ListenerBackOffDuration)
		if err != nil {
			log.FromContext(ctx).WithFields(logrus.Fields{
				"TenantID": chain.TenantID,
				"Name":     chain.Name,
			}).Errorf("cannot parse duration: %s", *chain.ListenerBackOffDuration)
		}

		msg.Configuration.Chains[chain.UUID] = &dynamic.Chain{
			UUID:     chain.UUID,
			TenantID: chain.TenantID,
			Name:     chain.Name,
			URL:      fmt.Sprintf("%v/%v", p.conf.ChainRegistryURL, chain.UUID),
			ChainID:  chain.ChainID,
			Listener: &dynamic.Listener{
				StartingBlock:     *chain.ListenerStartingBlock,
				CurrentBlock:      *chain.ListenerCurrentBlock,
				Depth:             *chain.ListenerDepth,
				Backoff:           duration,
				ExternalTxEnabled: *chain.ListenerExternalTxEnabled,
			},
		}
	}

	return msg
}
