package engine

import (
	"context"
	"fmt"
	"hash"
	"hash/fnv"
	"time"

	"math/rand"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"
)

// PartitionKeyFunc are functions returning a key for a message
type PartitionKeyFunc func(message interface{}) []byte

// Engine is an object that allows to consume go channels
type Engine struct {
	conf Config

	handlers []HandlerFunc    // Handlers to apply on each TxContext
	key      PartitionKeyFunc // Function used to create key for dispatching

	handling *sync.WaitGroup // WaitGroup to keep track of messages being consumed and gracefully stop
	ctxPool  *sync.Pool      // Pool used to re-cycle context

	hashPool *sync.Pool

	partitions  []chan interface{} // Partitions correspond to of message
	slots       chan struct{}      // Channel used to limit count of goroutine handling messages
	dying, done chan struct{}      // Channel used to indicate runner has terminated
	closeOnce   *sync.Once

	logger *log.Logger
}

// NewEngine creates a new Engine
func NewEngine(conf Config) *Engine {
	// Validate configuration
	conf.Validate()

	partitions := make([]chan interface{}, conf.Partitions)
	for i := range partitions {
		partitions[i] = make(chan interface{}, conf.Slots)
	}

	return &Engine{
		conf:     conf,
		handlers: []HandlerFunc{},
		handling: &sync.WaitGroup{},
		// By default key is randomly generated
		key:        func(message interface{}) []byte { return []byte(strconv.Itoa(rand.Int())) },
		ctxPool:    &sync.Pool{New: func() interface{} { return NewTxContext() }},
		hashPool:   &sync.Pool{New: func() interface{} { return fnv.New64() }},
		partitions: partitions,
		slots:      make(chan struct{}, conf.Slots),
		dying:      make(chan struct{}),
		done:       make(chan struct{}),
		closeOnce:  &sync.Once{},
		logger:     log.StandardLogger(), // TODO: make possible to use non-standard logrus logger
	}
}

// Partitionner register a function computing partition key on input messages
// Messages having same partition key are treated sequentially
// Messages from distinct partitions are treated in parallel
func (e *Engine) Partitionner(key PartitionKeyFunc) {
	e.key = key
}

// Use add a new handler
func (e *Engine) Use(handler HandlerFunc) {
	e.handlers = append(e.handlers, handler)
}

// dispatch send message to the corresponding partition
func (e *Engine) dispatch(message interface{}) {
	// Compute message partition key
	h := e.hashPool.Get().(hash.Hash64)
	defer e.hashPool.Put(h)
	h.Reset()
	_, err := h.Write(e.key(message))

	if err != nil {
		// Key must be computable for every message
		e.logger.Fatalf("engine: could not compute key on message %v", message)
	}

	// Dispatch message to the corresponding partition key
	e.partitions[h.Sum64()%uint64(len(e.partitions))] <- message
}

func (e *Engine) dispatcher(messages chan interface{}) {
	e.logger.Debugf("engine: start main loop")
dispatcherLoop:
	for {
		select {
		case <-e.dying:
			// Exit loop
			e.logger.Debug("engine: dying")
			break dispatcherLoop
		case msg, ok := <-messages:
			if !ok {
				// Message channel has been close so we also close
				e.logger.Debug("engine: input channel closed")
				e.Close()

				// Exit loop
				break dispatcherLoop
			} else {
				// Indicate that a new message is being handled
				e.handling.Add(1)

				// Acquire a goroutine slot
				e.slots <- struct{}{}

				// Dispatch message
				e.dispatch(msg)
			}
		}
	}
	e.logger.Debugf("engine: left main loop")

	// Close slots channel
	close(e.slots)

	// Close message channels
	for _, channel := range e.partitions {
		close(channel)
	}
}

// handler handles messages from a given parition
func (e *Engine) handler(partition <-chan interface{}) {
	for msg := range partition {
		e.logger.Trace("engine: handle msg")

		// Handle message
		e.handleMessage(msg)

		// Release a goroutine slot
		<-e.slots

		// Indicate that message has been handled
		e.handling.Done()
		e.logger.Trace("engine: msg handled")
	}
}

// Run Starts a Engine to consume messages
func (e *Engine) Run(messages chan interface{}) {
	// Start one handling goroutine per partition
	for _, channel := range e.partitions {
		go e.handler(channel)
	}

	// Start dispatching messages
	e.dispatcher(messages)

	// We wait until all messages have been properly handled
	e.logger.Debugf("engine: wait for all messages to be properly handled")
	e.handling.Wait()

	// We notify that we properly stopped
	close(e.done)
	e.logger.Debugf("engine: done")
}

// Close Engine
func (e *Engine) Close() {
	e.closeOnce.Do(func() {
		e.logger.Debugf("engine: closing...")
		close(e.dying)
	})
}

func (e *Engine) handleMessage(msg interface{}) {
	// Retrieve a re-cycled context
	c := e.ctxPool.Get().(*TxContext)

	// Re-cycle context object
	defer e.ctxPool.Put(c)

	// Prepare context
	c.Prepare(e.handlers, log.NewEntry(e.logger), msg)

	// And calls Next to trigger execution
	c.Next()
}

// Done returns a channel indicating if Engine is done running
func (e *Engine) Done() <-chan struct{} {
	return e.done
}

// TimeoutHandler returns a Handler that runs h with the given time limit
//
// Be careful that if h is a middleware then timeout should cover full execution of the handler
// including pending handlers
func TimeoutHandler(h HandlerFunc, timeout time.Duration, msg string) HandlerFunc {
	return func(ctx *TxContext) {
		// Create timeout context
		timeoutCtx, cancel := context.WithTimeout(ctx.Context(), timeout)
		defer cancel() // We always cancel to avoid memort leak

		// Attach time out context to TxContext
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
