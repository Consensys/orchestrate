package engine

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type TestHandler struct {
	mux     *sync.Mutex
	handled []*TxContext
}

// StringMsg is a dummy engine.StringMsg implementation
type StringMsg string

// Entrypoint is a dummy implementation of the method "Entrypoint of the dummy engine"
func (s StringMsg) Entrypoint() string { return "" }
func (s StringMsg) Header() Header     { return nil }
func (s StringMsg) Key() []byte        { return nil }
func (s StringMsg) Value() []byte      { return nil }

func (h *TestHandler) Handler(t *testing.T) HandlerFunc {
	return func(txctx *TxContext) {
		// We add some randomness in time execution
		r := rand.Intn(100)
		time.Sleep(time.Duration(r) * time.Millisecond)
		h.mux.Lock()
		defer h.mux.Unlock()
		h.handled = append(h.handled, txctx)
	}
}

func TestEngine(t *testing.T) {
	h := TestHandler{
		mux:     &sync.Mutex{},
		handled: []*TxContext{},
	}

	// Create new Engine and register test handler
	eng := NewEngine(&Config{Slots: 100})
	eng.Register(h.Handler(t))

	// Create input channels and prefills it
	ins := make([]chan Msg, 0)
	for i := 0; i < 50; i++ {
		in := make(chan Msg, 20)
		for j := 0; j < 20; j++ {
			s := StringMsg(fmt.Sprintf("test-%v-%v", i, j))
			in <- &s
		}
		close(in)
		ins = append(ins, in)
	}

	// Start consuming every input channel
	wg := &sync.WaitGroup{}
	for i := range ins {
		wg.Add(1)
		go func(in <-chan Msg) {
			eng.Run(context.Background(), in)
			wg.Done()
		}(ins[i])
	}

	// Wait for engine to finish consuming
	wg.Wait()

	assert.Len(t, h.handled, 1000, "All messages should have been processed")
}

func TestEngineStopped(t *testing.T) {
	h := TestHandler{
		mux:     &sync.Mutex{},
		handled: []*TxContext{},
	}

	// Create new Engine and register test handler
	eng := NewEngine(&Config{Slots: 100})
	eng.Register(h.Handler(t))

	// Create input channels and prefills it
	ins := make([]chan Msg, 0)
	for i := 0; i < 50; i++ {
		in := make(chan Msg, 20)
		for j := 0; j < 20; j++ {
			in <- StringMsg(fmt.Sprintf("test-%v-%v", i, j))
		}
		close(in)
		ins = append(ins, in)
	}

	// Start consuming every input channel
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	for i := range ins {
		wg.Add(1)
		go func(in <-chan Msg) {
			eng.Run(ctx, in)
			wg.Done()
		}(ins[i])
	}

	// Sleep for a short time and interrupt
	time.Sleep(100 * time.Millisecond)
	cancel()

	// Wait for Engine to consume messages
	wg.Wait()

	assert.True(t, len(h.handled) < 500, "Expected at least half of the message not to have been consumed")

	// We drain and count all messages that have not been consumed
	count := 0
	for i := range ins {
		for range ins[i] {
			count++
		}
	}

	assert.Equal(t, 1000, len(h.handled)+count, "Expected all message to have either been consumed or still be in input channel")
}

func testSleepingHandler(txctx *TxContext) {
	time.Sleep(txctx.Get("duration").(time.Duration))
}

func makeTimeoutContext(i int) *TxContext {
	txctx := NewTxContext()
	txctx.Reset()
	txctx.Prepare(log.NewEntry(log.StandardLogger()), nil)

	switch i % 2 {
	case 0:
		txctx.Set("duration", 50*time.Millisecond)
		txctx.Set("errors", 0)
	case 1:
		txctx.Set("duration", 100*time.Millisecond)
		txctx.Set("errors", 1)
	}
	return txctx
}

func TestTimeoutHandler(t *testing.T) {
	timeoutHandler := TimeoutHandler(testSleepingHandler, 80*time.Millisecond, "Test timeout")

	rounds := 100
	outs := make(chan *TxContext, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		txctx := makeTimeoutContext(i)
		go func(txctx *TxContext) {
			defer wg.Done()
			timeoutHandler(txctx)
			outs <- txctx
		}(txctx)
	}
	wg.Wait()
	close(outs)

	assert.Len(t, outs, rounds, "Timeout: processed contexts count should be correct")

	for out := range outs {
		errCount := out.Get("errors").(int)
		assert.Len(t, out.Envelope.Errors, errCount, "Timeout: expected correct count of errors")
	}
}
