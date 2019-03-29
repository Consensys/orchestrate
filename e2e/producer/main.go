package main

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	abipb "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/abi"
	commonpb "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
	ethpb "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/ethereum"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/trace"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	kafkaURL = []string{"localhost:9092"}
	inTopic  = "topic-tx-crafter"
	senders  = []string{
		"0xd71400daD07d70C976D6AAFC241aF1EA183a7236",
		//"0xf5956Eb46b377Ae41b41BDa94e6270208d8202bb",
		//"0x93f7274c9059e601be4512F656B57b830e019E41",
		//"0xbfc7137876d7Ac275019d70434B0f0779824a969",
		//"0xA8d8DB1d8919665a18212374d623fc7C0dFDa410",
	}
	// ERC1400Address of token contract to target
	ERC1400Address = "0x8f371DAA8A5325f53b754A7017Ac3803382bc847"
)

func newMessage(i int) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:     inTopic,
		Partition: -1,
	}

	var call *commonpb.Call

	switch i % 2 {
	case 0:
		bytecode := hexutil.MustDecode(bytecodeHex)
		call = &commonpb.Call{
			Contract: &abipb.Contract{Name: "ERC1400", Bytecode: bytecode},
			Method: &abipb.Method{Name: "constructor"},
			Args:   []string{"0xabcd", "0xabcd", "0x10", "[0xcd626bc764e1d553e0d75a42f5c4156b91a63f23,0xcd626bc764e1d553e0d75a42f5c4156b91a63f23]", "0xcd626bc764e1d553e0d75a42f5c4156b91a63f23", "0xabcd"},
		}

	case 1:
		call = &commonpb.Call{
			Contract: &abipb.Contract{Name: "ERC1400"},
			Method: &abipb.Method{Name: "setDocument"},
			Args:   []string{"0xabcd", "0xabcd", "0xabcd"},
		}
	}

	b, _ := proto.Marshal(
		&tracepb.Trace{
			Chain:  &commonpb.Chain{Id: "0x3"},
			Sender: &commonpb.Account{Addr: senders[i%len(senders)]},
			Call: call,
			Tx: &ethpb.Transaction{
				TxData: &ethpb.TxData{
					To: ERC1400Address,
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

	rounds := 50
	for i := 0; i < rounds; i++ {
		p.Input() <- newMessage(i)
	}
}
