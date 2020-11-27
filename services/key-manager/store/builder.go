package store

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/hashicorp"
)

func Build(ctx context.Context, cfg *Config) (Vault, error) {
	switch cfg.Type {
	case azureKeyVaultType:
		// TODO: Configure azure key vault
		return nil, errors.ConfigError("Azure key vault support not implemented yet")
	case hashicorpVaultType:
		return hashicorp.NewOrchestrateVaultClient(hashicorp.ConfigFromViper())
	case ukcVaultType:
		// TODO: Configure Unbound key vault
		return nil, errors.ConfigError("UKC key vault support not implemented yet")
	default:
		errMessage := "invalid key manager vault type"
		log.WithContext(ctx).WithField("vault_type", cfg.Type).Error(errMessage)
		return nil, errors.ConfigError(" %q", cfg.Type)
	}
}
