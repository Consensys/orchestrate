package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/jwt/generator"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
)

var cmdErr error

func main() {
	command := &cobra.Command{
		Use:              "run",
		TraverseChildren: true,
	}

	// Set pkglog flags
	log.Level(command.Flags())
	log.Format(command.Flags())

	// Register Kafka flags
	broker.InitKafkaFlags(command.Flags())
	broker.KafkaTopicTxCrafter(command.Flags())
	broker.KafkaTopicTxSigner(command.Flags())
	broker.KafkaTopicTxSender(command.Flags())
	broker.KafkaTopicTxDecoded(command.Flags())
	broker.KafkaTopicAccountGenerator(command.Flags())
	broker.KafkaTopicAccountGenerated(command.Flags())

	command.AddCommand(NewRunE2ECommand())
	command.AddCommand(NewRunStressTestCommand())

	// Register Multi-Tenancy flags
	multitenancy.Enabled(command.Flags())
	auth.Flags(command.Flags())
	generator.PrivateKey(command.Flags())

	if err := command.Execute(); err != nil {
		logrus.WithError(err).Fatalf("test: execution failed")
	}

	logrus.Infof("test: execution completed")
}
