package txsigner

import (
	"context"

	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/nonce"
	txscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/client"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/app"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	keymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client"
)

// New Utility function used to initialize a new service
func New(ctx context.Context) (*app.App, error) {
	config := NewConfig(viper.GetViper())

	// Initialize dependencies
	sarama.InitConsumerGroup(ctx, config.GroupName)
	sarama.InitSyncProducer(ctx)
	keymanager.Init()
	txscheduler.Init()
	ethclient.Init(ctx)
	nonce.Init(ctx)

	return NewTxSigner(
		config,
		sarama.GlobalConsumerGroup(),
		sarama.GlobalSyncProducer(),
		keymanager.GlobalClient(),
		txscheduler.GlobalClient(),
		ethclient.GlobalClient(),
		nonce.GlobalManager(),
	)
}
