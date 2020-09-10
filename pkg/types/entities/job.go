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
		// Ignore resending and warning statuses
		if logs[idx].Status == utils.StatusResending || logs[idx].Status == utils.StatusWarning {
			continue
		}
		// Ignore fail statuses if they come after a resending
		if logs[idx].Status == utils.StatusFailed && idx > 1 && logs[idx-1].Status == utils.StatusResending {
			continue
		}

		if logCreatedAt == nil || logs[idx].CreatedAt.After(*logCreatedAt) {
			status = logs[idx].Status
			logCreatedAt = &logs[idx].CreatedAt
		}
	}

	return status
}
