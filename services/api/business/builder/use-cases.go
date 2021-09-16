package builder

import (
	pkgsarama "github.com/ConsenSys/orchestrate/pkg/broker/sarama"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/ethclient"
	usecases "github.com/ConsenSys/orchestrate/services/api/business/use-cases"
	"github.com/ConsenSys/orchestrate/services/api/business/use-cases/faucets"
	"github.com/ConsenSys/orchestrate/services/api/metrics"
	"github.com/ConsenSys/orchestrate/services/api/store"
	qkmclient "github.com/consensys/quorum-key-manager/pkg/client"
	"github.com/Shopify/sarama"
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
	keyManagerClient qkmclient.EthClient,
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
	transactionUseCases := newTransactionUseCases(db, chainUseCases.SearchChains(), getFaucetCandidateUC, 
		scheduleUseCases, jobUseCases, contractUseCases.GetContract())
	accountUseCases := newAccountUseCases(db, keyManagerClient, chainUseCases.SearchChains(), 
		transactionUseCases.SendTransaction(), getFaucetCandidateUC)

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
