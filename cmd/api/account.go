package api

import (
	"context"
	"fmt"
	"strings"

	qkm "github.com/consensys/orchestrate/pkg/quorum-key-manager"
	"github.com/consensys/orchestrate/pkg/toolkit/database/postgres"
	qkmClient "github.com/consensys/quorum-key-manager/pkg/client"
	"github.com/go-pg/pg/v9"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newAccountCmd() *cobra.Command {
	var db *pg.DB

	accountCmd := &cobra.Command{
		Use:   "account",
		Short: "Account management",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Set database connection
			opts, err := postgres.NewConfig(viper.GetViper()).PGOptions()
			if err != nil {
				return err
			}
			db = pg.Connect(opts)

			// Init QKM client
			qkm.Init()
			return nil
		},
	}

	// Postgres flags
	postgres.PGFlags(accountCmd.Flags())
	qkm.Flags(accountCmd.Flags())

	importCmd := &cobra.Command{
		Use:   "import",
		Short: "import accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return importAccounts(cmd.Context(), db, qkm.GlobalClient(), qkm.GlobalStoreName())
		},
	}
	accountCmd.AddCommand(importCmd)

	return accountCmd
}

func importAccounts(ctx context.Context, db *pg.DB, client qkmClient.KeyManagerClient, storeName string) error {
	log.Debug("Loading accounts from Vault...")

	accounts, err := client.ListEthAccounts(ctx, storeName, 0, 0)
	if err != nil {
		log.WithError(err).Errorf("could not get list of accounts")
		return err
	}

	var queryInsertItems []string
	for _, accountID := range accounts {
		acc, err2 := client.GetEthAccount(ctx, storeName, accountID)
		if err2 != nil {
			log.WithField("account_id", accountID).WithError(err2).Error("Could not get account")
			return err2
		}

		tenantIDs := strings.Split(acc.Tags[qkm.TagIDAllowedTenants], qkm.TagSeparatorAllowedTenants)
		for _, tenantID := range tenantIDs {
			queryInsertItems = append(queryInsertItems, fmt.Sprintf("('%s', '%s', '%s', '%s', '{\"source\": \"kv-v2\"}')",
				tenantID,
				acc.Address,
				acc.PublicKey,
				acc.CompressedPublicKey,
			))
		}
	}

	if len(queryInsertItems) > 0 {
		_, err = db.Exec("INSERT INTO accounts (tenant_id, address, public_key, compressed_public_key, attributes) VALUES " +
			strings.Join(queryInsertItems, ", ") + " on conflict do nothing")
		if err != nil {
			log.WithError(err).Error("Could not import accounts")
			return err
		}
	}

	log.WithField("accounts", len(queryInsertItems)).Info("accounts imported successfully")
	return nil
}
