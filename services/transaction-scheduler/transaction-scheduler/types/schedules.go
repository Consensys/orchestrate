package types

import "time"

type ScheduleRequest struct {
	ChainID string `json:"chainID" validate:"required,uuid4"`
}

type ScheduleResponse struct {
	UUID      string         `json:"uuid" validate:"required,uuid4"`
	ChainID   string         `json:"chainID" validate:"required,uuid4"`
	Jobs      []*JobResponse `json:"jobs,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
}
