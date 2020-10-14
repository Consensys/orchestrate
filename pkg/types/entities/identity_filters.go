package entities

type IdentityFilters struct {
	Aliases []string `validate:"omitempty,unique"`
}
