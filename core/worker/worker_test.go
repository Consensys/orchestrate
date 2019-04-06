package worker

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
	handled []*Context
}

func (h *TestHandler) Handler(t *testing.T) HandlerFunc {
	return func(ctx *Context) {
		// We add some randomness in time execution
		r := rand.Intn(100)
		time.Sleep(time.Duration(r) * time.Millisecond)
		h.mux.Lock()
		defer h.mux.Unlock()
		h.handled = append(h.handled, ctx)
	}
}

func TestWorker(t *testing.T) {
	h := TestHandler{
		mux:     &sync.Mutex{},
		handled: []*Context{},
	}

	// Create new worker and register test handler
	w := NewWorker(
		Config{Slots: 100, Partitions: 100},
	)
	w.Use(h.Handler(t))

	// Create a Sarama message channel
	in := make(chan interface{})

	// Run worker
	go w.Run(in)

	// Feed sarama channel and then close it
	rounds := 1000
	for i := 1; i <= rounds; i++ {
		in <- "test"
	}
	close(in)

	// Wait for worker to be done
	<-w.Done()

	if len(h.handled) != rounds {
		t.Errorf("Worker: expected %v rounds but got %v", rounds, len(h.handled))
	}
}

func TestWorkerStopped(t *testing.T) {
	h := TestHandler{
		mux:     &sync.Mutex{},
		handled: []*Context{},
	}

	// Create new worker and register test handler
	w := NewWorker(
		Config{Slots: 100, Partitions: 100},
	)
	w.Use(h.Handler(t))

	// Create a Sarama message channel
	in := make(chan interface{})

	// Run worker
	go w.Run(in)

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
	w.Close()

	// Wait for worker to be done
	<-w.Done()

	if len(h.handled) > 500 {
		t.Errorf("Worker: expected max %v rounds but got %v", 500, len(h.handled))
	}

	msgCount := 0
	for range in {
		// We drain messages
		msgCount++
	}

	if len(h.handled)+msgCount != rounds {
		t.Errorf("Worker: expected all %v messages to have been consumed or drained but got consumed=%v drained=%v", rounds, len(h.handled), msgCount)
	}
}

func testSleepingHandler(ctx *Context) {
	time.Sleep(ctx.Keys["duration"].(time.Duration))
}

func makeTimeoutContext(i int) *Context {
	ctx := NewContext()
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
	outs := make(chan *Context, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		ctx := makeTimeoutContext(i)
		go func(ctx *Context) {
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
		assert.Len(t, out.T.Errors, errCount, "Timeout: expected correct count of errors")
	}
}
