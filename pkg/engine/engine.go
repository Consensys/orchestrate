package engine

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
)

// HandlerFunc is base type for an handler function processing a Context
type HandlerFunc func(txctx *TxContext)

// Engine is an object that allows to consume go channels
type Engine struct {
	// Configuration object
	conf *Config

	// chain of handlers to be to be executed
	handlers        []HandlerFunc
	wrapperHandlers []HandlerFunc

	// running keeps track of the number of running loops
	running   int64
	cleanOnce *sync.Once

	// ctxPool is a pool to re-cycle Envelope
	ctxPool *sync.Pool

	// slots is a channel to limit the number of messages treated concurrently by the Engine
	slots chan struct{}

	// Logger
	logger *log.Logger

	mux *sync.Mutex
}

// NewEngine creates a new Engine
func NewEngine(logger *log.Logger, conf *Config) (e *Engine) {
	e = &Engine{
		handlers:  []HandlerFunc{},
		running:   0,
		cleanOnce: &sync.Once{},
		ctxPool:   &sync.Pool{New: func() interface{} { return NewTxContext() }},
		mux:       &sync.Mutex{},
		logger:    logger,
	}

	if conf != nil {
		e.SetConfig(conf)
	}

	return
}

// SetConfig set Engine configuration
func (e *Engine) SetConfig(conf *Config) {
	if conf == nil {
		e.logger.Fatal("nil configuration passed to engine")
	}

	if err := conf.Validate(); err != nil {
		e.logger.WithError(err).Fatal("cannot validate engine configuration")
	}

	e.mux.Lock()
	e.conf = conf
	e.mux.Unlock()
}

// Register register a new handler
func (e *Engine) Register(handler HandlerFunc) {
	e.mux.Lock()
	e.handlers = append(e.handlers, handler)
	e.mux.Unlock()
}

// Register register a new WrapperHandler
func (e *Engine) RegisterWrapper(handler HandlerFunc) {
	e.mux.Lock()
	e.wrapperHandlers = append(e.wrapperHandlers, handler)
	e.mux.Unlock()
}

// Run starts consuming messages from an input channel
//
// Run will gracefully interrupt either if
// - provided ctx is canceled
// - input channel is closed
//
// Once you have stopped consuming from an input channel, you should not start to consuming
// from a new channel using Run() or it will panic (if you need to start consuming from a new channel you
// should create a new Engine)
func (e *Engine) Run(ctx context.Context, input <-chan Msg) {
	// Context must be not nil
	if ctx == nil {
		panic("nil context")
	}

	e.mux.Lock()
	// Ensure config has been attached
	if e.conf == nil {
		panic("nil configuration (call SetConfig() before running engine)")
	}

	// Initialize slots channel
	if e.slots == nil {
		e.slots = make(chan struct{}, e.conf.Slots)
	}
	e.mux.Unlock()

	// Increment count of input channels being consumed
	count := atomic.AddInt64(&e.running, 1)
	e.logger.WithField("loops.count", count).Debug("engine: start running loop")

	bckOff := backoff.NewExponentialBackOff()
	bckOff.MaxInterval = time.Second * 15
	bckOff.MaxElapsedTime = time.Minute * 5
runningLoop:
	for {
		select {
		case msg, ok := <-input:
			if !ok {
				// Input channel has been close so we leave the loop
				break runningLoop
			}

			// Acquire a message slot
			e.slots <- struct{}{}

			err := backoff.RetryNotify(
				func() error {
					// Handle message
					err := e.handleMessage(ctx, msg)
					if err == context.DeadlineExceeded || err == context.Canceled || ctx.Err() != nil {
						if err == nil {
							err = ctx.Err()
						}
						log.WithContext(ctx).WithError(err).Info("exiting listener session...")
						return backoff.Permanent(err)
					}

					return err
				},
				bckOff,
				func(err error, duration time.Duration) {
					e.logger.WithError(err).Warnf("error processing msg %q, retrying in %v...", msg.Key(), duration)
				},
			)

			if err != nil {
				e.logger.WithError(err).Error("engine: left running loop")
				break runningLoop
			}

			// Release a message slot
			<-e.slots
		case <-ctx.Done():
			// Context has timeout or been canceled so we leave the loop
			break runningLoop
		}
	}

	// Decrement count of input channels being consumed
	count = atomic.AddInt64(&e.running, -1)
	e.logger.WithField("loops.count", count).Debug("engine: left running loop")
}

// CleanUp clean Engine resources
//
// After completion of each Run() calls you should always call CleanUp
// to avoid memory leak and be able to re-initialize Engine
//
// Do not call CleanUp() before every calls to Run() have properly finished
// otherwise the Engine will panic
func (e *Engine) CleanUp() {
	e.cleanOnce.Do(func() {
		e.mux.Lock()
		if e.slots != nil {
			close(e.slots)
			e.slots = nil
		}
		e.mux.Unlock()
	})
}

func (e *Engine) handleMessage(ctx context.Context, msg Msg) error {
	// Retrieve a re-cycled context
	txctx := e.ctxPool.Get().(*TxContext)

	// Re-cycle context object
	defer e.ctxPool.Put(txctx)

	// Prepare context & calls Next to trigger execution
	txctx.
		Prepare(
			log.FromContext(ctx),
			msg,
		).
		WithContext(ctx).
		applyHandlers(e.handlers...).
		applyHandlers(e.wrapperHandlers...)

	if err := txctx.HasRetryMsgErr(); err != nil {
		return err
	}

	return nil
}

// TimeoutHandler returns a Handler that runs h with the given time limit
//
// Be careful that if h is a middleware then timeout should cover full execution of the handler
// including pending handlers
func TimeoutHandler(h HandlerFunc, timeout time.Duration, msg string) HandlerFunc {
	return func(txctx *TxContext) {
		// Create timeout context
		timeoutCtx, cancel := context.WithTimeout(txctx.Context(), timeout)
		defer cancel() // We always cancel to avoid memory leak

		// Attach time out context to Envelope
		txctx.WithContext(timeoutCtx)

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
			h(txctx)
			close(done)
		}()

		// Wait for handler execution to complete or for a timeout or panic
		select {
		case <-done:
			// Execution properly completed
		case <-timeoutCtx.Done():
			// Execution timed out
			_ = txctx.Error(fmt.Errorf(msg))
		case p := <-panicChan:
			// Execution panic so we forward
			panic(p)
		}
	}
}
