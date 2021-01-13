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

type TransactionRequestFilters struct {
	IdempotencyKeys []string `validate:"omitempty,unique"`
}

type FaucetFilters struct {
	Names     []string `validate:"omitempty,unique"`
	ChainRule string   `validate:"omitempty"`
}

type AccountFilters struct {
	Aliases []string `validate:"omitempty,unique"`
}

type ChainFilters struct {
	Names []string `validate:"omitempty,unique"`
}
