package stress

import (
	"context"
	"fmt"
	"sync"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	loader "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/handlers/loader/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/handlers/offset"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine"
	pkglog "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/handlers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/handlers/consumer"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/handlers/dispatcher"
	utils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/service/stress/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/utils/chanregistry"
)

var (
	workload *WorkLoadService
	cancel   func()
)

// Start starts application
func Start(ctx context.Context) error {
	log.FromContext(ctx).Info("Starting execution...")

	var gerr error
	// Create context for application
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	initComponents(ctx)
	registerHandlers()

	cfg, err := InitConfig(viper.GetViper())
	if err != nil {
		return err
	}

	workload = NewService(cfg,
		chanregistry.GlobalChanRegistry(),
		chainregistry.GlobalClient(),
		client.GlobalClient(),
		broker.GlobalSyncProducer())

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
		err := workload.Run(ctx)
		if err != nil {
			gerr = err
			log.FromContext(ctx).WithError(err).Error("error on workload test")
		}

		cancel()
		wg.Done()
	}()

	wg.Wait()
	return gerr
}

func Stop(ctx context.Context) error {
	log.FromContext(ctx).Info("Stopping Cucumber execution...")
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
		func() {
			broker.InitSyncProducer(ctx)
			chainregistry.Init(ctx)
			client.Init()
		},
		// Initialize ConsumerGroup
		func() {
			broker.InitConsumerGroup(ctx, fmt.Sprintf("group-cucumber-%s", utils.RandomString(10)))
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
	)
}
