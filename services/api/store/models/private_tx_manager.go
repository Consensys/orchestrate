package models

import "time"

type PrivateTxManager struct {
	tableName struct{} `pg:"private_tx_managers"` // nolint:unused,structcheck // reason

	UUID      string `pg:",pk"`
	ChainUUID string
	URL       string
	Type      string
	CreatedAt time.Time `pg:",default:now()"`
}
