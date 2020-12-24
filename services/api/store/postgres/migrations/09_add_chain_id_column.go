package migrations

import (
	"context"

	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/chain-registry/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/models"
)

func addChainIDColumn(db migrations.DB) error {
	log.Debugf("Adding chainID column on table %q...", "chains")

	_, err := db.Exec(`
ALTER TABLE chains
	ALTER COLUMN listener_depth SET DEFAULT 0,
	ALTER COLUMN listener_current_block SET DEFAULT 0,
	ALTER COLUMN listener_starting_block SET DEFAULT 0;

ALTER TABLE chains
	ALTER COLUMN listener_depth SET NOT NULL,
	ALTER COLUMN listener_current_block SET NOT NULL,
	ALTER COLUMN listener_starting_block SET NOT NULL;

ALTER TABLE chains
	ADD COLUMN chain_id BIGINT NOT NULL DEFAULT 0;
	`)

	if err != nil {
		return err
	}

	log.Infof("Added chainID columns on table %q", "chains")

	err = updateChainIDs(context.Background(), db)
	if err != nil {
		return err
	}

	return nil
}

func dropChainIDColumn(db migrations.DB) error {
	log.Debugf("Removing chainID chainID on table %q...", "chains")

	_, err := db.Exec(`
ALTER TABLE chains 
	DROP COLUMN chain_id;

ALTER TABLE chains
	ALTER COLUMN listener_depth DROP DEFAULT,
	ALTER COLUMN listener_current_block DROP DEFAULT,
	ALTER COLUMN listener_starting_block DROP DEFAULT;

ALTER TABLE chains
	ALTER COLUMN listener_depth DROP NOT NULL,
	ALTER COLUMN listener_current_block DROP NOT NULL,
	ALTER COLUMN listener_starting_block DROP NOT NULL;
	`)

	if err != nil {
		return err
	}

	log.Infof("Removed chainID column on table %q", "chains")

	return nil
}

func init() { Collection.MustRegisterTx(addChainIDColumn, dropChainIDColumn) }

func updateChainIDs(ctx context.Context, db migrations.DB) error {
	ethclient.Init(ctx)
	ec := ethclient.GlobalClient()

	log.Debugf("fetching chainIDs from rpc nodes")

	var chains []*models.Chain
	err := db.Model(&chains).Where(`chain_id = ?`, 0).Select()
	if err != nil {
		return err
	}

	for _, chain := range chains {
		chainID, err := utils.GetChainID(ctx, ec, chain.URLs)
		if err != nil {
			return err
		}

		_, err = db.Model(&models.Chain{ChainID: chainID.String()}).
			Where("uuid = ?", chain.UUID).UpdateNotZero()

		if err != nil {
			return err
		}

		log.WithField("chainName", chain.Name).
			WithField("chainUUID", chain.UUID).
			WithField("chainID", chainID).
			Infof("chain was updated correctly")
	}

	return nil
}
