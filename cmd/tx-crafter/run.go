package txcrafter

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/ethclient/rpc"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/http"
	registryclient "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/services/contract-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/tracing/opentracing/jaeger"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/utils"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/contract-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/controllers/amount"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/controllers/blacklist"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/controllers/cooldown"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/controllers/creditor"
	maxbalance "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/controllers/max-balance"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/faucet"
)

func NewRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
	}

	// Register Engine flags
	engine.InitFlags(runCmd.Flags())

	// Register OpenTracing flags
	jaeger.InitFlags(runCmd.Flags())

	// Register HTTP server flags
	http.Hostname(runCmd.Flags())

	// Register Ethereum client flags
	ethclient.URLs(runCmd.Flags())

	// Register Faucet flags
	faucet.Type(runCmd.Flags())
	amount.FaucetAmount(runCmd.Flags())
	blacklist.FaucetBlacklist(runCmd.Flags())
	cooldown.FaucetCooldown(runCmd.Flags())
	creditor.FaucetAddress(runCmd.Flags())
	maxbalance.FaucetMaxBalance(runCmd.Flags())

	// Register Crafter flags
	contractregistry.ABIs(runCmd.Flags())

	// Register Kafka flags
	broker.KafkaAddresses(runCmd.Flags())
	broker.KafkaGroup(runCmd.Flags())
	broker.KafkaTopicTxCrafter(runCmd.Flags())
	broker.KafkaTopicTxNonce(runCmd.Flags())
	broker.InitKafkaSASLTLSFlags(runCmd.Flags())
	broker.KafkaTopicTxRecover(runCmd.Flags())

	// Contract Registry
	registryclient.ContractRegistryGRPCTarget(runCmd.Flags())

	return runCmd
}

func run(_ *cobra.Command, _ []string) {
	// Create app
	ctx, cancel := context.WithCancel(context.Background())

	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) { cancel() })
	defer sig.Close()

	// Start application
	Start(ctx)
}
