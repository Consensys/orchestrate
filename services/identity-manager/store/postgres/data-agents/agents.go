package dataagents

import (
	pg "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/store"
)

type PGAgents struct {
	identity *PGIdentity
}

func New(db pg.DB) *PGAgents {
	return &PGAgents{
		identity: NewPGIdentity(db),
	}
}

func (a *PGAgents) Identity() store.IdentityAgent {
	return a.identity
}
