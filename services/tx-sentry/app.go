package txsentry

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/consensys/orchestrate/pkg/errors"
	orchestrateclient "github.com/consensys/orchestrate/pkg/sdk/client"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/services/tx-sentry/service/listeners"
	"github.com/consensys/orchestrate/services/tx-sentry/service/parsers"
	usecases "github.com/consensys/orchestrate/services/tx-sentry/tx-sentry/use-cases"
	backoffjob "github.com/traefik/traefik/v2/pkg/job"
)

const txSentryComponent = "application.tx-sentry"

type TxSentry struct {
	client         orchestrateclient.OrchestrateClient
	sessionManager listeners.SessionManager
	config         *Config
	logger         *log.Logger
}

func NewTxSentry(client orchestrateclient.OrchestrateClient, config *Config) *TxSentry {
	createChildJobUC := usecases.NewRetrySessionJobUseCase(client)
	return &TxSentry{
		client:         client,
		sessionManager: listeners.NewSessionManager(client, createChildJobUC),
		config:         config,
		logger:         log.NewLogger().SetComponent(txSentryComponent),
	}
}

func (sentry *TxSentry) Run(ctx context.Context) error {
	ctx = log.With(ctx, sentry.logger)

	backff := backoff.WithContext(backoffjob.NewBackOff(backoff.NewExponentialBackOff()), ctx)
	err := backoff.RetryNotify(
		func() error { return sentry.listen(ctx) },
		backff,
		func(err error, duration time.Duration) {
			sentry.logger.WithError(err).Warnf("error in job listening, restarting in %v...", duration)
		},
	)

	if err != nil && err != context.Canceled {
		sentry.logger.WithError(err).Errorf("sentry stopped after catching an error")
	}

	sentry.logger.Info("transaction sentry stopped without error")

	return nil
}

func (sentry *TxSentry) Close() error {
	return nil
}

func (sentry *TxSentry) listen(ctx context.Context) error {
	sentry.logger.Info("jobs listener started")

	// Initial job creation fetching all pending jobs
	jobFilters := &entities.JobFilters{
		Status:      entities.StatusPending,
		OnlyParents: true,
	}

	err := sentry.createSessions(ctx, jobFilters)
	if err != nil {
		return errors.FromError(err).ExtendComponent(txSentryComponent)
	}

	ticker := time.NewTicker(sentry.config.RefreshInterval)
	defer ticker.Stop()
	for {
		select {
		case t := <-ticker.C:
			lastCheckedAt := t.Add(-sentry.config.RefreshInterval)
			sentry.logger.WithField("updated_after", lastCheckedAt.Format("2006-01-02 15:04:05")).
				Debug("fetching new pending jobs")

			jobFilters.UpdatedAfter = lastCheckedAt
			err := sentry.createSessions(ctx, jobFilters)
			if err != nil {
				return errors.FromError(err).ExtendComponent(txSentryComponent)
			}
		case <-ctx.Done():
			sentry.logger.WithField("reason", ctx.Err().Error()).Info("gracefully stopping transaction sentry...")
			return nil
		}
	}
}

func (sentry *TxSentry) createSessions(ctx context.Context, filters *entities.JobFilters) error {
	// We get all the pending jobs updated_after the last tick
	jobResponses, err := sentry.client.SearchJob(ctx, filters)
	if err != nil {
		sentry.logger.WithError(err).Error("failed to fetch pending jobs")
		return err
	}

	for _, jobResponse := range jobResponses {
		jctx := multitenancy.WithTenantID(ctx, jobResponse.TenantID)
		sentry.sessionManager.Start(jctx, parsers.JobResponseToEntity(jobResponse))
	}

	return nil
}
