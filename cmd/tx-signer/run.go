package txsigner

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/tessera"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/controllers/amount"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/controllers/creditor"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/keystore"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/secretstore"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/secretstore/hashicorp"
)

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
	}

	// Register KeyStore flags
	hashicorp.InitFlags(runCmd.Flags())
	keystore.InitFlags(runCmd.Flags())
	secretstore.InitFlags(runCmd.Flags())
	tessera.InitFlags(runCmd.Flags())

	// Register Kafka flags
	broker.InitKafkaFlags(runCmd.Flags())
	broker.KafkaTopicTxSigner(runCmd.Flags())
	broker.KafkaTopicTxSender(runCmd.Flags())
	broker.KafkaTopicAccountGenerator(runCmd.Flags())
	broker.KafkaTopicAccountGenerated(runCmd.Flags())
	broker.KafkaTopicTxRecover(runCmd.Flags())

	// Register Faucet flags
	amount.FaucetAmount(runCmd.Flags())
	creditor.FaucetAddress(runCmd.Flags())

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
