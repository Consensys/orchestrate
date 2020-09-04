package listeners

import (
	"context"
	"sync"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-sentry/tx-sentry/use-cases"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"

	"github.com/cenkalti/backoff/v4"
	backoffjob "github.com/containous/traefik/v2/pkg/job"
	log "github.com/sirupsen/logrus"
)

//go:generate mockgen -source=session_manager.go -destination=mocks/session_manager.go -package=mocks

const sessionManagerComponent = "service.session-manager"

type SessionManager interface {
	Start(ctx context.Context, job *entities.Job)
}

// sessionManager is a manager of job sessions
type sessionManager struct {
	mutex                 *sync.RWMutex
	sessions              map[string]*entities.Job
	createChildJobUseCase usecases.CreateChildJobUseCase
}

// NewSessionManager creates a new SessionManager
func NewSessionManager(createChildJobUseCase usecases.CreateChildJobUseCase) SessionManager {
	return &sessionManager{
		mutex:                 &sync.RWMutex{},
		sessions:              make(map[string]*entities.Job),
		createChildJobUseCase: createChildJobUseCase,
	}
}

func (manager *sessionManager) Start(ctx context.Context, job *entities.Job) {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID)

	if manager.getSession(job.UUID) != nil {
		logger.Debug("job session already exists, skipping session creation")
		return
	}

	if job.InternalData.ParentJobUUID != "" {
		logger.Debug("job session is not a parent job, skipping session creation")
		return
	}

	manager.addSession(job)

	go func() {
		backff := backoff.WithContext(backoffjob.NewBackOff(backoff.NewExponentialBackOff()), ctx)
		err := backoff.RetryNotify(
			func() error { return manager.runSession(ctx, job) },
			backff,
			func(err error, duration time.Duration) {
				logger.WithError(err).Warnf("error in job listening session, restarting in %v...", duration)
			},
		)
		// At the moment, this should never happen as the session should either:
		// - fail and retry forever
		// - gracefully stop without error
		// A different failure strategy could be implemented to not retry at infinity in case of recurrent failure
		if err != nil {
			logger.WithError(err).Error("job listening session unexpectedly stopped")
		}

		manager.removeSession(job.UUID)
	}()
}

func (manager *sessionManager) addSession(job *entities.Job) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	manager.sessions[job.UUID] = job
}

func (manager *sessionManager) getSession(jobUUID string) *entities.Job {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()
	return manager.sessions[jobUUID]
}

func (manager *sessionManager) runSession(ctx context.Context, job *entities.Job) error {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID)
	logger.Info("starting job session")

	ticker := time.NewTicker(job.InternalData.RetryInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			childJobUUID, err := manager.createChildJobUseCase.Execute(ctx, job)
			if err != nil {
				return errors.FromError(err).ExtendComponent(sessionManagerComponent)
			}

			// If no child created but no error, we exit the session
			if childJobUUID == "" {
				return nil
			}
		case <-ctx.Done():
			logger.WithField("reason", ctx.Err().Error()).Info("session gracefully stopped")
			return nil
		}
	}
}

func (manager *sessionManager) removeSession(jobUUID string) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	delete(manager.sessions, jobUUID)
}
