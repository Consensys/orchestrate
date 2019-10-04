package txsender

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient/rpc"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	storeclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope-store/client"
)

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
	}

	// Register Ethereum client flags
	ethclient.URLs(runCmd.Flags())

	// Register Kafka flags
	broker.InitKafkaFlags(runCmd.Flags())
	broker.KafkaTopicTxSender(runCmd.Flags())
	broker.KafkaTopicTxRecover(runCmd.Flags())

	// Register StoreGRPC flags
	storeclient.EnvelopeStoreURL(runCmd.Flags())

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
