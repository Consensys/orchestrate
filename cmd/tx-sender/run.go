package txsender

import (
	"context"
	"os"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/spf13/cobra"
	noncechecker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/nonce/checker"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	chnregclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	storeclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/nonce"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/nonce/redis"
	txsender "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-sender"
)

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
	}

	// Register Kafka flags
	broker.InitKafkaFlags(runCmd.Flags())
	broker.KafkaTopicTxSender(runCmd.Flags())
	broker.KafkaTopicTxRecover(runCmd.Flags())

	// Chain Registry
	chnregclient.Flags(runCmd.Flags())

	// Register StoreGRPC flags
	storeclient.EnvelopeStoreURL(runCmd.Flags())

	// Register Nonce Manager flags
	nonce.Type(runCmd.Flags())
	redis.Flags(runCmd.Flags())
	noncechecker.Flags(runCmd.Flags())

	return runCmd
}

func run(_ *cobra.Command, _ []string) {
	_ = txsender.Start(context.Background())

	done := make(chan struct{})

	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) {
		err := txsender.Stop(context.Background())
		if err != nil {
			log.WithoutContext().WithError(err).Errorf("Application did not shutdown properly")
		} else {
			log.WithoutContext().WithError(err).Infof("Application gracefully closed")
		}
		close(done)
	})
	<-done

	sig.Close()
}
