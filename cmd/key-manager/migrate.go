package keymanager

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	kvv2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/hashicorp/kv-v2"
	keymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/store"
	migrations2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/store/migrations"
)

// newMigrateCmd create migrate command
func newMigrateCmd() *cobra.Command {
	var vault store.Vault
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migration of Vault secrets",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cfg := keymanager.NewConfig(viper.GetViper())
			var err error
			vault, err = store.Build(cmd.Context(), cfg.Store)
			if err != nil {
				return err
			}
			return vault.HealthCheck()
		},
	}

	// Register Init command
	importSecretCmd := &cobra.Command{
		Use:   "import-secrets",
		Short: "Import secrets store in old Hashicorp vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Initialize v2 client
			cfg := kvv2.ConfigFromViper()
			v2Client, err := kvv2.NewClient(cfg.Config, cfg.SecretPath)
			if err != nil {
				return err
			}

			return migrations2.Kvv2ImportSecrets(cmd.Context(), vault, v2Client)
		},
	}

	kvv2.InitFlags(importSecretCmd.Flags())
	migrateCmd.AddCommand(importSecretCmd)

	return migrateCmd
}
