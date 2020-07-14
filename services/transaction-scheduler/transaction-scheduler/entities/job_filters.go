package entities

type JobFilters struct {
	TxHashes  []string `validate:"omitempty,unique,dive,isHash"`
	ChainUUID string   `validate:"omitempty,uuid"`
	Status    string   `validate:"omitempty"`
}
