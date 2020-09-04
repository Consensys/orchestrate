package entities

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/ethereum"
)

type Job struct {
	UUID         string
	NextJobUUID  string
	ChainUUID    string
	ScheduleUUID string
	Type         string
	Labels       map[string]string
	InternalData *InternalData
	Transaction  *ETHTransaction
	Receipt      *ethereum.Receipt
	Logs         []*Log
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// GetStatus Computes the status of a Job by checking its logs
func (job *Job) GetStatus() string {
	return getStatus(job.Logs)
}

func getStatus(logs []*Log) string {
	var status string
	var logCreatedAt *time.Time
	for idx := range logs {
		if logCreatedAt == nil || logs[idx].CreatedAt.After(*logCreatedAt) {
			status = logs[idx].Status
			logCreatedAt = &logs[idx].CreatedAt
		}
	}

	// Recursive function until we reach a valid status
	if status == utils.StatusWarning {
		return getStatus(logs[:len(logs)-1])
	}

	return status
}
