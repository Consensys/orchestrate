package dataagents

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/secretstore/services"
)

const ethereumDAComponent = "data-agents.ethereum"

// HashicorpEthereum is a job data agent for Ethereum Hashicorp Vault
type HashicorpEthereum struct {
	vault       services.SecretStore
	generateKey func(address, namespace string) string
}

// NewHashicorpEthereum creates a new HashicorpEthereum
func NewHashicorpEthereum(vault services.SecretStore) *HashicorpEthereum {
	return &HashicorpEthereum{vault: vault, generateKey: generateKey}
}

func (agent *HashicorpEthereum) Insert(ctx context.Context, address, privKey, namespace string) error {
	key := agent.generateKey(address, namespace)
	err := agent.vault.Store(ctx, key, privKey)
	if err != nil {
		errMessage := "failed to store privateKey in Hashicorp Vault"
		log.WithContext(ctx).WithError(err).WithField("key", key).Error(errMessage)
		return errors.HashicorpVaultConnectionError(errMessage).ExtendComponent(ethereumDAComponent)
	}

	return nil
}

func generateKey(address, namespace string) string {
	key := address
	if namespace != "" {
		key = fmt.Sprintf("%s_%s", namespace, address)
	}

	return key
}
