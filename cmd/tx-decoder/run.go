package txdecoder

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	producer "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/producer/tx-decoder"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	registryclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry/client"
)

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
	}

	// Register Kafka flags
	broker.InitKafkaFlags(runCmd.Flags())
	broker.KafkaTopicTxDecoded(runCmd.Flags())
	broker.KafkaTopicTxDecoder(runCmd.Flags())

	// Registers local flags for handler producer
	producer.InitFlags(runCmd.Flags())

	// Contract Registry
	registryclient.ContractRegistryURL(runCmd.Flags())

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
