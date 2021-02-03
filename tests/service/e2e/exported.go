package e2e

import (
	"context"
	"fmt"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	loader "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/handlers/loader/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/handlers/offset"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/redis"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	pkglog "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/handlers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/handlers/consumer"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/handlers/dispatcher"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/service/e2e/cucumber"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/service/e2e/cucumber/alias"
	utils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/service/e2e/utils"
)

var (
	cancel func()
)

// Start starts application
func Start(ctx context.Context) error {
	log.FromContext(ctx).Info("Cucumber: starting execution...")

	var gerr error
	// Create context for application
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	initComponents(ctx)
	registerHandlers()

	if err := importTestIdentities(ctx); err != nil {
		return err
	}

	// Start consuming on every topics of interest
	var topics []string
	for _, topic := range utils2.TOPICS {
		topics = append(topics, viper.GetString(fmt.Sprintf("topic.%v", topic)))
	}

	cg := consumer.NewEmbeddingConsumerGroupHandler(engine.GlobalEngine())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		log.FromContext(ctx).WithFields(logrus.Fields{
			"topics": topics,
		}).Info("connecting")

		err := broker.Consume(ctx, topics, cg)
		if err != nil {
			log.FromContext(ctx).WithError(err).Error("error on consumer")
		}

		cancel()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		<-cg.IsReady()
		err := cucumber.Run(cucumber.GlobalOptions())
		if err != nil {
			gerr = err
			log.FromContext(ctx).WithError(err).Error("error on cucumber")
		}

		cancel()
		wg.Done()
	}()

	wg.Wait()
	return gerr
}

func Stop(ctx context.Context) error {
	log.FromContext(ctx).Info("Cucumber: stopping execution...")
	cancel()
	return nil
}

func registerHandlers() {
	// Generic handlers on every worker
	engine.Register(loader.Loader)
	engine.Register(offset.Marker)
	engine.Register(dispatcher.GlobalHandler())
}

func initComponents(ctx context.Context) {
	utils.InParallel(
		// Initialize Engine
		func() {
			engine.Init(ctx)
		},
		// Initialize ConsumerGroup
		func() {
			broker.InitConsumerGroup(ctx, viper.GetString(broker.ConsumerGroupNameViperKey))
		},
		// Initialize Handlers
		func() {
			// Prepare topics map for dispatcher
			topics := make(map[string]string)
			for _, topic := range utils2.TOPICS {
				topics[viper.GetString(fmt.Sprintf("topic.%v", topic))] = topic
			}
			dispatcher.SetKeyOfFuncs(
				dispatcher.LongKeyOf(topics),
				dispatcher.LabelKey(topics),
				dispatcher.ShortKeyOf(topics),
			)
			handlers.Init(ctx)
		},
		// Initialize logger
		func() {
			cfg := pkglog.NewConfig(viper.GetViper())
			// Create and configure logger
			logger := logrus.StandardLogger()
			_ = pkglog.ConfigureLogger(cfg, logger)
		},
		// Initialize Nonce Manager
		func() {
			redis.Init()
		},
		// Initialize cucumber handlers
		func() {
			cucumber.Init(ctx)
		},
	)
}

// We import account define at Global Aliases
func importTestIdentities(ctx context.Context) error {
	orchestrateClient := client.GlobalClient()
	aliases := alias.GlobalAliasRegistry()

	var privKeys []interface{}
	for _, netName := range []string{"besu", "quorum", "geth"} {
		if netNodes, ok := aliases.Get("global.nodes." + netName); ok {
			for _, node := range netNodes.([]interface{}) {
				if pKeys, ok := node.(map[string]interface{})["fundedPrivateKeys"].([]interface{}); ok {
					privKeys = append(privKeys, pKeys...)
				}
			}
		}
	}

	for _, privKey := range privKeys {
		resp, err := orchestrateClient.ImportAccount(ctx, &api.ImportAccountRequest{
			PrivateKey: privKey.(string),
		})

		if err != nil {
			if errors.IsAlreadyExistsError(err) {
				continue
			}

			return err
		}

		log.FromContext(ctx).WithField("address", resp.Address).Info("Account imported successfully")
	}

	return nil
}
