package postgres

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database"
	pg "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	dataagents "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/postgres/data-agents"
)

type PGDB struct {
	pg.DB
	*dataagents.PGAgents
}

type PGTX struct {
	pg.Tx
	*dataagents.PGAgents
}

func NewPGDB(db pg.DB) *PGDB {
	return &PGDB{
		DB:       db,
		PGAgents: dataagents.New(db),
	}
}

func (db *PGDB) Begin() (database.Tx, error) {
	db.Transaction()
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, errors.PostgresConnectionError("failed to start postgres DB transaction")
	}

	return &PGTX{
		Tx:       tx,
		PGAgents: dataagents.New(tx),
	}, nil
}

func (pgTx *PGTX) Begin() (database.Tx, error) {
	tx, err := pgTx.Tx.Begin()
	if err != nil {
		return nil, errors.PostgresConnectionError("failed to start nested postgres DB transaction")
	}

	return &PGTX{
		Tx:       tx,
		PGAgents: dataagents.New(tx),
	}, nil
}

func (pgTx *PGTX) Commit() error {
	err := pgTx.Tx.Commit()
	if err != nil {
		return errors.PostgresConnectionError("failed to commit postgres DB transaction")
	}

	return nil
}

func (pgTx *PGTX) Close() error {
	err := pgTx.Tx.Close()
	if err != nil {
		return errors.PostgresConnectionError("failed to close postgres DB transaction")
	}

	return nil
}

func (pgTx *PGTX) Rollback() error {
	err := pgTx.Tx.Rollback()
	if err != nil {
		return errors.PostgresConnectionError("failed to rollback postgres DB transaction")
	}

	return nil
}
