package entities

type AccountFilters struct {
	Aliases []string `validate:"omitempty,unique"`
}
