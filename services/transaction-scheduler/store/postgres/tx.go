package postgres

import (
	"github.com/go-pg/pg/v9"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	dataagents "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/postgres/data-agents"
)

type pgTx struct {
	Tx *pg.Tx
	*dataagents.PGAgents
}

func (pgTx *pgTx) Commit() error {
	err := pgTx.Tx.Commit()
	if err != nil {
		errMessage := "failed to commit postgres DB transaction"
		log.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage)
	}

	return nil
}
