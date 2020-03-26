package models

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
