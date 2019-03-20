package infra

import (
	"fmt"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	infSarama "gitlab.com/ConsenSys/client/fr/core-stack/infra/sarama.git"
	trace "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/trace"
)

func initSarama(infra *Infra, wait *sync.WaitGroup) {
	initClient(infra, wait)
	initProducer(infra, wait)
	initUnmarshaller(infra, wait)
	wait.Done()
}

func initClient(infra *Infra, wait *sync.WaitGroup) {
	// Init config
	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	config.Consumer.Return.Errors = true
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true

	// Create sarama client
	client, err := sarama.NewClient(viper.GetStringSlice("kafka.addresses"), config)
	if err != nil {
		log.WithError(err).Fatalf("infra-sarama: could not to start client")
		return
	}

	// Retrieve and log connectted brokers
	var brokers = make(map[int32]string)
	for _, v := range client.Brokers() {
		brokers[v.ID()] = v.Addr()
	}
	log.Infof("infra-sarama: client ready (connected to brokers: %v)", brokers)

	// Attach client
	infra.SaramaClient = client

	// Close when infra is cut
	go func() {
		<-infra.ctx.Done()
		client.Close()
	}()
}

func initProducer(infra *Infra, wait *sync.WaitGroup) {
	// Create sarama sync producer
	p, err := sarama.NewSyncProducerFromClient(infra.SaramaClient)
	if err != nil {
		log.WithError(err).Fatalf("infra-sarama: could not start producer")
	}
	log.Debug("infra-sarama: producer ready")

	// Initialize
	marshaller := infSarama.NewMarshaller()
	prepareMsg := func(t *trace.Trace, msg *sarama.ProducerMessage) error {
		err := marshaller.Marshal(t, msg)
		if err != nil {
			return err
		}

		// Set topic
		msg.Topic = fmt.Sprintf("%v-%v", viper.GetString("worker.out"), t.GetChain().GetId())

		return nil
	}

	// Attach producer
	infra.SaramaProducer = p
	infra.Producer = infSarama.NewProducer(
		p,
		prepareMsg,
	)

	// Close when infra is cut
	go func() {
		<-infra.ctx.Done()
		p.Close()
	}()
}

func initUnmarshaller(infra *Infra, wait *sync.WaitGroup) {
	infra.Unmarshaller = infSarama.NewUnmarshaller()
}

// GetLastRecord retrieve the last record that has been produced on a given topic/partition
func GetLastRecord(client sarama.Client, topic string, partition int32) (*sarama.Record, error) {
	// Retrieve last offset that has been produced for topic-partition
	lastOffset, err := client.GetOffset(topic, partition, -1)
	if err != nil {
		return nil, err
	}

	// Get broker Leader fo topic-partition
	broker, err := client.Leader(topic, partition)
	if err != nil {
		return nil, err
	}

	// Fetch block containing last produced record on topic partition
	req := &sarama.FetchRequest{
		MinBytes:    client.Config().Consumer.Fetch.Min,
		MaxWaitTime: int32(client.Config().Consumer.MaxWaitTime / time.Millisecond),
	}
	req.AddBlock(topic, 0, lastOffset-1, client.Config().Consumer.Fetch.Max)
	req.Version = 4
	req.Isolation = sarama.ReadUncommitted
	response, err := broker.Fetch(req)

	if err != nil {
		return nil, err
	}

	// Parse block to retrieve record of interest
	block := response.GetBlock(topic, partition)
	if len(block.RecordsSet) == 0 {
		return nil, nil
	}
	records := block.RecordsSet[0]
	record := records.RecordBatch.Records[0]
	return record, nil
}
