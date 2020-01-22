package chainregistry

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/containous/traefik/v2/pkg/job"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/containous/traefik/v2/pkg/safe"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
)

type Provider struct {
	Client chainregistry.Client
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
			var nodes []*types.Node
			nodes, err = p.Client.GetNodes(ctx)
			if err != nil {
				log.FromContext(logCtx).WithError(err).Errorf("failed to fetch nodes from chain registry")
				break loop
			}
			configInput <- p.buildConfiguration(nodes)
		case <-logCtx.Done():
		}
	}

	return
}

func (p *Provider) buildConfiguration(nodes []*types.Node) *dynamic.Message {
	msg := &dynamic.Message{
		Provider: "chain-registry",
		Configuration: &dynamic.Configuration{
			Nodes: make(map[string]*dynamic.Node),
		},
	}

	for _, node := range nodes {
		duration, err := time.ParseDuration(node.ListenerBackOffDuration)
		if err != nil {
			log.Errorf("cannot parse duration for node ID:%s - TenantID:%s - Name:%s", node.ID, node.TenantID, node.Name)
		}

		msg.Configuration.Nodes[node.ID] = &dynamic.Node{
			ID:       node.ID,
			TenantID: node.TenantID,
			Name:     node.Name,
			URL:      fmt.Sprintf("%v/%v", p.conf.ChainRegistryURL, node.ID),
			Listener: &dynamic.Listener{
				BlockPosition: node.ListenerBlockPosition,
				Depth:         node.ListenerDepth,
				Backoff:       duration,
			},
		}
	}

	return msg
}
