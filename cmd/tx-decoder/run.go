package main

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/ethclient/rpc"
	producer "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/producer/tx-decoder"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/tracing/opentracing/jaeger"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/utils"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/contract-registry"
)

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
	}

	// Register Engine flags
	engine.InitFlags(runCmd.Flags())

	// Register HTTP server flags
	http.Hostname(runCmd.Flags())

	// Register Ethereum client flags
	rpc.URLs(runCmd.Flags())

	// Register Decoder flags
	contractregistry.ABIs(runCmd.Flags())

	// Register Opentracing flags
	jaeger.InitFlags(runCmd.Flags())

	// Register Kafka flags
	broker.KafkaAddresses(runCmd.Flags())
	broker.KafkaGroup(runCmd.Flags())
	broker.KafkaTopicTxDecoded(runCmd.Flags())
	broker.KafkaTopicTxDecoder(runCmd.Flags())
	broker.InitKafkaSASLTLSFlags(runCmd.Flags())

	// Registers local flags for handler producer
	producer.InitFlags(runCmd.Flags())

	return runCmd
}

func run(cmd *cobra.Command, args []string) {
	// Create app
	ctx, cancel := context.WithCancel(context.Background())

	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) { cancel() })
	defer sig.Close()

	// Start application
	Start(ctx)
}
