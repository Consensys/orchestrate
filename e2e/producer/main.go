package main

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
)

/*
	As the tag 0.3.0, the available public keys are:
	[
		0x93f7274c9059e601be4512F656B57b830e019E41
		0x7E654d251Da770A068413677967F6d3Ea2FeA9E4
		0xdbb881a51CD4023E4400CEF3ef73046743f08da3
		0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff
		0xA8d8DB1d8919665a18212374d623fc7C0dFDa410
		0xffbBa394DEf3Ff1df0941c6429887107f58d4e9b
		0x664895b5fE3ddf049d2Fb508cfA03923859763C6
		0xfF778b716FC07D98839f48DdB88D8bE583BEB684
		0xf5956Eb46b377Ae41b41BDa94e6270208d8202bb
		0xbfc7137876d7Ac275019d70434B0f0779824a969
	]
*/
var (
	kafkaURL = []string{"localhost:9092"}
	topic    = "topic-tx-signer"
	senders  = []string{
		"0xd71400daD07d70C976D6AAFC241aF1EA183a7236", // As of 0.3.0, this address is not stored by default
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
		&envelope.Envelope{
			Chain:  &common.Chain{Id: "3"},
			Sender: &common.Account{Addr: senders[i%len(senders)]},
			Tx: &ethereum.Transaction{
				TxData: &ethereum.TxData{
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
