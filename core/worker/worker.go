package worker

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)

// PartitionKeyFunc are functions returning a key for a message
type PartitionKeyFunc func(message interface{}) []byte

// Worker allows to consume messages on an input channel
type Worker struct {
	// Worker configuration object
	conf Config

	// chain of handlers to be to be executed
	handlers []HandlerFunc

	// running keeps track of the number of running loops
	running   int64
	cleanOnce *sync.Once

	// ctxPool is a pool used to re-cycle context objects
	ctxPool *sync.Pool

	// slots is a channel used to limit the number of messages treated concurently by the worker
	slots chan struct{}

	// Worker logger
	logger *log.Logger
}

// NewWorker creates a new worker
// You indicate a count of goroutine that worker can occupy to process messages
// You must set `slots > 0`
func NewWorker(conf Config) *Worker {
	// Validate configuration
	conf.Validate()

	return &Worker{
		conf:      conf,
		handlers:  []HandlerFunc{},
		running:   0,
		cleanOnce: &sync.Once{},
		// By default key is randomly generated
		ctxPool: &sync.Pool{New: func() interface{} { return NewContext() }},
		slots:   make(chan struct{}, conf.Slots),
		logger:  log.StandardLogger(), // TODO: make possible to use non-standard logrus logger
	}
}

// Use add a new handler
func (w *Worker) Use(handler HandlerFunc) {
	w.handlers = append(w.handlers, handler)
}

// Run starts consuming messages from an input channel
//
// Run will gracefully interupt either if
// - provided ctx is cancelled
// - input channel is closed
//
// Once you have stopped consuming from an input channel, you should not start to consuming
// from a new channel using Run() or it will panic (if you need to start consuming from a new channel you
// should create a new worker)
func (w *Worker) Run(ctx context.Context, input <-chan interface{}) {
	// Context must be not nil
	if ctx == nil {
		panic("nil context")
	}

	// Increment count of input channels being consumed
	count := atomic.AddInt64(&w.running, 1)
	w.logger.WithFields(log.Fields{
		"inputs.count": count,
	}).Debugf("worker: start running loop")

runningLoop:
	for {
		select {
		case msg, ok := <-input:
			if !ok {
				// Input channel has been close so we leave the loop
				break runningLoop
			}

			// Acquire a message slot
			w.slots <- struct{}{}

			// Handle message
			w.handleMessage(msg)

			// Release a message slot
			<-w.slots
		case <-ctx.Done():
			// Context has timeout or been cancelled so we leave the loop
			break runningLoop
		}
	}

	// Decrement count of input channels being consumed
	count = atomic.AddInt64(&w.running, -1)
	w.logger.WithFields(log.Fields{
		"inputs.count": count,
	}).Debugf("worker: left running loop")
}

// CleanUp clean worker ressources
//
// After interupting execution of all Run() calls you should always call CleanUp
// to avoid memory leak
func (w *Worker) CleanUp() {
	w.cleanOnce.Do(func() {
		close(w.slots)
	})
}

func (w *Worker) handleMessage(msg interface{}) {
	// Retrieve a re-cycled context
	c := w.ctxPool.Get().(*Context)

	// Re-cycle context object
	defer w.ctxPool.Put(c)

	// Prepare context
	c.Prepare(w.handlers, log.NewEntry(w.logger), msg)

	// Calls Next to trigger execution
	c.Next()
}

// TimeoutHandler returns a Handler that runs h with the given time limit
//
// Be careful that if h is a middleware then timeout should cover full execution of the handler
// including pending handlers
func TimeoutHandler(h HandlerFunc, timeout time.Duration, msg string) HandlerFunc {
	return func(ctx *Context) {
		// Create timeout context
		timeoutCtx, cancel := context.WithTimeout(ctx.Context(), timeout)
		defer cancel() // We always cancel to avoid memort leak

		// Attach time out context to worker context
		ctx.WithContext(timeoutCtx)

		// Prepare channels
		done := make(chan struct{})
		panicChan := make(chan interface{}, 1)
		defer close(panicChan)

		// Execute handler
		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			h(ctx)
			close(done)
		}()

		// Wait for handler execution to complete or for a timeout or panic
		select {
		case <-done:
			// Execution properly completed
		case <-timeoutCtx.Done():
			// Execution timed out
			ctx.Error(fmt.Errorf(msg))
		case p := <-panicChan:
			// Execution panic so we forward
			panic(p)
		}
	}
}
