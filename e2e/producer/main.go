package main

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	ethpb "gitlab.com/ConsenSys/client/fr/core-stack/core.git/protobuf/ethereum"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core.git/protobuf/trace"
)

var (
	kafkaURL = []string{"localhost:9092"}
	topic    = "topic-tx-signer"
	senders  = []string{
		"0xd71400daD07d70C976D6AAFC241aF1EA183a7236",
		"0xf5956Eb46b377Ae41b41BDa94e6270208d8202bb",
		"0x93f7274c9059e601be4512F656B57b830e019E41",
		"0xbfc7137876d7Ac275019d70434B0f0779824a969",
		"0xA8d8DB1d8919665a18212374d623fc7C0dFDa410",
	}
)

func newMessage(i int) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
	}
	b, _ := proto.Marshal(
		&tracepb.Trace{
			Chain:  &tracepb.Chain{Id: "0x3"},
			Sender: &tracepb.Account{Address: senders[i%len(senders)]},
			Transaction: &ethpb.Transaction{
				TxData: &ethpb.TxData{
					Nonce:    1,
					To:       "0xfF778b716FC07D98839f48DdB88D8bE583BEB684",
					Value:    "0x2386f26fc10000",
					Gas:      21136,
					GasPrice: "0xee6b2800",
					Data:     "0xabcd",
				},
			},
		},
	)
	msg.Value = sarama.ByteEncoder(b)
	return msg
}

func main() {
	// Init config, specify appropriate version
	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	config.Consumer.Return.Errors = true

	// Create client

	client, err := sarama.NewClient(kafkaURL, config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() { client.Close() }()
	fmt.Println("Client ready")

	// Create producer
	p, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Producer ready")
	defer p.Close()

	rounds := 10
	for i := 0; i < rounds; i++ {
		p.Input() <- newMessage(i)
	}
}
