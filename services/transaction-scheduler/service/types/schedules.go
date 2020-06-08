package types

import (
	"time"
)

type CreateScheduleRequest struct{}

type ScheduleResponse struct {
	UUID      string         `json:"uuid"`
	Jobs      []*JobResponse `json:"jobs"`
	CreatedAt time.Time      `json:"createdAt"`
}
