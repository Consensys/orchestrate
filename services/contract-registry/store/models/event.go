package models

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
