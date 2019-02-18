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
	topic    = "topic-tx-decoder-3"
)

func newMessage(i int) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
	}
	b, _ := proto.Marshal(
		&tracepb.Trace{
			Chain: &tracepb.Chain{Id: "0x3"},
			Receipt: &ethpb.Receipt{
				TxHash:          "0xbf0b3048242aff8287d1dd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775afd",
				BlockHash:       "0x",
				BlockNumber:     uint64(0),
				TxIndex:         uint64(0),
				ContractAddress: "0x75d2917bD1E6C7c94d24dFd11C8EeAeFd3003C85",
				PostState:       "0x",
				Status:          uint64(0),
				Bloom:           "0x",
				Logs: []*ethpb.Log{
					&ethpb.Log{
						Address: "0x75d2917bD1E6C7c94d24dFd11C8EeAeFd3003C85",
						Topics: []string{
							"0xe8f0a47da72ca43153c7a5693a827aa8456f52633de9870a736e5605bff4af6d",
							"0x000000000000000000000000d71400dad07d70c976d6aafc241af1ea183a7236",
							"0x000000000000000000000000d71400dad07d70c976d6aafc241af1ea183a7236",
							"0x000000000000000000000000b5747835141b46f7c472393b31f8f5a57f74a44f",
						},
						Data:        "0x0000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000061000000000000000000000000000000000000000000000000000000006f0c7f50cd4b7e4466b726279b1506bc89d8e74ab9268a255eeb1c78f163d51a83c7380d54a8b597ee26351c15c83f922fd6b37334970d3f832e5e11e36acbecb460ffdb01000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
						BlockNumber: uint64(10),
						TxHash:      "0xbf0b3048242aff8287d1dd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775afd",
						TxIndex:     uint64(10),
						BlockHash:   "0xea2460a53299f7201d82483d891b26365ff2f49cd9c5c0c7686fd75599fda5b2",
					},
				},
				GasUsed:           uint64(10000),
				CumulativeGasUsed: uint64(10000),
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
