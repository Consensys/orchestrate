package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func renameBlockPositionAndFromBlock(db migrations.DB) error {
	log.Debug("Renaming column block_position")
	_, err := db.Exec(`
	ALTER TABLE chains
	RENAME COLUMN listener_block_position TO listener_current_block;`)

	if err != nil {
		log.WithError(err).Error("Could not rename column listener_block_position")
		return err
	}

	_, err = db.Exec(`
	ALTER TABLE chains
	RENAME COLUMN listener_from_block TO listener_starting_block;`)

	if err != nil {
		log.WithError(err).Error("Could not rename column listener_from_block")
		return err
	}

	log.Info("Renamed columns")

	return nil
}

func renameCurrentBlockAndStartingBlock(db migrations.DB) error {
	log.Debug("Renaming column current_block back")
	_, err := db.Exec(`
	ALTER TABLE chains
	RENAME COLUMN listener_current_block TO listener_block_position;`)
	if err != nil {
		log.WithError(err).Error("Could not rename column listener_current_block")
		return err
	}

	_, err = db.Exec(`
	ALTER TABLE chains
	RENAME COLUMN listener_starting_block TO listener_from_block;`)
	if err != nil {
		log.WithError(err).Error("Could not rename column listener_starting_block")
		return err
	}

	log.Info("Renamed columns")

	return nil
}

func init() {
	Collection.MustRegisterTx(renameBlockPositionAndFromBlock, renameCurrentBlockAndStartingBlock)
}
