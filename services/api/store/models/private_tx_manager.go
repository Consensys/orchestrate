package models

import (
	"time"

	"github.com/ConsenSys/orchestrate/pkg/types/entities"
)

type PrivateTxManager struct {
	tableName struct{} `pg:"private_tx_managers"` // nolint:unused,structcheck // reason

	UUID      string `pg:",pk"`
	ChainUUID string
	URL       string
	Type      entities.PrivateTxManagerType
	CreatedAt time.Time `pg:",default:now()"`
}
