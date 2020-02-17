package main

import (
	"fmt"
	"math/big"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/tx"

	"github.com/ethereum/go-ethereum/common"
	uuid "github.com/satori/go.uuid"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

var (
	kafkaURL = []string{"localhost:9092"}
)

func newMessage(i int) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:     "topic-tx-crafter",
		Partition: -1,
	}

	envelope := tx.NewEnvelope().
		SetID(uuid.NewV4().String()).
		SetNonce(uint64(i)).
		SetTo(common.HexToAddress("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C")).
		SetFrom(common.HexToAddress("0xdbb881a51cd4023e4400cef3ef73046743f08da3")).
		SetValue(big.NewInt(100000)).
		SetGas(21000).
		SetChainName("geth")

	e, _ := proto.Marshal(envelope.TxRequest())
	msg.Value = sarama.ByteEncoder(e)
	return msg
}

func main() {
	// Init config, specify appropriate version
	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true

	// Create client

	log.Info("Connecting to Kafka: ", kafkaURL)
	client, err := sarama.NewClient(kafkaURL, config)
	if err != nil {
		log.Info(err)
		return
	}
	defer func() {
		log.Info("Closing a client")
		e := client.Close()
		if e != nil {
			log.Info("Error while closing a client")
		}
	}()
	log.Info("Client ready")

	// Create producer
	p, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		log.Info(err)
		return
	}
	log.Info("Producer ready")
	defer func() {
		log.Info("Closing a producer")
		e := p.Close()
		if e != nil {
			log.Info("Error while closing a producer: ", e)
		}
	}()

	rounds := 10
	for i := 0; i < rounds; i++ {
		p.Input() <- newMessage(i)
	}

	for i := 0; i < rounds; i++ {
		select {
		case success := <-p.Successes():
			log.Info("Success: ", success.Topic, success.Key)
		case err := <-p.Errors():
			fmt.Println("Error", err)
		}
	}

}
