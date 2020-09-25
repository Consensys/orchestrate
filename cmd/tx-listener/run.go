package txlistener

import (
	"os"

	txsentry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-sentry"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	chnregclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	registryclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/client"
	txschedulerclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
	txlistener "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener"
	provider "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/providers/chain-registry"
)

var cmdErr error

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		RunE:  run,
		PreRun: func(cmd *cobra.Command, args []string) {
			utils.PreRunBindFlags(viper.GetViper(), cmd.Flags(), "tx-listener")
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			if err := errors.CombineErrors(cmdErr, cmd.Context().Err()); err != nil {
				os.Exit(1)
			}
		},
	}

	// Register Kafka flags
	broker.InitKafkaFlags(runCmd.Flags())

	// Listener flags
	provider.Flags(runCmd.Flags())
	chnregclient.Flags(runCmd.Flags())
	registryclient.ContractRegistryURL(runCmd.Flags())
	txschedulerclient.Flags(runCmd.Flags())
	txsentry.Flags(runCmd.Flags())

	return runCmd
}

func run(cmd *cobra.Command, _ []string) error {
	return txlistener.Run(cmd.Context())
}
