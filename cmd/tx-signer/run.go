package txsigner

import (
	"context"
	"os"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	chnregclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/keystore"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/secretstore"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/secretstore/hashicorp"
	txschedulerclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
	txsigner "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-signer"
)

var cmdErr error

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
		PreRun: func(cmd *cobra.Command, args []string) {
			utils.PreRunBindFlags(viper.GetViper(), cmd.Flags(), "tx-signer")
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			if err := errors.CombineErrors(cmdErr, cmd.Context().Err()); err != nil {
				os.Exit(1)
			}
		},
	}

	// Register KeyStore flags
	hashicorp.InitFlags(runCmd.Flags())
	keystore.InitFlags(runCmd.Flags())
	secretstore.InitFlags(runCmd.Flags())

	// Register Kafka flags
	broker.InitKafkaFlags(runCmd.Flags())
	broker.KafkaTopicTxSigner(runCmd.Flags())
	broker.KafkaTopicTxSender(runCmd.Flags())
	broker.KafkaTopicAccountGenerator(runCmd.Flags())
	broker.KafkaTopicAccountGenerated(runCmd.Flags())
	broker.KafkaTopicTxRecover(runCmd.Flags())

	// Internal API clients
	chnregclient.Flags(runCmd.Flags())
	txschedulerclient.Flags(runCmd.Flags())

	return runCmd
}

func run(cmd *cobra.Command, _ []string) {
	ctx, cancel := context.WithCancel(cmd.Context())
	logger := log.FromContext(ctx)

	// Start microservice
	go func() {
		if err := <-txsigner.Start(ctx); err != nil {
			cmdErr = errors.CombineErrors(cmdErr, err)
			log.WithoutContext().WithError(err).Errorf("Microservice raised an error")
		}
		cancel()
	}()

	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) {
		cancel()
	})

	// Stop when get context canceled
	<-ctx.Done()
	err := txsigner.Stop(ctx)
	if err != nil {
		cmdErr = errors.CombineErrors(cmdErr, err)
		logger.WithError(err).Errorf("Microservice did not shutdown properly")
	} else {
		logger.Infof("Microservice gracefully closed")
	}

	sig.Close()
}
