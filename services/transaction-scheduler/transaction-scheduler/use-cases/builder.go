package usecases

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/jobs"
)

const Component = "use-cases"

type UseCases struct {
	CreateJobs jobs.CreateJobUseCase
}

func NewUseCases(dataAgents *store.DataAgents) *UseCases {
	return &UseCases{
		CreateJobs: jobs.NewCreateJob(dataAgents.TransactionJob),
	}
}
