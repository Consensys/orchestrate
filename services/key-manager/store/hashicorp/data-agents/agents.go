package dataagents

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/secretstore/services"
)

type HashicorpAgents struct {
	ethereum *HashicorpEthereum
}

func New(vault services.SecretStore) *HashicorpAgents {
	return &HashicorpAgents{
		ethereum: NewHashicorpEthereum(vault),
	}
}

func (a *HashicorpAgents) Ethereum() store.EthereumAgent {
	return a.ethereum
}
