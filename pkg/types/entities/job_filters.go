package entities

import "time"

type JobFilters struct {
	TxHashes      []string  `validate:"omitempty,unique,dive,isHash"`
	ChainUUID     string    `validate:"omitempty,uuid"`
	Status        string    `validate:"omitempty,isJobStatus"`
	UpdatedAfter  time.Time `validate:"omitempty"`
	ParentJobUUID string    `validate:"omitempty"`
	OnlyParents   bool      `validate:"omitempty"`
}
