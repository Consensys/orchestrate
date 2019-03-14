package core

import (
	"fmt"
	"hash"
	"hash/fnv"

	"math/rand"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"
)

// PartitionKeyFunc are functions returning a key for a message
type PartitionKeyFunc func(message interface{}) []byte

// Worker allows to consume messages on an input channel
type Worker struct {
	handlers []HandlerFunc    // Handlers to apply on each context
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

// NewWorker creates a new worker
// You indicate a count of goroutine that worker can occupy to process messages
// You must set `slots > 0`
func NewWorker(slots uint) *Worker {
	if slots == 0 {
		panic(fmt.Errorf("Worker requires at least 1 goroutine slot"))
	}

	partitions := make([]chan interface{}, slots)
	for i := range partitions {
		partitions[i] = make(chan interface{})
	}

	return &Worker{
		handlers: []HandlerFunc{},
		handling: &sync.WaitGroup{},
		// By default key is randomly generated
		key:        func(message interface{}) []byte { return []byte(strconv.Itoa(rand.Int())) },
		ctxPool:    &sync.Pool{New: func() interface{} { return NewContext() }},
		hashPool:   &sync.Pool{New: func() interface{} { return fnv.New64() }},
		partitions: partitions,
		slots:      make(chan struct{}, slots),
		dying:      make(chan struct{}),
		done:       make(chan struct{}),
		closeOnce:  &sync.Once{},
		logger:     log.StandardLogger(), // TODO: make possible to use non-standard logrus logger
	}
}

// Partitionner register a function computing partition key on input messages
// Messages having same partition key are treated sequentially
// Messages from distinct partitions are treated in parallel
func (w *Worker) Partitionner(key PartitionKeyFunc) {
	w.key = key
}

// Use add a new handler
func (w *Worker) Use(handler HandlerFunc) {
	w.handlers = append(w.handlers, handler)
}

// dispatch send message to the corresponding partition
func (w *Worker) dispatch(message interface{}) {
	// Compute message partition key
	h := w.hashPool.Get().(hash.Hash64)
	defer w.hashPool.Put(h)
	h.Reset()
	_, err := h.Write(w.key(message))

	if err != nil {
		// Key must be computable for every message
		w.logger.Fatalf("worker: could not compute key on message %v", message)
	}

	// Dispatch message to the corresponding partition key
	w.partitions[h.Sum64()%uint64(len(w.partitions))] <- message
}

func (w *Worker) dispatcher(messages chan interface{}) {
	w.logger.Debugf("worker: start main loop")
dispatcherLoop:
	for {
		select {
		case <-w.dying:
			// Exit loop
			w.logger.Debug("worker: dying")
			break dispatcherLoop
		case msg, ok := <-messages:
			if !ok {
				// Message channel has been close so we also close
				w.logger.Debug("worker: input channel closed")
				w.Close()

				// Exit loop
				break dispatcherLoop
			} else {
				// Indicate that a new message is being handled
				w.handling.Add(1)

				// Acquire a goroutine slot
				w.slots <- struct{}{}

				// Dispatch message
				w.dispatch(msg)
			}
		}
	}
	w.logger.Debugf("worker: left main loop")

	// Close slots channel
	close(w.slots)

	// Close message channels
	for _, channel := range w.partitions {
		close(channel)
	}
}

// handler handles messages from a given parition
func (w *Worker) handler(partition <-chan interface{}) {
	for msg := range partition {
		w.logger.Trace("worker: handle msg")

		// Handle message
		w.handleMessage(msg)

		// Release a goroutine slot
		<-w.slots

		// Indicate that message has been handled
		w.handling.Done()
		w.logger.Trace("worker: msg handled")
	}
}

// Run Starts a worker to consume messages
func (w *Worker) Run(messages chan interface{}) {
	// Start one handling goroutine per partition
	for _, channel := range w.partitions {
		go w.handler(channel)
	}

	// Start dispatching messages
	w.dispatcher(messages)

	// We wait until all messages have been properly handled
	w.logger.Debugf("worker: wait for all messages to be properly handled")
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
	ctx := w.ctxPool.Get().(*Context)

	// Re-cycle context object
	defer w.ctxPool.Put(ctx)

	// Prepare context
	ctx.Prepare(w.handlers, log.NewEntry(w.logger), msg)

	// Handle context
	ctx.Next()
}

// Done returns a channel indicating if worker is done running
func (w *Worker) Done() <-chan struct{} {
	return w.done
}
