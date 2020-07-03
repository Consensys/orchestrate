package types

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
)

type CreateJobRequest struct {
	ScheduleUUID string                `json:"scheduleUUID" validate:"required,uuid4"`
	ChainUUID    string                `json:"chainUUID" validate:"required,uuid4"`
	Type         string                `json:"type" validate:"required"` //  @TODO validate Type is valid
	Labels       map[string]string     `json:"labels,omitempty"`
	Annotations  *types.Annotations    `json:"annotations,omitempty"`
	Transaction  *types.ETHTransaction `json:"transaction" validate:"required"`
}

type UpdateJobRequest struct {
	Labels      map[string]string     `json:"labels,omitempty"`
	Annotations *types.Annotations    `json:"annotations,omitempty"`
	Transaction *types.ETHTransaction `json:"transaction,omitempty"`
	Status      string                `json:"status,omitempty"`
	Message     string                `json:"message,omitempty"`
}

type JobResponse struct {
	UUID         string                `json:"uuid"`
	ChainUUID    string                `json:"chainUUID"`
	ScheduleUUID string                `json:"scheduleUUID"`
	Transaction  *types.ETHTransaction `json:"transaction"`
	Logs         []*types.Log          `json:"logs"`
	Labels       map[string]string     `json:"labels,omitempty"`
	Annotations  *types.Annotations    `json:"annotations,omitempty"`
	Status       string                `json:"status"`
	Type         string                `json:"type"`
	CreatedAt    time.Time             `json:"createdAt"`
}
