package stress

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/backoff"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http"

	"github.com/spf13/viper"
	loader "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/handlers/loader/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/handlers/offset"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	orchestrateclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/handlers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/handlers/consumer"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/handlers/dispatcher"
	utils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/service/stress/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/utils/chanregistry"
)

const component = "stress-test"

var (
	workload *WorkLoadService
	cancel   func()
)

// Start starts application
func Start(ctx context.Context) error {
	logger := log.WithContext(ctx).SetComponent(component)
	ctx = log.With(ctx, logger)
	logger.Info("starting execution...")

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

	httpClient := http.NewClient(http.NewConfig(viper.GetViper()))
	backoffConf := orchestrateclient.NewConfigFromViper(viper.GetViper(),
		backoff.IncrementalBackOff(time.Second*5, time.Minute))
	client := orchestrateclient.NewHTTPClient(httpClient, backoffConf)

	workload = NewService(cfg,
		chanregistry.GlobalChanRegistry(),
		client,
		ethclient.GlobalClient(),
		broker.GlobalSyncProducer())

	// Start consuming on every topics of interest
	var topics []string
	for _, viprTopicKey := range utils2.Topics {
		topics = append(topics, viper.GetString(viprTopicKey))
	}

	cg := consumer.NewEmbeddingConsumerGroupHandler(engine.GlobalEngine())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		logger.WithField("topics", topics).Info("connecting to kafka")

		err := broker.Consume(ctx, topics, cg)
		if err != nil {
			gerr = errors.CombineErrors(gerr, err)
			logger.WithError(err).Error("error on consumer")
		}

		cancel()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		<-cg.IsReady()
		err := workload.Run(ctx)
		if err != nil {
			gerr = errors.CombineErrors(gerr, err)
			logger.WithError(err).Error("error on workload test")
		}

		cancel()
		wg.Done()
	}()

	wg.Wait()
	return gerr
}

func Stop(ctx context.Context) error {
	log.WithContext(ctx).Info("stopping stress test execution...")
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
			ethclient.Init(ctx)
		},
		// Initialize ConsumerGroup
		func() {
			broker.InitConsumerGroup(ctx, viper.GetString(broker.ConsumerGroupNameViperKey))
		},
		// Initialize Handlers
		func() {
			// Prepare topics map for dispatcher
			topics := make(map[string]string)
			for topic, viprTopicKey := range utils2.Topics {
				topics[viper.GetString(viprTopicKey)] = topic
			}
			dispatcher.SetKeyOfFuncs(
				dispatcher.LongKeyOf(topics),
				dispatcher.LabelKey(topics),
			)
			handlers.Init(ctx)
		},
		// Initialize logger
		func() {
			cfg := log.NewConfig(viper.GetViper())
			// Create and configure logger
			logger := logrus.StandardLogger()
			_ = log.ConfigureLogger(cfg, logger)
		},
	)
}
