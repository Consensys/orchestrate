package txlistener

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/ethclient/rpc"
	handler "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/tx-listener/handler/base"
	listener "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/tx-listener/listener/base"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	storeclient "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/services/envelope-store/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/utils"
)

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
	}

	// Register Engine flags
	engine.InitFlags(runCmd.Flags())

	// Register Ethereum client flags
	ethclient.URLs(runCmd.Flags())

	// Register Kafka flags
	broker.KafkaAddresses(runCmd.Flags())
	broker.KafkaGroup(runCmd.Flags())
	broker.KafkaTopicTxDecoder(runCmd.Flags())
	broker.InitKafkaSASLTLSFlags(runCmd.Flags())

	// Register StoreGRPC flags
	storeclient.EnvelopeStoreGRPCTarget(runCmd.Flags())

	// Listener flags
	listener.InitFlags(runCmd.Flags())
	handler.InitFlags(runCmd.Flags())

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
