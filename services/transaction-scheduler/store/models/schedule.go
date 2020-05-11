package models

import "time"

type Schedule struct {
	tableName struct{} `pg:"schedules"` // nolint:unused,structcheck // reason

	ID        int
	UUID      string
	TenantID  string
	ChainUUID string
	Jobs      []*Job
	CreatedAt time.Time `pg:"default:now()"`
}
