package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func updateToStringColumns(db migrations.DB) error {
	log.Debug("Creating tables...")
	_, err := db.Exec(`
ALTER TABLE artifacts
	ALTER COLUMN abi TYPE TEXT USING convert_from(abi,'UTF-8'),
	ALTER COLUMN bytecode TYPE TEXT USING '0x'||encode(bytecode,'hex'),
	ALTER COLUMN deployed_bytecode TYPE TEXT USING '0x'||encode(deployed_bytecode,'hex');

ALTER TABLE methods
	ALTER COLUMN abi TYPE TEXT USING convert_from(abi,'UTF-8');

ALTER TABLE events
	ALTER COLUMN abi TYPE TEXT USING convert_from(abi,'UTF-8');
	`)
	if err != nil {
		log.WithError(err).Error("Could not create tables")
		return err
	}
	log.Info("Created tables")

	return nil
}

func downgradeToByteColumns(db migrations.DB) error {
	log.Debug("Dropping tables")
	_, err := db.Exec(`
ALTER TABLE artifacts
	ALTER COLUMN abi TYPE BYTEA USING abi::TEXT::BYTEA,
	ALTER COLUMN bytecode TYPE BYTEA USING decode(LTRIM(bytecode, '0x'), 'hex'),
	ALTER COLUMN deployed_bytecode TYPE BYTEA USING decode(LTRIM(deployed_bytecode, '0x'), 'hex');

ALTER TABLE methods
	ALTER COLUMN abi TYPE BYTEA USING abi::TEXT::BYTEA;

ALTER TABLE events
	ALTER COLUMN abi TYPE BYTEA USING abi::TEXT::BYTEA;
`)
	if err != nil {
		log.WithError(err).Error("Could not drop tables")
		return err
	}
	log.Info("Dropped tables")

	return nil
}

func init() {
	Collection.MustRegisterTx(updateToStringColumns, downgradeToByteColumns)
}
