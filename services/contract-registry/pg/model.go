package pg

import (
	"github.com/ethereum/go-ethereum/common"
)

// RepositoryModel represent a space where contract tags are listed
type RepositoryModel struct {
	tableName struct{} `sql:"repositories"` //nolint:unused,structcheck

	// ID technical identifier
	ID int

	// Repository name
	Name string
}

// TagModel represent a Tag on a Repository pointing towards a Source code
type TagModel struct {
	tableName struct{} `sql:"tags"` //nolint:unused,structcheck

	// ID technical identifier
	ID int

	// Tag name
	Name         string
	RepositoryID int

	ArtifactID int
}

// ArtifactModel represent an artifact compiled from a source code
type ArtifactModel struct {
	tableName struct{} `sql:"artifacts"` //nolint:unused,structcheck

	// ID technical identifier
	ID int

	// Artifact data
	Abi              []byte
	Bytecode         []byte
	DeployedBytecode []byte
	// Codehash stored on the Ethereum account. Correspond to the hash of the deployedBytecode
	Codehash []byte
}

// CodehashModel represent the codehash of smart contract addresses
type CodehashModel struct {
	tableName struct{} `sql:"codehashes"` //nolint:unused,structcheck

	// ID technical identifier
	ID int

	// Artifact data
	ChainID  string
	Address  []byte
	Codehash []byte
}

// MethodModel represent the codehash of smart contract addresses
type MethodModel struct {
	tableName struct{} `sql:"methods"` //nolint:unused,structcheck

	// ID technical identifier
	ID int

	// Artifact data
	Codehash common.Hash
	Selector [4]byte

	ABI []byte
}

// EventModel represent the codehash of smart contract addresses
type EventModel struct {
	tableName struct{} `sql:"events"` //nolint:unused,structcheck

	// ID technical identifier
	ID int

	// Artifact data
	Codehash          common.Hash
	SigHash           common.Hash
	IndexedInputCount uint

	ABI []byte
}
