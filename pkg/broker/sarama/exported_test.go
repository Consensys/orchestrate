// +build unit

package sarama

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	cgmock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine"
)

type ExportedTestSuite struct {
	suite.Suite
}

func (s *ExportedTestSuite) SetupTest() {
	SetGlobalConfig(nil)
	SetGlobalClient(nil)
	SetGlobalSyncProducer(nil)
	SetGlobalConsumerGroup(nil)
	initClientOnce = &sync.Once{}
	initProducerOnce = &sync.Once{}
	initConsumerGroupOnce = &sync.Once{}
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ExportedTestSuite))
}

func (s *ExportedTestSuite) TestInitConfig() {
	viper.Set(kafkaTLSEnableViperKey, true)
	viper.Set(kafkaTLSClientCertFilePathViperKey, "testdata/example-cert.pem")
	viper.Set(kafkaTLSClientKeyFilePathViperKey, "testdata/example-key.pem")
	viper.Set(kafkaTLSCACertFilePathViperKey, "testdata/example-ca-cert.pem")

	InitConfig()
	assert.NotNil(s.T(), GlobalConfig(), "Config should have been set")
}

func (s *ExportedTestSuite) TestNewTLSConfig() {

	tlsConfig, err := NewTLSConfig(
		"testdata/example-cert.pem",
		"testdata/example-key.pem",
		"testdata/example-ca-cert.pem",
	)
	assert.NoError(s.T(), err, "TLS should be instantiated without error")
	assert.NotNil(s.T(), tlsConfig, "TLS should be instantiated without error")

	_, err = NewTLSConfig(
		"testdata/example-cert-error.pem",
		"testdata/example-key.pem",
		"testdata/example-ca-cert.pem",
	)
	assert.Error(s.T(), err, "TLS should not be instantiated with a wrong file path")

	_, err = NewTLSConfig(
		"testdata/example-cert.pem",
		"testdata/example-key.pem",
		"testdata/example-ca-cert-error.pem",
	)
	assert.Error(s.T(), err, "TLS should not be instantiated with a wrong file path")

}

// func (s *ExportedTestSuite) TestInitClient() {
// 
// 	seedBroker := sarama.NewMockBroker(s.T(), 1)
// 	seedBroker.Returns(new(sarama.MetadataResponse))
// 
// 	config := sarama.NewConfig()
// 	SetGlobalConfig(config)
// 
// 	viper.Set(KafkaURLViperKey, []string{
// 		seedBroker.Addr(),
// 	})
// 	_ = InitClient(context.Background())
// 	assert.NotNil(s.T(), GlobalClient(), "Client should have been set")
// }

// func (s *ExportedTestSuite) TestInitSyncProducer() {
// 
// 	seedBroker := sarama.NewMockBroker(s.T(), 1)
// 	seedBroker.Returns(new(sarama.MetadataResponse))
// 
// 	config := sarama.NewConfig()
// 	config.Producer.Return.Successes = true
// 	SetGlobalConfig(config)
// 
// 	viper.Set(KafkaURLViperKey, []string{
// 		seedBroker.Addr(),
// 	})
// 
// 	_ = InitClient(context.Background())
// 
// 	InitSyncProducer(context.Background())
// 
// 	assert.NotNil(s.T(), GlobalSyncProducer(), "Producer should have been set")
// }

func (s *ExportedTestSuite) TestSetSyncProducer() {

	producer := &mocks.SyncProducer{}

	SetGlobalSyncProducer(producer)
	InitSyncProducer(context.Background())

	assert.NotNil(s.T(), GlobalSyncProducer(), "Producer should have been set")
}

func (s *ExportedTestSuite) TestSetConsumerGroup() {

	msgs := make(map[string]map[int32][]*sarama.ConsumerMessage)
	cg := cgmock.NewConsumerGroup("test-group", msgs)

	SetGlobalConsumerGroup(cg)

	InitConsumerGroup(context.Background(), "group-name")
	assert.NotNil(s.T(), GlobalConsumerGroup(), "ConsumerGroup should have been set")
}

func TestConsume(t *testing.T) {
	conf := engine.NewConfig()
	e := engine.NewEngine(&conf)

	counter := CounterHandler{}
	e.Register(counter.Handle)

	cgHandler := NewEngineConsumerGroupHandler(e)

	// Init messages in topics
	topics := []string{"test-topic-1", "test-topic-2"}
	msgs, count := MockConsumerMessages(topics, 3, 10)

	// Create consumer group
	cg := cgmock.NewConsumerGroup("test-group", msgs)
	SetGlobalConsumerGroup(cg)
	ctx, cancel := context.WithCancel(context.Background())

	// Run Consumer
	wg := sync.WaitGroup{}
	wg.Add(1)
	var err error
	go func() {
		err = Consume(ctx, topics, cgHandler)
		wg.Done()
	}()
	time.Sleep(100 * time.Millisecond)
	cancel()

	wg.Wait()

	assert.NoError(t, err, "No error expected")

	assert.Equal(t, int32(count), counter.counter, "Count of processed message should be correct")
}

// MockConsumerMessages Creates "nbTopics" topics, "nbPartition" partitions, and "nbMessageByPartition" messages by partition
func MockConsumerMessages(topics []string, nbPartition, nbMessageByPartition int) (msgs map[string]map[int32][]*sarama.ConsumerMessage, counter int) {

	m := make(map[string]map[int32][]*sarama.ConsumerMessage)
	count := 0

	for _, topic := range topics {
		m[topic] = make(map[int32][]*sarama.ConsumerMessage)
		for partition := range make([]int32, nbPartition) {
			m[topic][int32(partition)] = []*sarama.ConsumerMessage{}
			for i := range make([]int, nbMessageByPartition) {
				m[topic][int32(partition)] = append(m[topic][int32(partition)], &sarama.ConsumerMessage{
					Topic:     topic,
					Partition: int32(partition),
					Offset:    int64(i),
				})
				count++
			}
		}
	}
	return m, count
}
