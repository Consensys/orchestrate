package registry

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/job"
	"github.com/containous/traefik/v2/pkg/safe"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

// Provider holds configurations of the provider.
type Provider struct {
	Name            string
	RefreshInterval time.Duration
	Registry        types.ChainRegistryStore
}

func NewProvider(name string, registry types.ChainRegistryStore, refreshInterval time.Duration) *Provider {
	return &Provider{
		Name:            name,
		RefreshInterval: refreshInterval,
		Registry:        registry,
	}
}

func (p *Provider) Provide(configurationChan chan<- dynamic.Message, pool *safe.Pool) error {
	pool.GoCtx(func(routineCtx context.Context) {
		operation := func() error {
			ticker := time.NewTicker(p.RefreshInterval)

			for {
				select {
				case <-ticker.C:
					data, err := p.Registry.GetConfig(routineCtx)
					if err != nil {
						log.Errorf("error get chain registry data, %v", err)
						return err
					}

					configuration, err := types.BuildConfiguration(data)
					if err != nil {
						return err
					}

					configurationChan <- dynamic.Message{
						ProviderName:  p.Name,
						Configuration: configuration,
					}
				case <-routineCtx.Done():
					ticker.Stop()
					return nil
				}
			}
		}

		notify := func(err error, time time.Duration) {
			log.Errorf("Provider connection error %+v, retrying in %s", err, time)
		}

		err := backoff.RetryNotify(safe.OperationWithRecover(operation), backoff.WithContext(job.NewBackOff(backoff.NewExponentialBackOff()), routineCtx), notify)
		if err != nil {
			log.Errorf("Cannot connect to registry server %+v", err)
		}
	})

	return nil
}

func (p *Provider) Init() error {
	return nil
}
