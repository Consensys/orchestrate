package listeners

import (
	"context"
	"sync"
	"time"

	"github.com/consensys/orchestrate/pkg/errors"
	orchestrateclient "github.com/consensys/orchestrate/pkg/sdk/client"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	types "github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/formatters"
	usecases "github.com/consensys/orchestrate/services/tx-sentry/tx-sentry/use-cases"

	"github.com/consensys/orchestrate/pkg/types/entities"

	"github.com/cenkalti/backoff/v4"
	pkgbackoff "github.com/consensys/orchestrate/pkg/backoff"
)

//go:generate mockgen -source=session_manager.go -destination=mocks/session_manager.go -package=mocks

const sessionManagerComponent = "tx-sentry.service.session-manager"

type SessionManager interface {
	Start(ctx context.Context, job *entities.Job)
}

// sessionManager is a manager of job sessions
type sessionManager struct {
	mutex                  *sync.RWMutex
	sessions               map[string]bool
	retrySessionJobUseCase usecases.RetrySessionJobUseCase
	client                 orchestrateclient.OrchestrateClient
	logger                 *log.Logger
}

type sessionData struct {
	parentJob        *entities.Job
	nChildren        int
	retries          int
	lastChildJobUUID string
}

// NewSessionManager creates a new SessionManager
func NewSessionManager(client orchestrateclient.OrchestrateClient, retrySessionJobUseCase usecases.RetrySessionJobUseCase) SessionManager {
	return &sessionManager{
		mutex:                  &sync.RWMutex{},
		sessions:               make(map[string]bool),
		retrySessionJobUseCase: retrySessionJobUseCase,
		client:                 client,
		logger:                 log.NewLogger().SetComponent(sessionManagerComponent),
	}
}

func (manager *sessionManager) Start(ctx context.Context, job *entities.Job) {
	logger := manager.logger.WithContext(ctx).WithField("job", job.UUID).WithField("tenant", job.TenantID)
	ctx = log.With(ctx, logger)

	if manager.hasSession(job.UUID) {
		logger.Trace("job session already exists, skipping session creation")
		return
	}

	if job.InternalData.RetryInterval == 0 {
		logger.Trace("job session does not have any retry strategy")
		return
	}

	if job.InternalData.HasBeenRetried {
		logger.Warn("job session been already retried")
		return
	}

	ses, err := manager.retrieveJobSessionData(ctx, job)
	if err != nil {
		logger.WithError(err).Error("job listening session failed to start")
		return
	}

	if ses.retries >= types.SentryMaxRetries {
		logger.Warn("job already reached max retries")
		return
	}

	manager.addSession(job.UUID)

	go func() {
		err := backoff.RetryNotify(
			func() error {
				err := manager.runSession(ctx, ses)
				return err
			},
			pkgbackoff.IncrementalBackOff(time.Second, 5*time.Second, time.Minute),
			func(err error, d time.Duration) {
				logger.WithError(err).Warnf("error in job retry session, restarting in %v...", d)
			},
		)

		if err != nil {
			logger.WithError(err).Error("job listening session unexpectedly stopped")
		}

		annotations := formatters.FormatInternalDataToAnnotations(job.InternalData)
		annotations.HasBeenRetried = true
		_, err = manager.client.UpdateJob(ctx, job.UUID, &types.UpdateJobRequest{
			Annotations: &annotations,
		})

		if err != nil {
			logger.WithError(err).Error("failed to update job labels")
		}

		logger.Debug("job session was completed")
		manager.removeSession(job.UUID)
	}()
}

func (manager *sessionManager) addSession(jobUUID string) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	manager.sessions[jobUUID] = true
}

func (manager *sessionManager) hasSession(jobUUID string) bool {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()
	_, ok := manager.sessions[jobUUID]
	return ok
}

func (manager *sessionManager) runSession(ctx context.Context, ses *sessionData) error {
	logger := log.FromContext(ctx)
	logger.Info("job session started")

	ticker := time.NewTicker(ses.parentJob.InternalData.RetryInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			childJobUUID, err := manager.retrySessionJobUseCase.Execute(ctx, ses.parentJob.UUID, ses.lastChildJobUUID, ses.nChildren)
			if err != nil {
				return errors.FromError(err).ExtendComponent(sessionManagerComponent)
			}

			ses.retries++
			if ses.retries >= types.SentryMaxRetries {
				return nil
			}

			// If no child created but no error, we exit the session gracefully
			if childJobUUID == "" {
				return nil
			}

			if childJobUUID != ses.lastChildJobUUID {
				ses.nChildren++
				ses.lastChildJobUUID = childJobUUID
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

func (manager *sessionManager) retrieveJobSessionData(ctx context.Context, job *entities.Job) (*sessionData, error) {
	jobs, err := manager.client.SearchJob(ctx, &entities.JobFilters{
		ChainUUID:     job.ChainUUID,
		ParentJobUUID: job.UUID,
		WithLogs:      true,
	})

	if err != nil {
		return nil, err
	}

	nChildren := len(jobs) - 1
	lastJobRetry := jobs[len(jobs)-1]

	// we count the number of resending of last job as retries
	nRetries := nChildren
	for _, lg := range lastJobRetry.Logs {
		if lg.Status == entities.StatusResending {
			nRetries++
		}
	}

	return &sessionData{
		parentJob:        job,
		nChildren:        nChildren,
		retries:          nRetries,
		lastChildJobUUID: jobs[nChildren].UUID,
	}, nil
}
