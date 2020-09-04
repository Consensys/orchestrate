package txsentry

import (
	"context"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-sentry/service/parsers"

	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-sentry/tx-sentry/use-cases"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"

	"github.com/cenkalti/backoff/v4"
	backoffjob "github.com/containous/traefik/v2/pkg/job"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-sentry/service/listeners"
)

const txSentryComponent = "tx-sentry"

type txsentryDaemon struct {
	cancel            context.CancelFunc
	txSchedulerClient client.TransactionSchedulerClient
	sessionManager    listeners.SessionManager
	config            *Config
	done              chan struct{}
}

func NewTxSentry(txSchedulerClient client.TransactionSchedulerClient, config *Config) app.Daemon {
	// Create business layer
	createChildJobUC := usecases.NewCreateChildJobUseCase(txSchedulerClient)
	return &txsentryDaemon{
		txSchedulerClient: txSchedulerClient,
		sessionManager:    listeners.NewSessionManager(createChildJobUC),
		config:            config,
	}
}

func (daemon *txsentryDaemon) Start(ctx context.Context) {
	logger := log.WithContext(ctx)
	logger.Debug("starting transaction sentry")

	daemon.done = make(chan struct{})
	ctx, daemon.cancel = context.WithCancel(ctx)
	go func() {
		defer close(daemon.done)

		backff := backoff.WithContext(backoffjob.NewBackOff(backoff.NewExponentialBackOff()), ctx)
		err := backoff.RetryNotify(
			func() error { return daemon.listen(ctx) },
			backff,
			func(err error, duration time.Duration) {
				logger.WithError(err).Warnf("error in job listening, restarting in %v...", duration)
			},
		)
		// This should never happen as the sentry will either retry at infinity or return without error
		if err != nil {
			logger.WithError(err).Error("transaction sentry unexpectedly stopped with an error")
			return
		}

		log.WithContext(ctx).Info("transaction sentry gracefully stopped")
	}()
}

func (daemon *txsentryDaemon) Stop(ctx context.Context) {
	log.WithContext(ctx).Debug("waiting for transaction sentry to stop...")
	daemon.cancel()
	<-daemon.done
}

func (daemon *txsentryDaemon) listen(ctx context.Context) error {
	logger := log.WithContext(ctx)
	logger.Info("starting tx-sentry jobs listener")

	// Initial job creation fetching all pending jobs
	jobFilters := &entities.JobFilters{
		Status:      utils.StatusPending,
		OnlyParents: true,
	}
	err := daemon.createSessions(ctx, jobFilters)
	if err != nil {
		return errors.FromError(err).ExtendComponent(txSentryComponent)
	}

	ticker := time.NewTicker(daemon.config.RefreshInterval)
	defer ticker.Stop()
	for {
		select {
		case t := <-ticker.C:
			lastCheckedAt := t.Add(-daemon.config.RefreshInterval)
			logger.WithField("updated_after", lastCheckedAt).Debug("fetching new pending jobs")

			jobFilters.UpdatedAfter = lastCheckedAt
			err := daemon.createSessions(ctx, jobFilters)
			if err != nil {
				return errors.FromError(err).ExtendComponent(txSentryComponent)
			}
		case <-ctx.Done():
			logger.WithField("reason", ctx.Err().Error()).Info("gracefully stopping transaction sentry...")
			return nil
		}
	}
}

func (daemon *txsentryDaemon) createSessions(ctx context.Context, filters *entities.JobFilters) error {
	// We get all the pending jobs updated_after the last tick
	jobResponses, err := daemon.txSchedulerClient.SearchJob(ctx, filters)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to fetch pending jobs")
		return err
	}

	for _, jobResponse := range jobResponses {
		daemon.sessionManager.Start(ctx, parsers.JobResponseToEntity(jobResponse))
	}

	return nil
}
