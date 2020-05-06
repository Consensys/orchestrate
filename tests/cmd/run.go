package main

import (
	"context"
	"os"

	traefiklog "github.com/containous/traefik/v2/pkg/log"
	"github.com/spf13/cobra"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/jwt/generator"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	e2e "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/cucumber"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/cucumber/steps"
)

func NewRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
	}

	// Set pkglog flags
	log.Level(runCmd.Flags())
	log.Format(runCmd.Flags())

	// Register Kafka flags
	broker.InitKafkaFlags(runCmd.Flags())
	broker.KafkaTopicTxCrafter(runCmd.Flags())
	broker.KafkaTopicTxSigner(runCmd.Flags())
	broker.KafkaTopicTxSender(runCmd.Flags())
	broker.KafkaTopicTxDecoded(runCmd.Flags())
	broker.KafkaTopicAccountGenerator(runCmd.Flags())
	broker.KafkaTopicAccountGenerated(runCmd.Flags())

	// Register Cucumber flag
	cucumber.InitFlags(runCmd.Flags())
	steps.InitFlags(runCmd.Flags())

	// Register Multi-Tenancy flags
	multitenancy.Enabled(runCmd.Flags())
	auth.Flags(runCmd.Flags())
	generator.PrivateKey(runCmd.Flags())

	return runCmd
}

func run(_ *cobra.Command, _ []string) {

	done := make(chan struct{})
	go func() {
		_ = e2e.Start(context.Background())
		close(done)
	}()

	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) {
		err := e2e.Stop(context.Background())
		if err != nil {
			traefiklog.WithoutContext().WithError(err).Errorf("Application did not shutdown properly")
		} else {
			traefiklog.WithoutContext().WithError(err).Errorf("Application gracefully closed")
		}
	})

	<-done

	sig.Close()
}
