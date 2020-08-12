package listeners

import (
	"context"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/jobs"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
)

//go:generate mockgen -source=jobs_listener.go -destination=mocks/jobs_listener.go -package=mocks

const jobsListenerComponent = "service.jobs-listener"

type JobsListener interface {
	Listen(ctx context.Context) chan error
}

type jobsListener struct {
	refreshInterval   time.Duration
	sessionManager    SessionManager
	searchJobsUseCase jobs.SearchJobsUseCase
}

func NewJobsListener(refreshInterval time.Duration, sessionManager SessionManager, searchJobsUseCase jobs.SearchJobsUseCase) JobsListener {
	return &jobsListener{
		refreshInterval:   refreshInterval,
		sessionManager:    sessionManager,
		searchJobsUseCase: searchJobsUseCase,
	}
}

func (listener *jobsListener) Listen(ctx context.Context) chan error {
	log.WithContext(ctx).Debug("starting tx-sentry jobs listener")

	cerr := make(chan error)

	go func() {
		ticker := time.NewTicker(listener.refreshInterval)
		defer ticker.Stop()

		// Initial job push of ALL pending jobs
		// Then we fetch only new pending jobs that are older than the given refresh interval
		jobsFilter := &entities.JobFilters{
			Status: utils.StatusPending,
		}
		for {
			select {
			case t := <-ticker.C:
				err := listener.retrieveJobs(ctx, jobsFilter)
				if err != nil {
					cerr <- errors.FromError(err).ExtendComponent(jobsListenerComponent)
					return
				}
				jobsFilter.UpdatedAfter = t.Add(-listener.refreshInterval)
			case <-ctx.Done():
				return
			}
		}
	}()

	return cerr
}

func (listener *jobsListener) retrieveJobs(ctx context.Context, jobsFilter *entities.JobFilters) error {
	jobsArr, err := listener.searchJobsUseCase.Execute(ctx, jobsFilter, []string{multitenancy.Wildcard})
	if err != nil {
		errMessage := "failed to get latest pending jobs"
		log.WithContext(ctx).WithError(err).Error(errMessage)
		return errors.InternalError(errMessage)
	}

	for _, job := range jobsArr {
		err := listener.sessionManager.AddSession(ctx, job)
		if err != nil {
			errMessage := "failed to add job session"
			log.WithContext(ctx).WithError(err).Error(errMessage)
			return errors.InternalError(errMessage)
		}
	}

	return nil
}
