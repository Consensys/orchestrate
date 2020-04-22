package models

import "time"

type PrivateTxManagerModel struct {
	tableName struct{} `pg:"private_tx_managers"` // nolint:unused,structcheck // reason

	UUID      string     `pg:",pk" validate:"omitempty,uuid4"`
	ChainUUID string     `pg:",type:uuid,alias:chain_uuid,notnull" validate:"omitempty,uuid4"`
	URL       string     `json:"url" validate:"required,url"`
	Type      string     `json:"type" validate:"required,isPrivateTxManagerType"`
	CreatedAt *time.Time `pg:",default:now()"`
}
