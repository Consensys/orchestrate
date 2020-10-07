package store

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database"
)

//go:generate mockgen -source=store.go -destination=mocks/mock.go -package=mocks

type Agents interface {
	Identity() IdentityAgent
}

type DB interface {
	database.DB
	Agents
}

// Interfaces data agents
type IdentityAgent interface{}
