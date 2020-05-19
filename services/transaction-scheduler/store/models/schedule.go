package models

import (
	"time"
)

var (
	ScheduleStatic = Schedule{}
)

type Schedule struct {
	tableName struct{} `pg:"schedules"` // nolint:unused,structcheck // reason

	ID        int `pg:"alias:id"`
	UUID      string
	TenantID  string `pg:"alias:tenant_id"`
	ChainUUID string
	Jobs      []*Job
	CreatedAt time.Time `pg:"default:now()"`
}
