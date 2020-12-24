package dataagents

import (
	pg "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
)

type PGAgents struct {
	tx        *PGTransaction
	job       *PGJob
	log       *PGLog
	schedule  *PGSchedule
	txRequest *PGTransactionRequest
	account   *PGAccount
	faucet    *PGFaucet
}

func New(db pg.DB) *PGAgents {
	return &PGAgents{
		tx:        NewPGTransaction(db),
		job:       NewPGJob(db),
		log:       NewPGLog(db),
		schedule:  NewPGSchedule(db),
		txRequest: NewPGTransactionRequest(db),
		account:   NewPGAccount(db),
		faucet:    NewPGFaucet(db),
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

func (a *PGAgents) Account() store.AccountAgent {
	return a.account
}

func (a *PGAgents) Faucet() store.FaucetAgent {
	return a.faucet
}
