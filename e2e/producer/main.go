package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	commonpb "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
	ethpb "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/ethereum"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/trace"
)

var (
	kafkaURL = []string{"localhost:9092"}
)

func newMessage(i int) *sarama.ProducerMessage {
	var topic, chainID string
	switch i % 4 {
	case 0:
		topic = "topic-tx-decoder-1"
		chainID = "0x1"
	case 1:
		topic = "topic-tx-decoder-3"
		chainID = "0x3"
	case 2:
		topic = "topic-tx-decoder-4"
		chainID = "0x4"
	case 3:
		topic = "topic-tx-decoder-2a"
		chainID = "0x2a"
	}
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
	}
	Data, _ := hexutil.Decode("0x0000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000061000000000000000000000000000000000000000000000000000000006f0c7f50cd4b7e4466b726279b1506bc89d8e74ab9268a255eeb1c78f163d51a83c7380d54a8b597ee26351c15c83f922fd6b37334970d3f832e5e11e36acbecb460ffdb01000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	b, _ := proto.Marshal(
		&tracepb.Trace{
			Chain: &commonpb.Chain{Id: chainID},
			Receipt: &ethpb.Receipt{
				TxHash:          "0xbf0b3048242aff8287d1dd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775afd",
				BlockHash:       "0x",
				BlockNumber:     uint64(0),
				TxIndex:         uint64(0),
				ContractAddress: "0x75d2917bD1E6C7c94d24dFd11C8EeAeFd3003C85",
				PostState:       []byte{0x3b, 0x19, 0x8b, 0xfd, 0x5d, 0x29, 0x07, 0x28, 0x5a, 0xf0, 0x09, 0xe9, 0xae, 0x84, 0xa0, 0xec, 0xd6, 0x36, 0x77, 0x11, 0x0d, 0x89, 0xd7, 0xe0, 0x30, 0x25, 0x1a, 0xcb, 0x87, 0xf6, 0x48, 0x7e},
				Status:          uint64(0),
				Bloom:           []byte{0x1a, 0x05, 0x56, 0x90, 0xd9, 0xdb, 0x80, 0x00},
				Logs: []*ethpb.Log{
					&ethpb.Log{
						Address: "0x75d2917bD1E6C7c94d24dFd11C8EeAeFd3003C85",
						Topics: []string{
							"0xe8f0a47da72ca43153c7a5693a827aa8456f52633de9870a736e5605bff4af6d",
							"0x000000000000000000000000d71400dad07d70c976d6aafc241af1ea183a7236",
							"0x000000000000000000000000d71400dad07d70c976d6aafc241af1ea183a7236",
							"0x000000000000000000000000b5747835141b46f7c472393b31f8f5a57f74a44f",
						},
						Data:        Data,
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
