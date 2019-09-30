package main

import (
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/types/chain"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/types/ethereum"
)

var (
	kafkaURL = []string{"localhost:9092"}
)

func newMessage(i int) *sarama.ProducerMessage {
	var topic string
	var chainID int64
	switch i % 4 {
	case 0:
		topic = "topic-tx-decoder-1"
		chainID = 1
	case 1:
		topic = "topic-tx-decoder-3"
		chainID = 3
	case 2:
		topic = "topic-tx-decoder-4"
		chainID = 4
	case 3:
		topic = "topic-tx-decoder-42"
		chainID = 42
	}
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
	}
	b, _ := proto.Marshal(
		&envelope.Envelope{
			Chain: chain.FromInt(chainID),
			Receipt: &ethereum.Receipt{
				TxHash:          ethereum.HexToHash("0xbf0b3048242aff8287d1dd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775afd"),
				BlockHash:       nil,
				BlockNumber:     uint64(0),
				TxIndex:         uint64(0),
				ContractAddress: ethereum.HexToAccount("0x75d2917bD1E6C7c94d24dFd11C8EeAeFd3003C85"),
				PostState:       hexutil.MustDecode("0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e"),
				Status:          uint64(0),
				Bloom:           hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001a055690d9db80000"),
				Logs: []*ethereum.Log{
					{
						Address: ethereum.HexToAccount("0x75d2917bD1E6C7c94d24dFd11C8EeAeFd3003C85"),
						Topics: []*ethereum.Hash{
							ethereum.HexToHash("0xe8f0a47da72ca43153c7a5693a827aa8456f52633de9870a736e5605bff4af6d"),
							ethereum.HexToHash("0x000000000000000000000000d71400dad07d70c976d6aafc241af1ea183a7236"),
							ethereum.HexToHash("0x000000000000000000000000d71400dad07d70c976d6aafc241af1ea183a7236"),
							ethereum.HexToHash("0x000000000000000000000000b5747835141b46f7c472393b31f8f5a57f74a44f"),
						},
						Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000061000000000000000000000000000000000000000000000000000000006f0c7f50cd4b7e4466b726279b1506bc89d8e74ab9268a255eeb1c78f163d51a83c7380d54a8b597ee26351c15c83f922fd6b37334970d3f832e5e11e36acbecb460ffdb01000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
						BlockNumber: uint64(10),
						TxHash:      ethereum.HexToHash("0xbf0b3048242aff8287d1dd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775afd"),
						TxIndex:     uint64(10),
						BlockHash:   ethereum.HexToHash("0xea2460a53299f7201d82483d891b26365ff2f49cd9c5c0c7686fd75599fda5b2"),
					},
				},
				GasUsed:           uint64(10000),
				CumulativeGasUsed: uint64(10000),
			},
			Metadata: &envelope.Metadata{
				Id: "bc430299-4d34-4920-81be-b9fc5b6ef0f9",
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
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true

	// Create client

	fmt.Println("Connecting to Kafka: ", kafkaURL)
	client, err := sarama.NewClient(kafkaURL, config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		fmt.Println("Closing a client")
		e := client.Close()
		if e != nil {
			fmt.Println("Error while closing a client")
		}
	}()
	fmt.Println("Client ready")

	// Create producer
	p, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Producer ready")
	defer func() {
		fmt.Println("Closing a producer")
		e := p.Close()
		if e != nil {
			fmt.Println("Error while closing a producer: ", e)
		}
	}()

	rounds := 10
	for i := 0; i < rounds; i++ {
		p.Input() <- newMessage(i)
	}

	for i := 0; i < rounds; i++ {
		select {
		case success := <-p.Successes():
			fmt.Println("Success: ", success.Topic, success.Key)
		case err := <-p.Errors():
			fmt.Println("Error", err)
		}
	}

}
