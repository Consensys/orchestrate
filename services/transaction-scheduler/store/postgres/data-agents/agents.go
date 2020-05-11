package dataagents

import (
	"github.com/go-pg/pg/v9/orm"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/interfaces"
)

type PGAgents struct {
	job       *PGJob
	log       *PGLog
	schedule  *PGSchedule
	txRequest *PGTransactionRequest
}

func New(db orm.DB) *PGAgents {
	return &PGAgents{
		job:       NewPGJob(db),
		log:       NewPGLog(db),
		schedule:  NewPGSchedule(db),
		txRequest: NewPGTransactionRequest(db),
	}
}

func (a *PGAgents) Job() interfaces.JobAgent {
	return a.job
}

func (a *PGAgents) Log() interfaces.LogAgent {
	return a.log
}

func (a *PGAgents) Schedule() interfaces.ScheduleAgent {
	return a.schedule
}

func (a *PGAgents) TransactionRequest() interfaces.TransactionRequestAgent {
	return a.txRequest
}
