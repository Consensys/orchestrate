package entities

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/ethereum"

	"github.com/ethereum/go-ethereum/common"
)

type Job struct {
	UUID         string
	ChainUUID    string
	ScheduleUUID string
	Type         string
	Labels       map[string]string
	Status       string
	Transaction  *ETHTransaction
	Receipt      *ethereum.Receipt
	CreatedAt    time.Time
}

type JobFilters struct {
	TxHashes []common.Hash
}
