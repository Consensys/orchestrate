package pg

// RepositoryModel represent a space where contract tags are listed
type RepositoryModel struct {
	tableName struct{} `pg:"repositories"` //nolint:unused,structcheck // reason

	// UUID technical identifier
	ID int

	// Repository name
	Name string
}

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

// ArtifactModel represent an artifact compiled from a source code
type ArtifactModel struct {
	tableName struct{} `pg:"artifacts"` //nolint:unused,structcheck // reason

	// UUID technical identifier
	ID int

	// Artifact data
	Abi              string
	Bytecode         string
	DeployedBytecode string
	// Codehash stored on the Ethereum account. Correspond to the hash of the deployedBytecode
	Codehash string
}

// CodehashModel represent the codehash of smart contract addresses
type CodehashModel struct {
	tableName struct{} `pg:"codehashes"` //nolint:unused,structcheck // reason

	// UUID technical identifier
	ID int

	// Artifact data
	ChainID  string
	Address  string
	Codehash string
}

// MethodModel represent the codehash of smart contract addresses
type MethodModel struct {
	tableName struct{} `pg:"methods"` //nolint:unused,structcheck // reason

	// UUID technical identifier
	ID int

	// Artifact data
	Codehash string
	Selector [4]byte

	ABI string
}

// EventModel represent the codehash of smart contract addresses
type EventModel struct {
	tableName struct{} `pg:"events"` //nolint:unused,structcheck // reason

	// UUID technical identifier
	ID int

	// Artifact data
	Codehash          string
	SigHash           string
	IndexedInputCount uint `pg:",use_zero"`

	ABI string
}
