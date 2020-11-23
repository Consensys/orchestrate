package dataagents

import (
	pg "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/identity-manager/store"
)

type PGAgents struct {
	identity *PGAccount
}

func New(db pg.DB) *PGAgents {
	return &PGAgents{
		identity: NewPGAccount(db),
	}
}

func (a *PGAgents) Account() store.AccountAgent {
	return a.identity
}
