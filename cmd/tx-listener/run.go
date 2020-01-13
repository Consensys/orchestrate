package txlistener

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient/rpc"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	chnregclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	provider "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/providers/chain-registry"
	hook "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/session/ethereum/hooks/kafka"
	registryclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry/client"
	storeclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope-store/client"
)

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
	}

	// Register Ethereum client flags
	ethclient.Flags(runCmd.Flags())

	// Register Kafka flags
	broker.InitKafkaFlags(runCmd.Flags())
	broker.KafkaTopicTxDecoder(runCmd.Flags())

	// Register StoreGRPC flags
	storeclient.EnvelopeStoreURL(runCmd.Flags())

	// Listener flags
	provider.Flags(runCmd.Flags())
	chnregclient.Flags(runCmd.Flags())
	registryclient.ContractRegistryURL(runCmd.Flags())

	// Register local flags for handler producer
	hook.InitFlags(runCmd.Flags())

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
