package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/opentracing/jaeger"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/multi-vault.git/keystore"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/multi-vault.git/secretstore"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/multi-vault.git/secretstore/hashicorp"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-signer.git/app"
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

	// Register Opentracing flags
	jaeger.InitFlags(runCmd.Flags())

	// Register KeyStore flags
	hashicorp.InitFlags(runCmd.Flags())
	keystore.InitFlags(runCmd.Flags())
	secretstore.InitFlags(runCmd.Flags())

	// Register Kafka flags
	broker.KafkaAddresses(runCmd.Flags())
	broker.KafkaGroup(runCmd.Flags())
	broker.KafkaTopicTxSigner(runCmd.Flags())
	broker.KafkaTopicTxSender(runCmd.Flags())
	broker.KafkaTopicWalletGenerator(runCmd.Flags())
	broker.KafkaTopicWalletGenerated(runCmd.Flags())

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
