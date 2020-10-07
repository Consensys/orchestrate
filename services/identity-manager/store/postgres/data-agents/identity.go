package dataagents

import (
	pg "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
)

// const identityDAComponent = "data-agents.identity"

// PGIdentity is an Identity data agent for PostgreSQL
type PGIdentity struct {
	db pg.DB
}

// NewPGIdentity creates a new PGIdentity
func NewPGIdentity(db pg.DB) *PGIdentity {
	return &PGIdentity{db: db}
}
