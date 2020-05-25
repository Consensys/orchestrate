package types

import (
	"time"
)

type CreateScheduleRequest struct {
	ChainUUID string `json:"chainUUID" validate:"required,uuid4"`
}

type ScheduleResponse struct {
	UUID      string         `json:"uuid" validate:"required,uuid4"`
	ChainUUID string         `json:"chainUUID" validate:"required,uuid4"`
	Jobs      []*JobResponse `json:"jobs,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
}
