package main

import (
	"github.com/ConsenSys/orchestrate/pkg/auth"
	"github.com/ConsenSys/orchestrate/pkg/auth/jwt/generator"
	broker "github.com/ConsenSys/orchestrate/pkg/broker/sarama"
	"github.com/ConsenSys/orchestrate/pkg/log"
	"github.com/ConsenSys/orchestrate/pkg/multitenancy"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	command := &cobra.Command{
		Use:              "run",
		TraverseChildren: true,
		SilenceUsage:     true,
	}

	// Set pkglog flags
	log.Flags(command.Flags())

	// Register Kafka flags
	broker.KafkaConsumerFlags(command.Flags())
	broker.KafkaTopicTxSender(command.Flags())
	broker.KafkaTopicTxDecoded(command.Flags())

	command.AddCommand(NewRunE2ECommand())
	command.AddCommand(NewRunStressTestCommand())

	// Register Multi-Tenancy flags
	multitenancy.Flags(command.Flags())
	auth.Flags(command.Flags())
	generator.PrivateKey(command.Flags())

	if err := command.Execute(); err != nil {
		logrus.WithError(err).Fatalf("test: execution failed")
	}

	logrus.Infof("test: execution completed")
}
