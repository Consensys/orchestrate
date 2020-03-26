package models

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
