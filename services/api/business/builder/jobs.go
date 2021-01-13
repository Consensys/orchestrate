package builder

import (
	"github.com/Shopify/sarama"
	pkgsarama "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases/jobs"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
)

type jobUseCases struct {
	createJob   usecases.CreateJobUseCase
	getJob      usecases.GetJobUseCase
	startJob    usecases.StartJobUseCase
	resendJobTx usecases.ResendJobTxUseCase
	updateJob   usecases.UpdateJobUseCase
	searchJobs  usecases.SearchJobsUseCase
}

func newJobUseCases(
	db store.DB,
	appMetrics metrics.TransactionSchedulerMetrics,
	producer sarama.SyncProducer,
	topicsCfg *pkgsarama.KafkaTopicConfig,
	getChainUC usecases.GetChainUseCase,
) *jobUseCases {
	startJobUC := jobs.NewStartJobUseCase(db, producer, topicsCfg, appMetrics)
	updateChildrenUC := jobs.NewUpdateChildrenUseCase(db)
	startNextJobUC := jobs.NewStartNextJobUseCase(db, startJobUC)

	return &jobUseCases{
		createJob:   jobs.NewCreateJobUseCase(db, getChainUC),
		getJob:      jobs.NewGetJobUseCase(db),
		searchJobs:  jobs.NewSearchJobsUseCase(db),
		updateJob:   jobs.NewUpdateJobUseCase(db, updateChildrenUC, startNextJobUC, appMetrics),
		startJob:    startJobUC,
		resendJobTx: jobs.NewResendJobTxUseCase(db, producer, topicsCfg),
	}
}

func (u *jobUseCases) CreateJob() usecases.CreateJobUseCase {
	return u.createJob
}

func (u *jobUseCases) GetJob() usecases.GetJobUseCase {
	return u.getJob
}

func (u *jobUseCases) SearchJobs() usecases.SearchJobsUseCase {
	return u.searchJobs
}

func (u *jobUseCases) StartJob() usecases.StartJobUseCase {
	return u.startJob
}

func (u *jobUseCases) ResendJobTx() usecases.ResendJobTxUseCase {
	return u.resendJobTx
}

func (u *jobUseCases) UpdateJob() usecases.UpdateJobUseCase {
	return u.updateJob
}
