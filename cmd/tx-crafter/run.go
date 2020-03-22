package txcrafter

import (
	"context"
	"os"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/spf13/cobra"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	chnregclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	registryclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/client"
	txcrafter "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-crafter"
)

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
	}

	// Register Kafka flags
	broker.InitKafkaFlags(runCmd.Flags())
	broker.KafkaTopicTxCrafter(runCmd.Flags())
	broker.KafkaTopicTxRecover(runCmd.Flags())

	// Chain Registry
	chnregclient.Flags(runCmd.Flags())

	// Contract Registry
	registryclient.ContractRegistryURL(runCmd.Flags())

	return runCmd
}

func run(_ *cobra.Command, _ []string) {
	_ = txcrafter.Start(context.Background())

	done := make(chan struct{})

	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) {
		err := txcrafter.Stop(context.Background())
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
