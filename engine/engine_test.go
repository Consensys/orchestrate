package engine

import (
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

func (h *TestHandler) Handler(t *testing.T) HandlerFunc {
	return func(ctx *TxContext) {
		// We add some randomness in time execution
		r := rand.Intn(100)
		time.Sleep(time.Duration(r) * time.Millisecond)
		h.mux.Lock()
		defer h.mux.Unlock()
		h.handled = append(h.handled, ctx)
	}
}

func TestEngine(t *testing.T) {
	h := TestHandler{
		mux:     &sync.Mutex{},
		handled: []*TxContext{},
	}

	// Create new Engine and register test handler
	e := NewEngine(
		Config{Slots: 100, Partitions: 100},
	)
	e.Use(h.Handler(t))

	// Create a Sarama message channel
	in := make(chan interface{})

	// Run Engine
	go e.Run(in)

	// Feed sarama channel and then close it
	rounds := 1000
	for i := 1; i <= rounds; i++ {
		in <- "test"
	}
	close(in)

	// Wait for Engine to be done
	<-e.Done()

	if len(h.handled) != rounds {
		t.Errorf("Engine: expected %v rounds but got %v", rounds, len(h.handled))
	}
}

func TestEngineStopped(t *testing.T) {
	h := TestHandler{
		mux:     &sync.Mutex{},
		handled: []*TxContext{},
	}

	// Create new Engine and register test handler
	e := NewEngine(
		Config{Slots: 100, Partitions: 100},
	)
	e.Use(h.Handler(t))

	// Create a Sarama message channel
	in := make(chan interface{})

	// Run Engine
	go e.Run(in)

	// Feed sarama channel and then close it
	rounds := 1000
	go func() {
		for i := 1; i <= rounds; i++ {
			in <- "test"
			time.Sleep(time.Millisecond)
		}
		close(in)
	}()

	// Sleep and close
	time.Sleep(300 * time.Millisecond)
	e.Close()

	// Wait for Engine to be done
	<-e.Done()

	if len(h.handled) > 500 {
		t.Errorf("Engine: expected max %v rounds but got %v", 500, len(h.handled))
	}

	msgCount := 0
	for range in {
		// We drain messages
		msgCount++
	}

	if len(h.handled)+msgCount != rounds {
		t.Errorf("Engine: expected all %v messages to have been consumed or drained but got consumed=%v drained=%v", rounds, len(h.handled), msgCount)
	}
}

func testSleepingHandler(ctx *TxContext) {
	time.Sleep(ctx.Keys["duration"].(time.Duration))
}

func makeTimeoutContext(i int) *TxContext {
	ctx := NewTxContext()
	ctx.Reset()
	ctx.Prepare([]HandlerFunc{}, log.NewEntry(log.StandardLogger()), nil)

	switch i % 2 {
	case 0:
		ctx.Keys["duration"] = 50 * time.Millisecond
		ctx.Keys["errors"] = 0
	case 1:
		ctx.Keys["duration"] = 100 * time.Millisecond
		ctx.Keys["errors"] = 1
	}
	return ctx
}

func TestTimeoutHandler(t *testing.T) {
	timeoutHandler := TimeoutHandler(testSleepingHandler, 60*time.Millisecond, "Test timeout")

	rounds := 100
	outs := make(chan *TxContext, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		ctx := makeTimeoutContext(i)
		go func(ctx *TxContext) {
			defer wg.Done()
			timeoutHandler(ctx)
			outs <- ctx
		}(ctx)
	}
	wg.Wait()
	close(outs)

	assert.Len(t, outs, rounds, "Timeout: processed contexts count should be correct")

	for out := range outs {
		errCount := out.Keys["errors"].(int)
		assert.Len(t, out.Envelope.Errors, errCount, "Timeout: expected correct count of errors")
	}
}
