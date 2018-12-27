package handlers

import (
	"fmt"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	"github.com/golang/protobuf/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf"
	ethpb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/ethereum"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
)

func testCtxToProducerMessage(ctx *infra.Context) *sarama.ProducerMessage {
	msg := sarama.ProducerMessage{}
	ctx.Pb.Reset()
	protobuf.DumpTrace(ctx.T, ctx.Pb)
	b, _ := proto.Marshal(ctx.Pb)
	msg.Value = sarama.ByteEncoder(b)
	return &msg
}

func newProducerTestMessage() *tracepb.Trace {
	var pb tracepb.Trace
	pb.Chain = &tracepb.Chain{Id: "0x1", IsEIP155: true}
	pb.Sender = &tracepb.Account{Id: "", Address: "0xdbb881a51CD4023E4400CEF3ef73046743f08da3"}
	pb.Receiver = &tracepb.Account{Id: "toto", Address: "0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff"}
	pb.Transaction = &ethpb.Transaction{
		TxData: &ethpb.TxData{
			Nonce: 10,
			To:    "0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff",
			Value: "0xa2bfe3",
			Data:  "0xbabe",
			GasPrice: "0x0",
			Gas: 1234,
		},
		Raw: "0xbeef",
		Hash: "0x0000000000000000000000000000000000000000000000000000000000000000",
	}
	pb.Call = &tracepb.Call{
		MethodId: "abcde",
		Args: []string{"0xfF778b716FC07D98839f48DdB88D8bE583BEB684", "0x2386f26fc10000"},
	}
	pb.Errors = []*tracepb.Error{}
	return &pb
}

var expected, _ = proto.Marshal(newProducerTestMessage())

func checker(val []byte) error {
	if string(val) != string(expected) {
		return fmt.Errorf("Expected %q but got %q", string(expected), string(val))
	}
	return nil
}

func TestProducer(t *testing.T) {
	// Create worker
	w := infra.NewWorker(100)
	w.Use(Loader(&TraceProtoUnmarshaller{}))

	//Register mock (for time randomness)
	w.Use(NewMockHandler(50).Handler())

	// Register producer
	mp := mocks.NewSyncProducer(t, nil)
	p := Producer(NewSaramaProducer(mp, testCtxToProducerMessage))
	w.Use(p)

	// Register mock (to track output)
	mockH := NewMockHandler(1)
	w.Use(mockH.Handler())

	rounds := 1000
	for i := 1; i <= rounds; i++ {
		mp.ExpectSendMessageWithCheckerFunctionAndSucceed(checker)
	}

	// Create a Sarama message channel
	in := make(chan interface{})

	// Run worker
	go w.Run(in)

	// Feed channel channel and then close it

	for i := 1; i <= rounds; i++ {
		in <- newProducerTestMessage()
	}
	close(in)

	// Wait for worker to be done
	<-w.Done()

	// Run worker
	go w.Run(in)

	if len(mockH.handled) != rounds {
		t.Errorf("Gas: expected %v rounds but got %v", rounds, len(mockH.handled))
	}

}
