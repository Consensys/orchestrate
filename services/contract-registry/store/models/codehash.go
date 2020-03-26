package models

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
