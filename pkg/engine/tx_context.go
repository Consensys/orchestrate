package engine

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	ierror "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/error"
)

// Envelope is the most important part of an engine.
// It allows to pass variables between handlers
type TxContext struct {
	// Envelope stores all information about transaction lifecycle
	Envelope *tx.Envelope

	// Input message
	In Msg

	// Array of sequences of handlers to execute on a given context
	stack []*sequence

	// Logger logrus log entry for this Envelope execution
	Logger *log.Entry

	// ctx is a go context that is attached to the Envelope
	// It allows to carry deadlines, cancellation signals, etc. between handlers
	//
	// This approach is not recommended by go context documentation
	// txctx.f. https://golang.org/pkg/context/#pkg-overview
	//
	// Still this recommendation against has been actively questioned
	// (txctx.f https://github.com/golang/go/issues/22602)
	// Also net/http has been following this implementation for the Envelope object
	// (txctx.f. https://github.com/golang/go/blob/master/src/net/http/request.go#L107)
	ctx context.Context
}

// NewTxContext creates a new Envelope
func NewTxContext() *TxContext {
	return &TxContext{
		Envelope: tx.NewEnvelope(),
	}
}

// Reset re-initialize Envelope
func (txctx *TxContext) Reset() {
	txctx.ctx = nil
	txctx.Envelope = tx.NewEnvelope()
	txctx.In = nil
	txctx.stack = nil
	txctx.Logger = nil
}

// Next should be used only inside middleware
// It executes the pending handlers in the chain inside the calling handler
func (txctx *TxContext) Next() {
	if len(txctx.stack) > 0 {
		txctx.stack[len(txctx.stack)-1].next()
	}
}

// Error attaches an error to Envelope
func (txctx *TxContext) Error(err error) *ierror.Error {
	if err == nil {
		panic("err is nil")
	}

	_ = txctx.Envelope.AppendError(errors.FromError(err))

	return txctx.Envelope.GetErrors()[len(txctx.Envelope.Errors)-1]
}

// Abort prevents pending handlers to be executed
func (txctx *TxContext) Abort() {
	for s := len(txctx.stack) - 1; s >= 0; s-- {
		txctx.stack[s].abort()
	}
}

// AbortWithError calls `Abort()` and `Error()``
func (txctx *TxContext) AbortWithError(err error) *ierror.Error {
	txctx.Abort()
	return txctx.Error(err)
}

// Prepare re-initializes Envelope, set handlers, set logger and set message
func (txctx *TxContext) Prepare(logger *log.Entry, msg Msg) *TxContext {
	txctx.Reset()
	txctx.In = msg
	txctx.Logger = logger
	return txctx
}

type txCtxKey string

// Set is used to store a new key/value pair exclusively for this context
func (txctx *TxContext) Set(key string, value interface{}) {
	txctx.WithContext(context.WithValue(txctx.Context(), txCtxKey(key), value))
}

// Get returns the value for the given key
func (txctx *TxContext) Get(key string) interface{} {
	return txctx.Context().Value(txCtxKey(key))
}

// Context returns the go context attached to Envelope.
// To change the go context, use WithContext.
//
// The returned context is always non-nil; it defaults to the background context.
func (txctx *TxContext) Context() context.Context {
	if txctx.ctx != nil {
		return txctx.ctx
	}
	return context.Background()
}

// WithContext attach a go context to Envelope
// The go context provided as argument must be non nil or WithContext will panic
func (txctx *TxContext) WithContext(ctx context.Context) *TxContext {
	if ctx == nil {
		panic("nil context")
	}
	txctx.ctx = ctx
	return txctx
}

func (txctx *TxContext) applyHandlers(handlers ...HandlerFunc) {
	// Recycle sequence
	seq := seqPool.Get().(*sequence)
	defer seqPool.Put(seq)

	// Initialize sequence
	seq.index = -1
	seq.handlers = handlers
	seq.txctx = txctx

	// Attach the sequence to the Envelope
	txctx.stack = append(txctx.stack, seq)

	// Execute sequence
	seq.next()

	// Once executed remove the sequence
	txctx.stack = txctx.stack[:len(txctx.stack)-1]
}

type sequence struct {
	// chain of handlers to be executed in the sequence
	handlers []HandlerFunc

	// index of the handler being executed
	index int

	// context the sequence is attached to
	txctx *TxContext
}

// sequences are pooled to relieve pressure on garbage collector
var seqPool = sync.Pool{
	New: func() interface{} { return &sequence{index: -1} },
}

func (seq *sequence) next() {
	seq.index++
	for s := len(seq.handlers); seq.index < s; seq.index++ {
		seq.handlers[seq.index](seq.txctx)
	}
}

func (seq *sequence) abort() {
	seq.index = len(seq.handlers)
}

// CombineHandlers returns a composite of several handlers
func CombineHandlers(handlers ...HandlerFunc) HandlerFunc {
	return func(txctx *TxContext) {
		txctx.applyHandlers(handlers...)
	}
}
