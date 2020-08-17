package sessions

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"

	entities2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/tx-sentry/entities"

	log "github.com/sirupsen/logrus"
)

//go:generate mockgen -source=create_session.go -destination=mocks/create_session.go -package=mocks

// const createSessionComponent = "use-cases.create-session"

type CreateSessionUseCase interface {
	Execute(ctx context.Context, job *entities.Job) *entities2.JobSession
}

// createSessionUseCase is a use case to create a new transaction job
type createSessionUseCase struct {
}

// NewCreateSessionUseCase creates a new CreateSessionUseCase
func NewCreateSessionUseCase() CreateSessionUseCase {
	return &createSessionUseCase{}
}

// Execute validates and creates a new job session
func (uc *createSessionUseCase) Execute(ctx context.Context, job *entities.Job) *entities2.JobSession {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID)
	logger.Debug("creating new job session")

	_, cancel := context.WithCancel(ctx)
	session := &entities2.JobSession{
		Job:    job,
		Cancel: cancel,
	}

	logger.Info("new job session successfully created")
	return session
}
