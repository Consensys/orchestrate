package entities

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/ethereum"
)

type Job struct {
	UUID         string
	NextJobUUID  string
	ChainUUID    string
	ScheduleUUID string
	TenantID     string
	Type         string
	Status       string
	Labels       map[string]string
	InternalData *InternalData
	Transaction  *ETHTransaction
	Receipt      *ethereum.Receipt
	Logs         []*Log
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
