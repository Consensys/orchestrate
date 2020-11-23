package multi

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/secretstore/hashicorp"
	hashicorpstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/store/hashicorp"

	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/store"
)

func Build(ctx context.Context, cfg *Config) (store.Vault, error) {
	switch cfg.Type {
	case azureKeyVaultType:
		// TODO: Configure azure key vault
		return nil, errors.ConfigError("Azure key vault support not implemented yet")
	case hashicorpVaultType:
		// TODO: Refactor Initialization and move to this MS or pkg if reused
		hashicorp.Init(ctx)
		vault := hashicorpstore.NewHashicorpVault(hashicorp.GlobalStore(), hashicorp.GlobalChecker())

		log.WithContext(ctx).Info("Hashicorp vault initialized")
		return vault, nil
	case ukcVaultType:
		// TODO: Configure Unbound key vault
		return nil, errors.ConfigError("UKC key vault support not implemented yet")
	default:
		return nil, errors.ConfigError("invalid key manager vault type %q", cfg.Type)
	}
}
