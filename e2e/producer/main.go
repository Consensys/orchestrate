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
	topic    = "topic-tx-sender"
)

func newMessage(i int) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
	}
	b, _ := proto.Marshal(
		&tracepb.Trace{
			Chain: &tracepb.Chain{Id: "0x3"},
			Transaction: &ethpb.Transaction{
				TxData: &ethpb.TxData{Nonce: 1, To: "0xfF778b716FC07D98839f48DdB88D8bE583BEB684", Value: "0x2386f26fc10000", Gas: 21136, GasPrice: "0xee6b2800", Data: "0xabcd"},
				Raw:    "0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80",
				Hash:   "0xbf0b3048242aff8287d1dd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775afd",
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
