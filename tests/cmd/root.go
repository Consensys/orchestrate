package main

import (
	broker "github.com/consensys/orchestrate/pkg/broker/sarama"
	"github.com/consensys/orchestrate/pkg/multitenancy"
	"github.com/consensys/orchestrate/pkg/toolkit/app/auth"
	"github.com/consensys/orchestrate/pkg/toolkit/app/auth/jwt/generator"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
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
