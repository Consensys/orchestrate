package stress

import (
	"context"
	"sync"
	"time"

	"github.com/consensys/orchestrate/pkg/backoff"
	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app"
	ethclient "github.com/consensys/orchestrate/pkg/toolkit/ethclient/rpc"
	"github.com/sirupsen/logrus"

	loader "github.com/consensys/orchestrate/handlers/loader/sarama"
	"github.com/consensys/orchestrate/handlers/offset"
	broker "github.com/consensys/orchestrate/pkg/broker/sarama"
	"github.com/consensys/orchestrate/pkg/engine"
	orchestrateclient "github.com/consensys/orchestrate/pkg/sdk/client"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/consensys/orchestrate/tests/handlers"
	"github.com/consensys/orchestrate/tests/handlers/consumer"
	"github.com/consensys/orchestrate/tests/handlers/dispatcher"
	utils2 "github.com/consensys/orchestrate/tests/service/stress/utils"
	"github.com/consensys/orchestrate/tests/utils/chanregistry"
	"github.com/spf13/viper"
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

	httpClient := app.NewHTTPClient(viper.GetViper())
	backoffConf := orchestrateclient.NewConfigFromViper(viper.GetViper(),
		backoff.IncrementalBackOff(time.Second, time.Second*5, time.Minute))
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
