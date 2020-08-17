package txschedulertypes

import "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"

type CreateJobRequest struct {
	ScheduleUUID string                  `json:"scheduleUUID" validate:"required,uuid4" example:"b4374e6f-b28a-4bad-b4fe-bda36eaf849c"`
	ChainUUID    string                  `json:"chainUUID" validate:"required,uuid4" example:"b4374e6f-b28a-4bad-b4fe-bda36eaf849c"`
	Type         string                  `json:"type" validate:"required,isJobType" example:"eth://ethereum/transaction"`
	Labels       map[string]string       `json:"labels,omitempty"`
	Annotations  Annotations             `json:"annotations,omitempty"`
	Transaction  entities.ETHTransaction `json:"transaction" validate:"required"`
}

type UpdateJobRequest struct {
	Labels      map[string]string        `json:"labels,omitempty"`
	Annotations *Annotations             `json:"annotations,omitempty"`
	Transaction *entities.ETHTransaction `json:"transaction,omitempty"`
	Status      string                   `json:"status,omitempty" validate:"isJobStatus" example:"MINED"`
	Message     string                   `json:"message,omitempty" example:"Update message"`
}
