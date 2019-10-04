package txcrafter

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient/rpc"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/controllers/amount"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/controllers/blacklist"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/controllers/cooldown"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/controllers/creditor"
	maxbalance "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/controllers/max-balance"
	registryclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry/client"
)

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
	}

	// Register Ethereum client flags
	ethclient.URLs(runCmd.Flags())

	// Register Faucet flags
	amount.FaucetAmount(runCmd.Flags())
	blacklist.FaucetBlacklist(runCmd.Flags())
	cooldown.FaucetCooldown(runCmd.Flags())
	creditor.FaucetAddress(runCmd.Flags())
	maxbalance.FaucetMaxBalance(runCmd.Flags())

	// Register Kafka flags
	broker.InitKafkaFlags(runCmd.Flags())
	broker.KafkaTopicTxCrafter(runCmd.Flags())
	broker.KafkaTopicTxNonce(runCmd.Flags())
	broker.KafkaTopicTxRecover(runCmd.Flags())

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
