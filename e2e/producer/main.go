package main

import (
	"fmt"
	"math/rand"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/chain"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
)

var (
	kafkaURL = []string{"localhost:9092"}
	topic    = "topic-tx-sender"
)

var letterRunes = []rune("abcdef0123456789")

// RandString creates a random string
func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func newMessage() *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
	}
	b, _ := proto.Marshal(
		&envelope.Envelope{
			Chain: chain.CreateChainInt(888),
			Tx: &ethereum.Transaction{
				TxData: &ethereum.TxData{
					Nonce:    1,
					To:       ethereum.HexToAccount("0xfF778b716FC07D98839f48DdB88D8bE583BEB684"),
					Value:    ethereum.HexToQuantity("0x2386f26fc10000"),
					Gas:      21136,
					GasPrice: ethereum.HexToQuantity("0xee6b2800"),
					Data:     ethereum.HexToData("0xabcd"),
				},
				Raw:  ethereum.HexToData("0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80"),
				Hash: ethereum.HexToHash("0x" + RandString(64)),
			},
			Metadata: &envelope.Metadata{
				Id: RandString(32),
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
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	// Create client
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
		p.Input() <- newMessage()
	}

	for i := 0; i < rounds; i++ {
		select {
		case success := <-p.Successes():
			fmt.Println("Success", success)
		case err := <-p.Errors():
			fmt.Println("Error", err)
		}
	}
}
