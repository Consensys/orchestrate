package steps

import (
	"context"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/jwt/generator"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient/rpc"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/client"
	noncememory "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/nonce/memory"
	txscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/cucumber/alias"
)

// Init initialize handlers
func Init(ctx context.Context) {
	broker.InitSyncProducer(ctx)
	generator.Init(ctx)
	chainregistry.Init(ctx)
	alias.Init(ctx)
	contractregistry.Init(ctx, viper.GetString(contractregistry.ContractRegistryURLViperKey))
	txscheduler.Init()
	noncememory.Init(ctx)
	rpc.Init(ctx)
}
