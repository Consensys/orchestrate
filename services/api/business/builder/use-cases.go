package builder

import (
	"github.com/Shopify/sarama"
	pkgsarama "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases/faucets"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
	keymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client"
)

type useCases struct {
	*jobUseCases
	*scheduleUseCases
	*transactionUseCases
	*faucetUseCases
	*chainUseCases
	*contractUseCases
	*accountUseCases
}

func NewUseCases(
	db store.DB,
	appMetrics metrics.TransactionSchedulerMetrics,
	keyManagerClient keymanager.KeyManagerClient,
	ec ethclient.Client,
	producer sarama.SyncProducer,
	topicsCfg *pkgsarama.KafkaTopicConfig,
) usecases.UseCases {

	chainUseCases := newChainUseCases(db, ec)
	contractUseCases := newContractUseCases(db)
	faucetUseCases := newFaucetUseCases(db)
	getFaucetCandidateUC := faucets.NewGetFaucetCandidateUseCase(faucetUseCases.SearchFaucets(), ec)
	scheduleUseCases := newScheduleUseCases(db)
	jobUseCases := newJobUseCases(db, appMetrics, producer, topicsCfg, chainUseCases.GetChain())
	transactionUseCases := newTransactionUseCases(db, chainUseCases.SearchChains(), getFaucetCandidateUC, scheduleUseCases, jobUseCases, contractUseCases.GetContract())
	accountUseCases := newAccountUseCases(db, keyManagerClient, chainUseCases.SearchChains(), transactionUseCases.SendTransaction(), getFaucetCandidateUC)

	return &useCases{
		jobUseCases:         jobUseCases,
		scheduleUseCases:    scheduleUseCases,
		transactionUseCases: transactionUseCases,
		faucetUseCases:      faucetUseCases,
		chainUseCases:       chainUseCases,
		contractUseCases:    contractUseCases,
		accountUseCases:     accountUseCases,
	}
}
