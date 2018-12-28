package infra

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	"github.com/golang/protobuf/proto"
	ethpb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/ethereum"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
)

func testPbToProducerMessage(pb *tracepb.Trace) *sarama.ProducerMessage {
	msg := sarama.ProducerMessage{}
	b, err := proto.Marshal(pb)
	if err != nil {
		return nil
	}
	msg.Value = sarama.ByteEncoder(b)
	return &msg
}

func newSaramaProducerTestMessage() *tracepb.Trace {
	var pb tracepb.Trace
	pb.Chain = &tracepb.Chain{Id: "0x1", IsEIP155: true}
	pb.Sender = &tracepb.Account{Id: "", Address: "0xdbb881a51CD4023E4400CEF3ef73046743f08da3"}
	pb.Receiver = &tracepb.Account{Id: "toto", Address: "0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff"}
	pb.Transaction = &ethpb.Transaction{
		TxData: &ethpb.TxData{
			Nonce:    10,
			To:       "0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff",
			Value:    "0xa2bfe3",
			Data:     "0xbabe",
			GasPrice: "0x0",
			Gas:      1234,
		},
		Raw:  "0xbeef",
		Hash: "0x0000000000000000000000000000000000000000000000000000000000000000",
	}
	pb.Call = &tracepb.Call{
		MethodId: "abcde",
		Args:     []string{"0xfF778b716FC07D98839f48DdB88D8bE583BEB684", "0x2386f26fc10000"},
	}
	pb.Errors = []*tracepb.Error{}
	return &pb
}

var expected, _ = proto.Marshal(newSaramaProducerTestMessage())
var counter int64

func checker(val []byte) error {
	atomic.AddInt64(&counter, 1)
	if string(val) != string(expected) {
		return fmt.Errorf("Expected %q but got %q", string(expected), string(val))
	}
	return nil
}

func TestSaramaProducerConcurrent(t *testing.T) {
	mp := mocks.NewSyncProducer(t, nil)
	p := NewSaramaProducer(mp, testPbToProducerMessage)

	rounds := 1000
	for i := 1; i <= rounds; i++ {
		mp.ExpectSendMessageWithCheckerFunctionAndSucceed(checker)
	}

	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			p.Produce(newSaramaProducerTestMessage())
		}(i)
	}
	wg.Wait()

	if counter != int64(rounds) {
		t.Errorf("SaramaProducer: expected %v rounds but got %v", rounds, counter)
	}
}
