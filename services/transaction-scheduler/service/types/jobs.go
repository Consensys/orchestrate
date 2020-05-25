package types

import (
	"time"
)

type CreateJobRequest struct {
	ScheduleUUID string            `json:"scheduleUUID" validate:"required"`
	Type         string            `json:"type" validate:"required"` //  @TODO validate Type is valid
	Labels       map[string]string `json:"labels,omitempty"`
	Transaction  ETHTransaction    `json:"transaction" validate:"required"`
}

type UpdateJobRequest struct {
	Labels      map[string]string `json:"labels,omitempty"`
	Transaction ETHTransaction    `json:"transaction"`
}

type JobResponse struct {
	UUID        string         `json:"uuid" validate:"required,uuid4"`
	Transaction ETHTransaction `json:"transaction" validate:"required"`
	Status      string         `json:"status" validate:"required"`
	CreatedAt   time.Time      `json:"createdAt"`
}

type ETHTransaction struct {
	Hash           string   `json:"hash,omitempty"`
	From           string   `json:"from,omitempty"`
	To             string   `json:"to,omitempty"`
	Nonce          string   `json:"nonce,omitempty"`
	Value          string   `json:"value,omitempty"`
	GasPrice       string   `json:"gasPrice,omitempty"`
	GasLimit       string   `json:"gasLimit,omitempty"`
	Data           string   `json:"data,omitempty"`
	Raw            string   `json:"raw,omitempty"`
	PrivateFrom    string   `json:"privateFrom,omitempty"`
	PrivateFor     []string `json:"privateFor,omitempty"`
	PrivacyGroupID string   `json:"privacyGroupID,omitempty"`
}
