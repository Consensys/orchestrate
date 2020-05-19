package dataagents

import (
	pg "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
)

type PGAgents struct {
	tx        *PGTransaction
	job       *PGJob
	log       *PGLog
	schedule  *PGSchedule
	txRequest *PGTransactionRequest
}

func New(db pg.DB) *PGAgents {
	return &PGAgents{
		tx:        NewPGTransaction(db),
		job:       NewPGJob(db),
		log:       NewPGLog(db),
		schedule:  NewPGSchedule(db),
		txRequest: NewPGTransactionRequest(db),
	}
}

func (a *PGAgents) Job() store.JobAgent {
	return a.job
}

func (a *PGAgents) Log() store.LogAgent {
	return a.log
}

func (a *PGAgents) Schedule() store.ScheduleAgent {
	return a.schedule
}

func (a *PGAgents) Transaction() store.TransactionAgent {
	return a.tx
}

func (a *PGAgents) TransactionRequest() store.TransactionRequestAgent {
	return a.txRequest
}
