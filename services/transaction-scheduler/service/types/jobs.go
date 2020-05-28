package types

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
)

type CreateJobRequest struct {
	ScheduleUUID string                   `json:"scheduleUUID" validate:"required"`
	Type         string                   `json:"type" validate:"required"` //  @TODO validate Type is valid
	Labels       map[string]string        `json:"labels,omitempty"`
	Transaction  *entities.ETHTransaction `json:"transaction" validate:"required"`
}

type UpdateJobRequest struct {
	Labels      map[string]string        `json:"labels,omitempty"`
	Transaction *entities.ETHTransaction `json:"transaction"`
}

type JobResponse struct {
	UUID        string                   `json:"uuid" validate:"required,uuid4"`
	Transaction *entities.ETHTransaction `json:"transaction" validate:"required"`
	Status      string                   `json:"status" validate:"required"`
	CreatedAt   time.Time                `json:"createdAt"`
}
