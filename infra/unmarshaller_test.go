package infra

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core.git/protobuf/trace"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func newProtoMessage() *tracepb.Trace {
	b := make([]rune, 5)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return &tracepb.Trace{
		Sender: &tracepb.Account{Id: string(b)},
	}
}

func newSaramaMessage() *sarama.ConsumerMessage {
	msg := sarama.ConsumerMessage{}
	msg.Value, _ = proto.Marshal(
		newProtoMessage(),
	)
	return &msg
}

func TestTraceProtoUnmarshallerConcurrent(t *testing.T) {
	u := TraceProtoUnmarshaller{}
	pbs := make([]*tracepb.Trace, 0)
	rounds := 1000
	wg := &sync.WaitGroup{}
	for i := 1; i < rounds; i++ {
		pb := &tracepb.Trace{}
		pbs = append(pbs, pb)
		wg.Add(1)
		go func() {
			defer wg.Done()
			u.Unmarshal(newProtoMessage(), pb)
		}()
	}
	wg.Wait()

	for _, pb := range pbs {
		if len(pb.GetSender().GetId()) != 5 {
			t.Errorf("TraceProtoUnmarshaller: expected a 5 long string but got %q", pb.GetSender().GetId())
		}
	}
}

func TestSaramaUnmarshallerConcurrent(t *testing.T) {
	u := SaramaUnmarshaller{}
	pbs := make([]*tracepb.Trace, 0)
	rounds := 1000
	wg := &sync.WaitGroup{}
	for i := 1; i < rounds; i++ {
		pb := &tracepb.Trace{}
		pbs = append(pbs, pb)
		wg.Add(1)
		go func() {
			defer wg.Done()
			u.Unmarshal(newSaramaMessage(), pb)
		}()
	}
	wg.Wait()

	for _, pb := range pbs {
		if len(pb.GetSender().GetId()) != 5 {
			t.Errorf("SaramaUnmarshaller: expected a 5 long string but got %q", pb.GetSender().GetId())
		}
	}

}
