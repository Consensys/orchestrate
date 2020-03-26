package models

// RepositoryModel represents a space where contract tags are listed
type RepositoryModel struct {
	tableName struct{} `pg:"repositories"` //nolint:unused,structcheck // reason

	// UUID technical identifier
	ID int

	// Repository name
	Name string
}
