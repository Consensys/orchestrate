package models

import (
	"time"
)

type Schedule struct {
	tableName struct{} `pg:"schedules"` // nolint:unused,structcheck // reason

	ID                   int `pg:"alias:id"`
	UUID                 string
	TenantID             string `pg:"alias:tenant_id"`
	ChainUUID            string
	Jobs                 []*Job
	TransactionRequestID *int
	TransactionRequest   *TransactionRequest
	CreatedAt            time.Time `pg:"default:now()"`
}
