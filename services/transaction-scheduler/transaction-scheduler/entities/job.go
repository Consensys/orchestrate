package entities

import (
	"time"
)

type Job struct {
	UUID         string
	ScheduleUUID string
	Type         string
	Labels       map[string]string
	Status       string
	Transaction  *Transaction
	CreatedAt    time.Time
}

// type JobFilter struct {
// 	UUID     string
// 	Filters  map[string]string
// 	TenantID string
// }
