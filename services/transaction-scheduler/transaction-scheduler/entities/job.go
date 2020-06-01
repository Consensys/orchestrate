package entities

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Job struct {
	UUID         string
	ScheduleUUID string
	Type         string
	Labels       map[string]string
	Status       string
	Transaction  *ETHTransaction
	CreatedAt    time.Time
}

type JobFilters struct {
	TxHashes []common.Hash
}
