package models

type ArtifactModel struct {
	tableName struct{} `pg:"artifacts"` // nolint:unused,structcheck // reason

	// UUID technical identifier
	ID int

	// Artifact data
	ABI              string `pg:"alias:abi"`
	Bytecode         string
	DeployedBytecode string
	// Codehash stored on the Ethereum account. Correspond to the hash of the deployedBytecode
	Codehash string
}

type CodehashModel struct {
	tableName struct{} `pg:"codehashes"` // nolint:unused,structcheck // reason

	// UUID technical identifier
	ID int

	// Artifact data
	ChainID  string `pg:"alias:chain_id"`
	Address  string
	Codehash string
}

type EventModel struct {
	tableName struct{} `pg:"events"` // nolint:unused,structcheck // reason

	// UUID technical identifier
	ID int

	// Artifact data
	Codehash          string
	SigHash           string
	IndexedInputCount uint `pg:",use_zero"`

	ABI string
}

type MethodModel struct {
	tableName struct{} `pg:"methods"` // nolint:unused,structcheck // reason

	// UUID technical identifier
	ID int

	// Artifact data
	Codehash string
	Selector [4]byte

	ABI string
}

type RepositoryModel struct {
	tableName struct{} `pg:"repositories"` // nolint:unused,structcheck // reason

	// UUID technical identifier
	ID int

	// Repository name
	Name string
}

type TagModel struct {
	tableName struct{} `pg:"tags"` // nolint:unused,structcheck // reason

	// UUID technical identifier
	ID int

	// Tag name
	Name         string
	RepositoryID int

	ArtifactID int
}
