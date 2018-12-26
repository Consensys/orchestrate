package main

import (
	// "os"
	// "os/signal"
	// "syscall"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
)

func newMessage() *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:     "test",
		Partition: -1,
	}
	b, _ := proto.Marshal(
		&tracepb.Trace{
			Sender: &tracepb.Account{Id: "sender-id", Address: "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4"},
			Call:   &tracepb.Call{MethodId: "method-id", Args: []string{"0x71a556C033cD4beB023eb2baa734d0e8304CA88a", "0x2386f26fc10000"}},
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
	client, err := sarama.NewClient([]string{"localhost:9092"}, config)
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

	msg := newMessage()
	p.Input() <- msg
}
