package txsentry

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
	backoffjob "github.com/containous/traefik/v2/pkg/job"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
	orchestrateclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sentry/service/listeners"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sentry/service/parsers"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sentry/tx-sentry/use-cases"
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
	logger := sentry.logger.WithContext(ctx)

	backff := backoff.WithContext(backoffjob.NewBackOff(backoff.NewExponentialBackOff()), ctx)
	err := backoff.RetryNotify(
		func() error { return sentry.listen(ctx) },
		backff,
		func(err error, duration time.Duration) {
			logger.WithError(err).Warnf("error in job listening, restarting in %v...", duration)
		},
	)

	if err != nil && err != context.Canceled {
		logger.WithError(err).Errorf("sentry stopped after catching an error")
	}

	logger.Infof("transaction sentry stopped without error")

	return nil
}

func (sentry *TxSentry) Close() error {
	return nil
}

func (sentry *TxSentry) listen(ctx context.Context) error {
	logger := sentry.logger.WithContext(ctx)
	logger.Info("jobs listener started")

	// Initial job creation fetching all pending jobs
	jobFilters := &entities.JobFilters{
		Status:      utils.StatusPending,
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
			logger.WithField("updated_after", lastCheckedAt.Format("2006-01-02 15:04:05")).
				Debug("fetching new pending jobs")

			jobFilters.UpdatedAfter = lastCheckedAt
			err := sentry.createSessions(ctx, jobFilters)
			if err != nil {
				return errors.FromError(err).ExtendComponent(txSentryComponent)
			}
		case <-ctx.Done():
			logger.WithField("reason", ctx.Err().Error()).Info("gracefully stopping transaction sentry...")
			return nil
		}
	}
}

func (sentry *TxSentry) createSessions(ctx context.Context, filters *entities.JobFilters) error {
	// We get all the pending jobs updated_after the last tick
	jobResponses, err := sentry.client.SearchJob(ctx, filters)
	if err != nil {
		sentry.logger.WithContext(ctx).WithError(err).Error("failed to fetch pending jobs")
		return err
	}

	for _, jobResponse := range jobResponses {
		jctx := multitenancy.WithTenantID(ctx, jobResponse.TenantID)
		sentry.sessionManager.Start(jctx, parsers.JobResponseToEntity(jobResponse))
	}

	return nil
}
