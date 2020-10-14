package hashicorp

import (
	healthz "github.com/heptiolabs/healthcheck"
	dataagents "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/store/hashicorp/data-agents"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/secretstore/services"
)

type Vault struct {
	*dataagents.HashicorpAgents
	healthChecker healthz.Check
}

func NewHashicorpVault(vault services.SecretStore, healthChecker healthz.Check) *Vault {
	return &Vault{
		HashicorpAgents: dataagents.New(vault),
		healthChecker:   healthChecker,
	}
}

func (vault *Vault) HealthCheck() healthz.Check {
	return vault.healthChecker
}
