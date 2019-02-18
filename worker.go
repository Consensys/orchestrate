package core

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// Worker for consuming Sarama messages
type Worker struct {
	handlers    []types.HandlerFunc // Handlers to apply on each context
	handling    *sync.WaitGroup     // WaitGroup to keep track of messages being consumed and gracefully stop
	pool        *sync.Pool          // Pool used to re-cycle context
	slots       chan struct{}       // Channel used to limit count of goroutine handling messages
	dying, done chan struct{}       // Channel used to indicate runner has terminated
	closeOnce   *sync.Once
	logger      *log.Logger
}

// NewWorker creates a new worker
// You indicate a count of goroutine that worker can occupy to process messages
// You must set `slots > 0`
func NewWorker(slots uint) *Worker {
	if slots == 0 {
		panic(fmt.Errorf("Worker requires at least 1 goroutine slots"))
	}

	return &Worker{
		handlers:  []types.HandlerFunc{},
		handling:  &sync.WaitGroup{},
		pool:      &sync.Pool{New: func() interface{} { return types.NewContext() }},
		slots:     make(chan struct{}, slots),
		dying:     make(chan struct{}),
		done:      make(chan struct{}),
		closeOnce: &sync.Once{},
		logger:    log.StandardLogger(), // TODO: make possible to use non-standard logrus logger
	}
}

// Use add a new handler
func (w *Worker) Use(handler types.HandlerFunc) {
	w.handlers = append(w.handlers, handler)
}

// Run Starts a worker to consume sarama messages
func (w *Worker) Run(messages chan interface{}) {
runningLoop:
	for {
		select {
		case msg, ok := <-messages:
			if !ok {
				// Message channel has been close so we also close
				w.Close()
			} else {
				// Indicate that a new message is being handled
				w.handling.Add(1)

				// Acquire a goroutine slot
				w.slots <- struct{}{}

				go func(msg interface{}) {
					// Handle message in a dedicated goroutine
					w.handleMessage(msg)

					// Release a goroutine slot
					<-w.slots

					// Indicate that message has been handled
					w.handling.Done()
				}(msg)
			}
		case <-w.dying:
			// We wait until all messages have been properly handled
			w.handling.Wait()

			// Exit loop
			break runningLoop

		}
	}

	// We notify that we properly stopped
	close(w.done)
}

// Close worker
func (w *Worker) Close() {
	w.closeOnce.Do(func() {
		close(w.dying)
	})
}

func (w *Worker) handleMessage(msg interface{}) {
	// Retrieve a re-cycled context
	ctx := w.pool.Get().(*types.Context)

	// Re-cycle context object
	defer w.pool.Put(ctx)

	// Prepare context
	ctx.Prepare(w.handlers, log.NewEntry(w.logger), msg)

	// Handle context
	ctx.Next()
}

// Done returns a channel indicating if worker is done running
func (w *Worker) Done() <-chan struct{} {
	return w.done
}
