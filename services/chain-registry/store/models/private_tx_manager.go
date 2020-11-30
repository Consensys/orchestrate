package models

import "time"

type PrivateTxManagerModel struct {
	tableName struct{} `pg:"private_tx_managers"` // nolint:unused,structcheck // reason

	UUID      string     `pg:",pk" validate:"omitempty,uuid4" example:"b4374e6f-b28a-4bad-b4fe-bda36eaf849c"`
	ChainUUID string     `pg:",type:uuid,alias:chain_uuid,notnull" validate:"omitempty,uuid4" example:"b4374e6f-b28a-4bad-b4fe-bda36eaf849c"`
	URL       string     `json:"url" validate:"required,url" example:"http://tessera:3000"`
	Type      string     `json:"type" validate:"required,isPrivateTxManagerType" example:"Tessera"`
	CreatedAt *time.Time `pg:",default:now()" example:"2020-07-09T12:35:42.115395Z"`
}
