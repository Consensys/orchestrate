package dataagents

import (
	pg "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
)

type PGAgents struct {
	tx               store.TransactionAgent
	job              store.JobAgent
	log              store.LogAgent
	schedule         store.ScheduleAgent
	txRequest        store.TransactionRequestAgent
	account          store.AccountAgent
	faucet           store.FaucetAgent
	artifact         store.ArtifactAgent
	codeHash         store.CodeHashAgent
	event            store.EventAgent
	method           store.MethodAgent
	repository       store.RepositoryAgent
	tag              store.TagAgent
	chain            store.ChainAgent
	privateTxManager store.PrivateTxManagerAgent
}

func New(db pg.DB) *PGAgents {
	return &PGAgents{
		tx:               NewPGTransaction(db),
		job:              NewPGJob(db),
		log:              NewPGLog(db),
		schedule:         NewPGSchedule(db),
		txRequest:        NewPGTransactionRequest(db),
		account:          NewPGAccount(db),
		faucet:           NewPGFaucet(db),
		artifact:         NewPGArtifact(db),
		codeHash:         NewPGCodeHash(db),
		event:            NewPGEvent(db),
		method:           NewPGMethod(db),
		repository:       NewPGRepository(db),
		tag:              NewPGTag(db),
		chain:            NewPGChain(db),
		privateTxManager: NewPGPrivateTxManager(db),
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

func (a *PGAgents) Artifact() store.ArtifactAgent {
	return a.artifact
}

func (a *PGAgents) CodeHash() store.CodeHashAgent {
	return a.codeHash
}

func (a *PGAgents) Event() store.EventAgent {
	return a.event
}

func (a *PGAgents) Method() store.MethodAgent {
	return a.method
}

func (a *PGAgents) Repository() store.RepositoryAgent {
	return a.repository
}

func (a *PGAgents) Tag() store.TagAgent {
	return a.tag
}

func (a *PGAgents) Chain() store.ChainAgent {
	return a.chain
}

func (a *PGAgents) PrivateTxManager() store.PrivateTxManagerAgent {
	return a.privateTxManager
}
