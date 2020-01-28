package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func createContextTable(db migrations.DB) error {
	log.Debug("Creating tables...")
	_, err := db.Exec(`
CREATE TABLE repositories (
	id SERIAL PRIMARY KEY,
	name VARCHAR(66) NOT NULL,
	UNIQUE(name)
);

CREATE TABLE artifacts (
    id SERIAL PRIMARY KEY,
    abi BYTEA,
    bytecode BYTEA,
    deployed_bytecode BYTEA,
    codehash CHAR(66)
);
CREATE UNIQUE INDEX unique_abi_bytecode ON artifacts (md5(abi), md5(deployed_bytecode));

CREATE TABLE tags (
	id SERIAL PRIMARY KEY,
	name VARCHAR(66) NOT NULL,
	repository_id INTEGER REFERENCES repositories,
	artifact_id INTEGER REFERENCES artifacts,
	UNIQUE(name, repository_id)
);

CREATE TABLE codehashes (
	id SERIAL PRIMARY KEY,
	chain_id VARCHAR(66) NOT NULL,
	address CHAR(42) NOT NULL,
	codehash CHAR(66) NOT NULL,
	UNIQUE(chain_id, address)
);

CREATE TABLE methods (
	id SERIAL PRIMARY KEY,
	codehash CHAR(66),
	selector CHAR(10) NOT NULL,
	abi BYTEA NOT NULL
);

CREATE TABLE events (
	id SERIAL PRIMARY KEY,
	codehash CHAR(66),
	sig_hash CHAR(66) NOT NULL,
	indexed_input_count INTEGER NOT NULL,
	abi BYTEA NOT NULL
);
CREATE INDEX ON events (sig_hash, indexed_input_count, codehash);
	`)
	if err != nil {
		log.WithError(err).Error("Could not create tables")
		return err
	}
	log.Info("Created tables")

	return nil
}

func dropContextTable(db migrations.DB) error {
	log.Debug("Dropping tables")
	_, err := db.Exec(`
DROP TABLE events;
DROP TABLE methods;
DROP TABLE codehashes;
DROP TABLE tags;
DROP TABLE artifacts;
DROP TABLE repositories;
`)
	if err != nil {
		log.WithError(err).Error("Could not drop tables")
		return err
	}
	log.Info("Dropped tables")

	return nil
}

func init() {
	Collection.MustRegisterTx(createContextTable, dropContextTable)
}
