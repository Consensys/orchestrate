package entities

type TransactionFilters struct {
	IdempotencyKeys []string `validate:"omitempty,unique"`
}
