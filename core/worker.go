package core

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
)

// Worker for consuming messages
type Worker struct {
	handlers    []HandlerFunc   // Handlers to apply on each context
	handling    *sync.WaitGroup // WaitGroup to keep track of messages being consumed and gracefully stop
	pool        *sync.Pool      // Pool used to re-cycle context
	slots       chan struct{}   // Channel used to limit count of goroutine handling messages
	dying, done chan struct{}   // Channel used to indicate runner has terminated
	closeOnce   *sync.Once
	logger      *log.Logger
}

// NewWorker creates a new worker
// You indicate a count of goroutine that worker can occupy to process messages
// You must set `slots > 0`
func NewWorker(slots uint) *Worker {
	if slots == 0 {
		panic(fmt.Errorf("Worker requires at least 1 goroutine slot"))
	}

	return &Worker{
		handlers:  []HandlerFunc{},
		handling:  &sync.WaitGroup{},
		pool:      &sync.Pool{New: func() interface{} { return NewContext() }},
		slots:     make(chan struct{}, slots),
		dying:     make(chan struct{}),
		done:      make(chan struct{}),
		closeOnce: &sync.Once{},
		logger:    log.StandardLogger(), // TODO: make possible to use non-standard logrus logger
	}
}

// Use add a new handler
func (w *Worker) Use(handler HandlerFunc) {
	w.handlers = append(w.handlers, handler)
}

// Run Starts a worker to consume messages
func (w *Worker) Run(messages chan interface{}) {
	w.logger.Debugf("worker: start main loop")
runningLoop:
	for {
		select {
		case msg, ok := <-messages:
			if !ok {
				// Message channel has been close so we also close
				w.logger.Debug("worker: input channel closed")
				w.Close()

				// Exit loop
				break runningLoop
			} else {
				// Indicate that a new message is being handled
				w.handling.Add(1)

				// Acquire a goroutine slot
				w.slots <- struct{}{}

				go func(msg interface{}) {
					w.logger.Trace("worker: handle msg")
					// Handle message in a dedicated goroutine
					w.handleMessage(msg)

					// Release a goroutine slot
					<-w.slots

					// Indicate that message has been handled
					w.handling.Done()
					w.logger.Trace("worker: msg handled")
				}(msg)
			}
		case <-w.dying:
			// Exit loop
			w.logger.Debug("worker: dying")
			break runningLoop
		}
	}

	// Close slots channel
	close(w.slots)

	// We wait until all messages have been properly handled
	w.logger.Debugf("worker: left main loop, wait for messages to be properly handled")
	w.handling.Wait()

	// We notify that we properly stopped
	close(w.done)
	w.logger.Debugf("worker: done")
}

// Close worker
func (w *Worker) Close() {
	w.closeOnce.Do(func() {
		w.logger.Debugf("worker: closing...")
		close(w.dying)
	})
}

func (w *Worker) handleMessage(msg interface{}) {
	// Retrieve a re-cycled context
	ctx := w.pool.Get().(*Context)

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
