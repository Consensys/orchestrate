package multi

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/store"
)

func Build(ctx context.Context, cfg *Config) (store.Vault, error) {
	switch cfg.Type {
	case azureKeyVaultType:
		// TODO: Configure azure key vault
		log.WithContext(ctx).Info("Azure key vault initialized")
		return nil, nil
	case hashicorpVaultType:
		// TODO: Configure hashicorp Vault
		log.WithContext(ctx).Info("Hashicorp vault initialized")
		return nil, nil
	case unboundType:
		// TODO: Configure Unbound key vault
		log.WithContext(ctx).Info("Unbound key vault initialized")
		return nil, nil
	default:
		return nil, fmt.Errorf("invalid tx-signer vault type %q", cfg.Type)
	}
}
