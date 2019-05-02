package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	storegrpc "gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/store/grpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/opentracing/jaeger"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-sender.git/app"
)

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
	}

	// Register Engine flags
	engine.InitFlags(runCmd.Flags())

	// Register Opentracing flags
	jaeger.InitFlags(runCmd.Flags())

	// Register HTTP server flags
	http.Hostname(runCmd.Flags())

	// Register Ethereum client flags
	ethclient.URLs(runCmd.Flags())

	// Register Kafka flags
	broker.KafkaAddresses(runCmd.Flags())
	broker.KafkaGroup(runCmd.Flags())
	broker.KafkaTopicTxSender(runCmd.Flags())

	// Register StoreGRPC flags
	storegrpc.StoreTarget(runCmd.Flags())

	return runCmd
}

func run(cmd *cobra.Command, args []string) {
	// Create app
	ctx, cancel := context.WithCancel(context.Background())

	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) { cancel() })
	defer sig.Close()

	// Start application
	app.Start(ctx)
}
