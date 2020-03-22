package traefik

import (
	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"github.com/containous/traefik/v2/pkg/log"
	traefikprovider "github.com/containous/traefik/v2/pkg/provider"
	"github.com/containous/traefik/v2/pkg/safe"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider/aggregator"
)

// NewProviderAggregator returns an aggregate of all the providers configured in the static configuration.
func NewProvidersAggregator(conf *traefikstatic.Providers, pool *safe.Pool) *aggregator.Provider {
	p := aggregator.New()

	if conf.File != nil {
		initProvider(conf.File)
		p.AddProvider(New(conf.File, pool))
	}

	if conf.Docker != nil {
		initProvider(conf.Docker)
		p.AddProvider(New(conf.Docker, pool))
	}

	if conf.Marathon != nil {
		initProvider(conf.Marathon)
		p.AddProvider(New(conf.Marathon, pool))
	}

	if conf.Rest != nil {
		initProvider(conf.Rest)
		p.AddProvider(New(conf.Rest, pool))
	}

	if conf.KubernetesIngress != nil {
		initProvider(conf.KubernetesIngress)
		p.AddProvider(New(conf.KubernetesIngress, pool))
	}

	if conf.KubernetesCRD != nil {
		initProvider(conf.KubernetesCRD)
		p.AddProvider(New(conf.KubernetesCRD, pool))
	}

	if conf.Rancher != nil {
		initProvider(conf.Rancher)
		p.AddProvider(New(conf.Rancher, pool))
	}

	return p
}

func initProvider(prvdr traefikprovider.Provider) {
	err := prvdr.Init()
	if err != nil {
		log.WithoutContext().WithError(err).Errorf("Error while initializing Traefik provider %T", prvdr)
	}
}
