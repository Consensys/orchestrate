package infra

import (
	"fmt"
	"sync"

	"github.com/Shopify/sarama"
)

// Worker for consuming Sarama messages
type Worker struct {
	handlers []HandlerFunc
	handling sync.WaitGroup
	pool     sync.Pool

	// Channel used to limit count of goroutine handling messages
	slots chan struct{}
	// Channel used to indicate runner has terminated
	done chan struct{}
}

// NewWorker creates a new worker
// Slots indicate a count of goroutine that worker can occupy to process message
// You must set `slots > 0`
func NewWorker(handlers []HandlerFunc, slots uint) *Worker {
	if slots == 0 {
		panic(fmt.Errorf("New worker requires at least 1 goroutine slots"))
	}
	return &Worker{
		handlers: handlers,
		handling: sync.WaitGroup{},
		pool:     sync.Pool{New: func() interface{} { return NewContext() }},
		slots:    make(chan struct{}, slots),
		done:     make(chan struct{}, 1),
	}
}

// Run Starts a worker to consume sarama messages
func (w *Worker) Run(messages <-chan *sarama.ConsumerMessage) {
	for {
		msg, ok := <-messages
		if !ok {
			// Message channel has been close.
			// We wait until all messages have been properly consumed
			w.handling.Wait()
			// We indicate that we have gracefully stop
			w.done <- struct{}{}
			// Exit loop
			break
		} else {
			// Acquire a goroutine slot
			w.slots <- struct{}{}

			// Execute message handling in a dedicated goroutine
			go func(msg *sarama.ConsumerMessage) {
				defer func() {
					// Release a goroutine slot
					<-w.slots
				}()
				w.handleMessage(msg)
			}(msg)
		}
	}
}

func (w *Worker) handleMessage(msg *sarama.ConsumerMessage) {
	// Indicate that a new message is being handled
	w.handling.Add(1)

	// Retrieve a re-cycled context
	ctx := w.pool.Get().(*Context)

	defer func(ctx *Context) {
		// Indicate that message has been handled
		w.handling.Done()

		// Re-cycle context object
		w.pool.Put(ctx)
	}(ctx)

	// Prepare context
	ctx.Prepare(w.handlers, msg)

	// Handle context
	ctx.Next()
}

// Done returns a channel indicating if worker is done running
func (w *Worker) Done() <-chan struct{} {
	return w.done
}
