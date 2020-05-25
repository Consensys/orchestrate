package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func createContextTable(db migrations.DB) error {
	log.Debug("Creating tables...")
	_, err := db.Exec(`
CREATE TABLE transactions (
	id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL,
	hash TEXT,
	sender TEXT,
	recipient TEXT,
	nonce INTEGER,
	value BIGINT,
	gas_price BIGINT,
	gas_limit INTEGER,
	data TEXT,
	raw TEXT,
	private_from TEXT,
	private_for TEXT [],
	privacy_group_id TEXT,
	created_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL, 
	updated_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL,
	UNIQUE(uuid),
	UNIQUE(hash)
);

CREATE TABLE transaction_requests (
    id SERIAL PRIMARY KEY,
    idempotency_key TEXT NOT NULL,
	request_hash TEXT NOT NULL,
    params jsonb NOT NULL,
	created_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL, 
	UNIQUE(idempotency_key)
);

CREATE TABLE schedules (
    id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL,
	tenant_id TEXT NOT NULL,
	chain_uuid UUID NOT NULL,
	transaction_request_id INTEGER REFERENCES transaction_requests(id),
	created_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL, 
	UNIQUE(uuid)
);

CREATE TABLE jobs (
    id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL,
	schedule_id INTEGER NOT NULL REFERENCES schedules,
    type TEXT NOT NULL,
    transaction_id INTEGER NOT NULL REFERENCES transactions(id),
	labels jsonb,
	created_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL,
	UNIQUE(uuid)
);

CREATE TABLE logs (
    id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL,
    job_id INTEGER NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
	status TEXT NOT NULL,
	message TEXT,
	created_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL,
	UNIQUE(uuid)
);
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
DROP TABLE logs;
DROP TABLE jobs;
DROP TABLE schedules;
DROP TABLE transaction_requests;
DROP TABLE transactions;
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
