package transactionscheduler

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
	backoffjob "github.com/containous/traefik/v2/pkg/job"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/tx-sentry/use-cases/sessions"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/listeners"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/jobs"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/multi"
)

type txsentryDaemon struct {
	jobListener listeners.JobsListener
	ctx         context.Context
	cancel      context.CancelFunc
	isReady     bool
}

func NewTxSentryDaemon(pgmngr postgres.Manager, config *Config) (app.Daemon, error) {
	// Create Data agents
	db, err := multi.Build(context.Background(), config.Store, pgmngr)
	if err != nil {
		return nil, err
	}

	searchJobsUC := jobs.NewSearchJobsUseCase(db)
	createSessionUC := sessions.NewCreateSessionUseCase()

	sessionManager := listeners.NewSessionManager(createSessionUC)
	jobListener := listeners.NewJobsListener(config.Sentry.RefreshInterval, sessionManager, searchJobsUC)

	return &txsentryDaemon{
		jobListener: jobListener,
	}, nil
}

func (daemon *txsentryDaemon) IsReady() bool {
	return daemon.isReady
}

func (daemon *txsentryDaemon) Start(ctx context.Context) chan error {
	log.WithContext(ctx).Debug("starting transaction sentry")

	daemon.ctx, daemon.cancel = context.WithCancel(ctx)
	backff := backoff.WithContext(backoffjob.NewBackOff(backoff.NewExponentialBackOff()), ctx)

	cerr := make(chan error)
	go func() {
		err := backoff.RetryNotify(daemon.listenJobs, backff, daemon.logOnError)
		if err != nil {
			cerr <- err
		}

		close(cerr)
	}()

	return cerr
}

func (daemon *txsentryDaemon) Stop(ctx context.Context) {
	log.WithContext(ctx).Debug("stopping transaction sentry")

	if !daemon.isReady {
		log.WithContext(ctx).Warn("transaction sentry service is not running")
		return
	}

	daemon.cancel()

	log.WithContext(ctx).Info("stopped transaction sentry")
}

func (daemon *txsentryDaemon) listenJobs() error {
	daemon.isReady = true
	log.WithContext(daemon.ctx).Info("transaction sentry service is ready")

	select {
	case err := <-daemon.jobListener.Listen(daemon.ctx):
		return err
	case <-daemon.ctx.Done():
		daemon.isReady = false
		log.WithContext(daemon.ctx).WithError(daemon.ctx.Err()).Warn("transaction sentry stopped listening")
		return backoff.Permanent(daemon.ctx.Err())
	}
}

func (daemon *txsentryDaemon) logOnError(err error, d time.Duration) {
	log.WithContext(daemon.ctx).WithError(err).Warnf("transaction sentry restarting in... %v", d)
}
