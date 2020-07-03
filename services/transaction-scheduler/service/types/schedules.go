package types

import (
	"time"
)

type CreateScheduleRequest struct{}

type ScheduleResponse struct {
	UUID      string         `json:"uuid"`
	TenantID  string         `json:"tenantID"`
	Jobs      []*JobResponse `json:"jobs"`
	CreatedAt time.Time      `json:"createdAt"`
}
