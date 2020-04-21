package models

import "time"

type PrivateTxManagerModel struct {
	tableName struct{} `pg:"private_tx_manager"` // nolint:unused,structcheck // reason

	ID        string `pg:",pk"`
	ChainUUID string
	URL       string
	Type      string
	CreatedAt *time.Time
}
