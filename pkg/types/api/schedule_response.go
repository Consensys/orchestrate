package api

import "time"

type ScheduleResponse struct {
	UUID      string         `json:"uuid" example:"b4374e6f-b28a-4bad-b4fe-bda36eaf849c"`
	TenantID  string         `json:"tenantID" example:"tenant_id"`
	Jobs      []*JobResponse `json:"jobs"`
	CreatedAt time.Time      `json:"createdAt" example:"2020-07-09T12:35:42.115395Z"`
}
