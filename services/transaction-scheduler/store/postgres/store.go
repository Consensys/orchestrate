package postgres

import (
	"github.com/go-pg/pg/v9"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/interfaces"
	dataagents "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/postgres/data-agents"
)

type PGDB struct {
	*pg.DB
	*dataagents.PGAgents
}

func NewPGDB(db *pg.DB) *PGDB {
	return &PGDB{
		DB:       db,
		PGAgents: dataagents.New(db),
	}
}

func (db *PGDB) Begin() (interfaces.Tx, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		errMessage := "failed to start postgres DB transaction"
		log.WithError(err).Error(errMessage)
		return nil, errors.PostgresConnectionError(errMessage)
	}

	return &pgTx{
		Tx:       tx,
		PGAgents: dataagents.New(tx),
	}, nil
}
