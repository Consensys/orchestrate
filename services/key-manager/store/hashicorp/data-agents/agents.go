package dataagents

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/secretstore"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/store"
)

type HashicorpAgents struct {
	ethereum *HashicorpEthereum
}

func New(vault secretstore.SecretStore) *HashicorpAgents {
	return &HashicorpAgents{
		ethereum: NewHashicorpEthereum(vault),
	}
}

func (a *HashicorpAgents) Ethereum() store.EthereumAgent {
	return a.ethereum
}
