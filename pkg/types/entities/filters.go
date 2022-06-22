package entities

import "time"

type JobFilters struct {
	TxHashes      []string  `validate:"omitempty,unique,dive,isHash"`
	ChainUUID     string    `validate:"omitempty,uuid"`
	Status        JobStatus `validate:"omitempty,isJobStatus"`
	UpdatedAfter  time.Time `validate:"omitempty"`
	ParentJobUUID string    `validate:"omitempty"`
	OnlyParents   bool      `validate:"omitempty"`
	WithLogs      bool      `validate:"omitempty"`
}

type PaginationFilters struct {
	Limit int `validate:"omitempty"`
	Page  int `validate:"omitempty"`
}

type TransactionRequestFilters struct {
	IdempotencyKeys []string          `validate:"omitempty,unique"`
	Pagination      PaginationFilters `validate:"omitempty"`
}

type FaucetFilters struct {
	Names     []string `validate:"omitempty,unique"`
	ChainRule string   `validate:"omitempty"`
	TenantID  string   `validate:"omitempty"`
}

type AccountFilters struct {
	Aliases    []string          `validate:"omitempty,unique"`
	TenantID   string            `validate:"omitempty"`
	Pagination PaginationFilters `validate:"omitempty"`
}

type ChainFilters struct {
	Names    []string `validate:"omitempty,unique"`
	TenantID string   `validate:"omitempty"`
}
