package entities

import "time"

type Log struct {
	Status    JobStatus `json:"status" example:"MINED"`
	Message   string    `json:"message,omitempty" example:"Log message"`
	CreatedAt time.Time `json:"at" example:"2020-07-09T12:35:42.115395Z"`
}
