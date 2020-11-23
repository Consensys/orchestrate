package hashicorp

import (
	healthz "github.com/heptiolabs/healthcheck"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/secretstore"
	dataagents "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/store/hashicorp/data-agents"
)

type Vault struct {
	*dataagents.HashicorpAgents
	healthChecker healthz.Check
}

func NewHashicorpVault(vault secretstore.SecretStore, healthChecker healthz.Check) *Vault {
	return &Vault{
		HashicorpAgents: dataagents.New(vault),
		healthChecker:   healthChecker,
	}
}

func (vault *Vault) HealthCheck() healthz.Check {
	return vault.healthChecker
}
