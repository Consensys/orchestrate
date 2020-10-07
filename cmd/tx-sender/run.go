package txsender

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	noncechecker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/nonce/checker"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/redis"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	chnregclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/nonce"
	txschedulerclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
	txsender "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-sender"
)

var cmdErr error

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		RunE:  run,
		PreRun: func(cmd *cobra.Command, args []string) {
			utils.PreRunBindFlags(viper.GetViper(), cmd.Flags(), "tx-sender")
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			if err := errors.CombineErrors(cmdErr, cmd.Context().Err()); err != nil {
				os.Exit(1)
			}
		},
	}

	// Register Kafka flags
	broker.InitKafkaFlags(runCmd.Flags())
	broker.KafkaTopicTxSender(runCmd.Flags())
	broker.KafkaTopicTxRecover(runCmd.Flags())

	// Internal API clients
	chnregclient.Flags(runCmd.Flags())
	txschedulerclient.Flags(runCmd.Flags())

	// Register Nonce Manager flags
	nonce.Type(runCmd.Flags())
	redis.Flags(runCmd.Flags())
	noncechecker.Flags(runCmd.Flags())

	return runCmd
}

func run(cmd *cobra.Command, _ []string) error {
	return txsender.Run(cmd.Context())
}
