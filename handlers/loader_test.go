package handlers

import (
	"math/rand"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func newSaramaLoaderMessage() *sarama.ConsumerMessage {
	msg := &sarama.ConsumerMessage{}
	b := make([]rune, 5)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	msg.Value, _ = proto.Marshal(
		&tracepb.Trace{
			Sender: &tracepb.Account{Id: string(b)},
		},
	)
	return msg
}

func TestSaramaLoader(t *testing.T) {
	// Create worker
	w := infra.NewWorker(100)

	// Create Sarama loader
	h := SaramaLoader()
	w.Use(h)

	// Register mock handler
	mockH := NewMockHandler(50)
	w.Use(mockH.Handler())

	// Create a input channel
	in := make(chan interface{})

	// Run worker
	go w.Run(in)

	// Feed sarama channel and then close it
	rounds := 1000
	for i := 1; i <= rounds; i++ {
		in <- newSaramaLoaderMessage()
	}
	close(in)

	// Wait for worker to be done
	<-w.Done()

	if len(mockH.handled) != rounds {
		t.Errorf("Loader: expected %v rounds but got %v", rounds, len(mockH.handled))
	}

	for _, ctx := range mockH.handled {
		if len(ctx.T.Sender().ID) != 5 {
			t.Errorf("Loader: expected Sender ID to have lenght 5 but got %q", ctx.T.Sender().ID)
		}
	}
}
