package postgres

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database"
	pg "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	dataagents "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/postgres/data-agents"
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
		errMessage := "failed to start postgres DB transaction"
		log.WithError(err).Error(errMessage)
		return nil, errors.PostgresConnectionError(errMessage)
	}

	return &PGTX{
		Tx:       tx,
		PGAgents: dataagents.New(tx),
	}, nil
}

func (pgTx *PGTX) Begin() (database.Tx, error) {
	tx, err := pgTx.Tx.Begin()
	if err != nil {
		errMessage := "failed to start nested postgres DB transaction"
		log.WithError(err).Error(errMessage)
		return nil, errors.PostgresConnectionError(errMessage)
	}

	return &PGTX{
		Tx:       tx,
		PGAgents: dataagents.New(tx),
	}, nil
}

func (pgTx *PGTX) Commit() error {
	err := pgTx.Tx.Commit()
	if err != nil {
		errMessage := "failed to commit postgres DB transaction"
		log.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage)
	}

	return nil
}

func (pgTx *PGTX) Close() error {
	err := pgTx.Tx.Close()
	if err != nil {
		errMessage := "failed to close postgres DB transaction"
		log.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage)
	}

	return nil
}

func (pgTx *PGTX) Rollback() error {
	err := pgTx.Tx.Rollback()
	if err != nil {
		errMessage := "failed to rollback postgres DB transaction"
		log.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage)
	}

	return nil
}
