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
	UUID      string
	ChainUUID string
	URL       string
	Type      PrivateTxManagerType
	CreatedAt time.Time
}
