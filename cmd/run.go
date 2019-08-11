package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/tracing/opentracing/jaeger"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/abi/registry"
	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/ethclient/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/faucet.git/controllers/amount"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/faucet.git/controllers/blacklist"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/faucet.git/controllers/cooldown"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/faucet.git/controllers/creditor"
	maxbalance "gitlab.com/ConsenSys/client/fr/core-stack/service/faucet.git/controllers/max-balance"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/faucet.git/faucet"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-crafter.git/app"
)

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
	}

	// Register Engine flags
	engine.InitFlags(runCmd.Flags())

	// Register Opentracing flags
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
	registry.ABIs(runCmd.Flags())

	// Register Kafka flags
	broker.KafkaAddresses(runCmd.Flags())
	broker.KafkaGroup(runCmd.Flags())
	broker.KafkaTopicTxCrafter(runCmd.Flags())
	broker.KafkaTopicTxNonce(runCmd.Flags())
	broker.InitKafkaSASLTLSFlags(runCmd.Flags())
	broker.KafkaTopicTxRecover(runCmd.Flags())

	return runCmd
}

func run(cmd *cobra.Command, args []string) {
	// Create app
	ctx, cancel := context.WithCancel(context.Background())

	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) { cancel() })
	defer sig.Close()

	// Start application
	app.Start(ctx)
}
