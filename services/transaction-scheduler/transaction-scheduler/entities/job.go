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
	Transaction  *ETHTransaction
	CreatedAt    time.Time
}
