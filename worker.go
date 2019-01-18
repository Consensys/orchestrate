package core

import (
	"fmt"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// Worker for consuming Sarama messages
type Worker struct {
	handlers []types.HandlerFunc // Handlers to apply on each context
	handling *sync.WaitGroup     // WaitGroup to keep track of messages being consumed and gracefully stop
	pool     *sync.Pool          // Pool used to re-cycle context
	slots    chan struct{}       // Channel used to limit count of goroutine handling messages
	done     chan struct{}       // Channel used to indicate runner has terminated
}

// NewWorker creates a new worker
// Slots indicate a count of goroutine that worker can occupy to process messages
// You must set `slots > 0`
func NewWorker(slots uint) *Worker {
	if slots == 0 {
		panic(fmt.Errorf("Worker requires at least 1 goroutine slots"))
	}
	return &Worker{
		handlers: []types.HandlerFunc{},
		handling: &sync.WaitGroup{},
		pool:     &sync.Pool{New: func() interface{} { return types.NewContext() }},
		slots:    make(chan struct{}, slots),
		done:     make(chan struct{}, 1),
	}
}

// Use add a new handler
func (w *Worker) Use(handler types.HandlerFunc) {
	w.handlers = append(w.handlers, handler)
}

// Run Starts a worker to consume sarama messages
func (w *Worker) Run(messages chan interface{}) {
	for {
		msg, ok := <-messages
		if !ok {
			// Message channel has been close.
			// We wait until all messages have been properly consumed
			w.handling.Wait()
			// We indicate that we have gracefully stoped
			w.done <- struct{}{}
			// Exit loop
			break
		} else {
			// Indicate that a new message is being handled
			w.handling.Add(1)

			// Acquire a goroutine slot
			w.slots <- struct{}{}
			// Execute message handling in a dedicated goroutine
			go func(msg interface{}) {
				defer func() {
					// Release a goroutine slot
					<-w.slots
				}()
				w.handleMessage(msg)

				// Indicate that message has been handled
				w.handling.Done()
			}(msg)
		}
	}
}

func (w *Worker) handleMessage(msg interface{}) {
	// Retrieve a re-cycled context
	ctx := w.pool.Get().(*types.Context)

	defer func(ctx *types.Context) {
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
