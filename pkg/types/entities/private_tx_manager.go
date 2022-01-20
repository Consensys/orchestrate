package entities

import "time"

type PrivateTxType string
type PrivateTxManagerType string

const (
	PrivateTxTypeRestricted PrivateTxType = "restricted"

	TesseraChainType PrivateTxManagerType = "Tessera"
	EEAChainType     PrivateTxManagerType = "EEA"

	// Minimum gas is calculated by the size of the enclaveKey
	TesseraGasLimit = 60000
)

type PrivateTxManager struct {
	UUID      string               // UUID of the private transaction manager.
	ChainUUID string               // UUID of the registered chain.
	URL       string               // Transaction manager endpoint.
	Type      PrivateTxManagerType // Currently supports `Tessera` and `EEA`.
	CreatedAt time.Time            // Date and time that the private transaction manager was registered with the chain.
}
