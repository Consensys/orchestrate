package entities

import (
	"time"
)

type Schedule struct {
	UUID      string
	TenantID  string
	OwnerID   string
	Jobs      []*Job
	CreatedAt time.Time
}
