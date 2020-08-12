package listeners

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/tx-sentry/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/tx-sentry/use-cases/sessions"
)

//go:generate mockgen -source=session_manager.go -destination=mocks/session_manager.go -package=mocks

// const sessionManagerComponent = "service.session-manager"

type SessionManager interface {
	AddSession(ctx context.Context, job *types.Job) error
}

// sessionManager is a manager of job sessions
type sessionManager struct {
	sessionsMap          map[string]*entities.JobSession
	createSessionUseCase sessions.CreateSessionUseCase
}

// NewSessionManager creates a new SessionManager
func NewSessionManager(createSessionUseCase sessions.CreateSessionUseCase) SessionManager {
	return &sessionManager{
		sessionsMap:          make(map[string]*entities.JobSession),
		createSessionUseCase: createSessionUseCase,
	}
}

func (manager *sessionManager) AddSession(ctx context.Context, job *types.Job) error {
	if manager.sessionsMap[job.UUID] != nil {
		log.WithContext(ctx).WithField("job_uuid", job.UUID).Warn("job session already exists")
		return nil
	}

	manager.sessionsMap[job.UUID] = manager.createSessionUseCase.Execute(ctx, job)
	return nil
}
