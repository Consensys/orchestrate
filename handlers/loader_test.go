package handlers

import (
	"math/rand"
	"sync"
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

func newLoaderMessage() *sarama.ConsumerMessage {
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

type testHandler struct {
	mux     *sync.Mutex
	handled []*infra.Context
}

func (h *testHandler) Handler(maxtime int, t *testing.T) infra.HandlerFunc {
	return func(ctx *infra.Context) {
		// We add some randomness in time execution
		r := rand.Intn(maxtime)
		time.Sleep(time.Duration(r) * time.Millisecond)
		h.mux.Lock()
		defer h.mux.Unlock()
		h.handled = append(h.handled, ctx)
	}
}

func TestSaramaLoader(t *testing.T) {
	testH := &testHandler{
		mux:     &sync.Mutex{},
		handled: []*infra.Context{},
	}
	h := SaramaLoader()

	w := infra.NewWorker([]infra.HandlerFunc{h, testH.Handler(50, t)}, 100)

	// Create a input channel
	in := make(chan interface{})

	// Run worker
	go w.Run(in)

	// Feed sarama channel and then close it
	rounds := 1000
	for i := 1; i <= rounds; i++ {
		in <- newLoaderMessage()
	}
	close(in)

	// Wait for worker to be done
	<-w.Done()

	if len(testH.handled) != rounds  {
		t.Errorf("Loader: expected %v rounds but got %v", rounds, len(testH.handled))
	}

	for _, ctx := range testH.handled {
		if len(ctx.T.Sender().ID) != 5 {
			t.Errorf("Loader: expected Sender ID to have lenght 5 but got %q", ctx.T.Sender().ID)
		}
	}
}
