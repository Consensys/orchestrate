package models

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
)

type PrivateTxManager struct {
	tableName struct{} `pg:"private_tx_managers"` // nolint:unused,structcheck // reason

	UUID      string `pg:",pk"`
	ChainUUID string
	URL       string
	Type      entities.PrivateTxManagerType
	CreatedAt time.Time `pg:",default:now()"`
}
