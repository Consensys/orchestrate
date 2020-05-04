package models

import "time"

type Schedule struct {
	tableName struct{} `pg:"schedules"` // nolint:unused,structcheck // reason

	ID        int
	UUID      string
	TenantID  string
	ChainID   string
	CreatedAt time.Time `pg:"default:now()"`
}
