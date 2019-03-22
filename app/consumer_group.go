package app

import (
	"sync"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	infSarama "gitlab.com/ConsenSys/client/fr/core-stack/infra/sarama.git"
	coreworker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-signer.git/app/worker"
)

type handler struct {
	app *App

	cleanOnce *sync.Once
	in        chan *sarama.ConsumerMessage
	worker    *coreworker.Worker
	logger    *log.Entry
}

// Setup configure handler
func (h *handler) Setup(s sarama.ConsumerGroupSession) error {
	h.worker = worker.CreateWorker(h.app.infra, infSarama.NewSimpleOffsetMarker(s))

	// Pipe sarama message channel into worker
	in := make(chan interface{})
	go func() {
		// Pipe channels for interface compatibility
		for i := range h.in {
			in <- i
		}
		close(in)
	}()

	// Run worker
	go h.worker.Run(in)

	h.logger.WithFields(log.Fields{
		"kafka.generation_id": s.GenerationID(),
		"kafka.member_id":     s.MemberID(),
	}).Infof("consumer-group: ready to consume claims %v", s.Claims())

	return nil
}

// ConsumeClaim consume messages from queue
func (h *handler) ConsumeClaim(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) error {
	logger := h.logger.WithFields(log.Fields{
		"kafka.topic":     c.Topic(),
		"kafka.partition": c.Partition(),
	})

	logger.Infof("consumer-group: start consuming claim (offset=%v)", c.InitialOffset())

	// Pipe messages into input channel
consumeLoop:
	for {
		select {
		case msg, ok := <-c.Messages():
			if !ok {
				break consumeLoop
			}
			select {
			case h.in <- msg:
				continue
			case <-s.Context().Done():
				break consumeLoop
			}
		case <-s.Context().Done():
			break consumeLoop
		}
	}
	// Close worker
	h.worker.Close()

	// Wait for worker to be done then leave
	<-h.worker.Done()

	logger.Infof("consumer-group: stoped consuming claim")

	return nil
}

// Cleanup cleans handler
func (h *handler) Cleanup(s sarama.ConsumerGroupSession) error {
	h.cleanOnce.Do(func() {
		close(h.in)
	})
	return nil
}

func initConsumerGroup(app *App) {
	// Retrieve group name from config
	group := viper.GetString("worker.group")

	// Add group field in logger
	logger := log.StandardLogger().WithFields(log.Fields{
		"kafka.group": group,
	})

	// Create group
	g, err := sarama.NewConsumerGroupFromClient(group, app.infra.SaramaClient)
	if err != nil {
		logger.WithError(err).Fatalf("consumer-group: error creating consumer group")
	}

	// Attach consumer group and handler to app
	app.saramaConsumerGroup = g
	app.saramaHandler = &handler{
		app:       app,
		cleanOnce: &sync.Once{},
		in:        make(chan *sarama.ConsumerMessage),
		logger:    logger,
	}

	// Wait for app to be done and then close
	go func() {
		<-app.Done()
		g.Close()
	}()
}
