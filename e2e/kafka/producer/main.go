package main

func main() {}

// import (
// 	// "os"
// 	// "os/signal"
// 	// "syscall"
// 	"fmt"

// 	"github.com/Shopify/sarama"
// 	"github.com/golang/protobuf/proto"
// 	ethpb "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/protobuf/ethereum"
// 	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/protobuf/trace"
// )

// var (
// 	chainURL = "https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7"
// 	kafkaURL = []string{"localhost:9092"}
// 	group    = "test-group"
// 	inTopic  = "test-in"
// 	outTopic = "test-out"
// 	senders  = []string{
// 		"0x664895b5fE3ddf049d2Fb508cfA03923859763C6",
// 		"0xf5956Eb46b377Ae41b41BDa94e6270208d8202bb",
// 		"0x93f7274c9059e601be4512F656B57b830e019E41",
// 		"0xbfc7137876d7Ac275019d70434B0f0779824a969",
// 		"0xA8d8DB1d8919665a18212374d623fc7C0dFDa410",
// 	}
// 	// ERC20Address of token contract to target
// 	ERC20Address = "0x6AFE55b2b5CcA4920182a70c71e793A7Bf44a547"
// )

// func newTokenCredit(i int) *sarama.ProducerMessage {
// 	msg := &sarama.ProducerMessage{
// 		Topic:     inTopic,
// 		Partition: -1,
// 	}
// 	b, _ := proto.Marshal(
// 		&tracepb.Trace{
// 			Sender: &tracepb.Account{Address: "0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff"},
// 			Call:   &tracepb.Call{MethodId: "transfer@ERC20", Args: []string{senders[i%len(senders)], "0x400"}},
// 			Transaction: &ethpb.Transaction{
// 				TxData: &ethpb.TxData{
// 					To: ERC20Address,
// 				},
// 			},
// 		},
// 	)
// 	msg.Value = sarama.ByteEncoder(b)
// 	return msg
// }

// func newMessage(i int) *sarama.ProducerMessage {
// 	msg := &sarama.ProducerMessage{
// 		Topic:     inTopic,
// 		Partition: -1,
// 	}
// 	b, _ := proto.Marshal(
// 		&tracepb.Trace{
// 			Sender: &tracepb.Account{Address: senders[i%len(senders)]},
// 			Call:   &tracepb.Call{MethodId: "transfer@ERC20", Args: []string{"0x71a556C033cD4beB023eb2baa734d0e8304CA88a", "0x200"}},
// 			Transaction: &ethpb.Transaction{
// 				TxData: &ethpb.TxData{
// 					To: ERC20Address,
// 				},
// 			},
// 		},
// 	)
// 	msg.Value = sarama.ByteEncoder(b)
// 	return msg
// }

// func main() {
// 	// Init config, specify appropriate version
// 	config := sarama.NewConfig()
// 	config.Version = sarama.V1_0_0_0
// 	config.Consumer.Return.Errors = true

// 	// Create client
// 	client, err := sarama.NewClient(kafkaURL, config)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer func() { client.Close() }()
// 	fmt.Println("Client ready")

// 	// Create producer
// 	p, err := sarama.NewAsyncProducerFromClient(client)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println("Producer ready")
// 	defer p.Close()

// 	rounds := 10
// 	for i := 0; i < rounds; i++ {
// 		// p.Input() <- newTokenCredit(i)
// 		p.Input() <- newMessage(i)
// 	}
// }
