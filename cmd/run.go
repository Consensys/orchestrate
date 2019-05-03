package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/nonce.git/nonce"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/nonce.git/nonce/redis"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/opentracing/jaeger"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/utils"

	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-nonce.git/app"
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
	ethclient.URLs(runCmd.Flags())

	// Register Opentracing flags
	jaeger.Host(runCmd.Flags())
	jaeger.Port(runCmd.Flags())
	jaeger.SamplerParam(runCmd.Flags())
	jaeger.SamplerType(runCmd.Flags())

	// Register Kafka flags
	broker.KafkaAddresses(runCmd.Flags())
	broker.KafkaGroup(runCmd.Flags())
	broker.KafkaTopicTxNonce(runCmd.Flags())
	broker.KafkaTopicTxSigner(runCmd.Flags())

	// Register Nonce Manager flags
	nonce.Type(runCmd.Flags())
	redis.Address(runCmd.Flags())
	redis.LockTimeout(runCmd.Flags())
	redis.RedisNonceExpirationTime(runCmd.Flags())

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
