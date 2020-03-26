package models

// TagModel represent a Tag on a Repository pointing towards a Source code
type TagModel struct {
	tableName struct{} `pg:"tags"` //nolint:unused,structcheck // reason

	// UUID technical identifier
	ID int

	// Tag name
	Name         string
	RepositoryID int

	ArtifactID int
}
