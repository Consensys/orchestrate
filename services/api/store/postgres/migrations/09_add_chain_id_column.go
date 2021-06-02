package migrations

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	ethclient "github.com/ConsenSys/orchestrate/pkg/toolkit/ethclient/rpc"
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

// Minimum version of model to perform migration actions
type chain struct {
	tableName struct{} `pg:"chains"` // nolint

	UUID     string `pg:",pk"`
	Name     string
	TenantID string
	ChainID  string
	URLs     []string `pg:"urls,array"`
}

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

	var chains []*chain
	err := db.Model(&chains).Where(`chain_id = ?`, 0).Select()
	if err != nil {
		return err
	}

	for _, c := range chains {
		chainID, err := getChainID(ctx, ec, c.URLs)
		if err != nil {
			return err
		}

		_, err = db.Model(&chain{ChainID: chainID}).
			Where("uuid = ?", c.UUID).UpdateNotZero()

		if err != nil {
			return err
		}

		log.WithField("chainName", c.Name).
			WithField("chainUUID", c.UUID).
			WithField("chainID", chainID).
			Infof("chain was updated correctly")
	}

	return nil
}

func getChainID(ctx context.Context, ethClient *ethclient.Client, uris []string) (string, error) {
	var prevChainID string
	for i, uri := range uris {
		chainID, err := ethClient.Network(ctx, uri)
		if err != nil {
			errMessage := "failed to fetch chain id"
			log.WithContext(ctx).WithField("url", uri).WithError(err).Error(errMessage)
			return "", errors.InvalidParameterError(errMessage)
		}

		if i > 0 && chainID.String() != prevChainID {
			errMessage := "URLs in the list point to different networks"
			log.WithContext(ctx).
				WithField("url", uri).
				WithField("previous_chain_id", prevChainID).
				WithField("chain_id", chainID.String()).
				Error(errMessage)
			return "", errors.InvalidParameterError(errMessage)
		}

		prevChainID = chainID.String()
	}

	return prevChainID, nil
}
