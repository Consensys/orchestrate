package engine

import (
	"context"

	log "github.com/sirupsen/logrus"
	common "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
)

// HandlerFunc is base type for an handler function processing a Context
type HandlerFunc func(txctx *TxContext)

// TxContext is the most important part of an engine.
// It allows to pass variables between handlers
type TxContext struct {
	// Envelope stores all information about transaction lifecycle
	Envelope *envelope.Envelope

	// Message that triggered TxContext execution (typically a sarama.ConsumerMessage)
	Msg interface{}

	// Keys is a key/value pair
	Keys map[string]interface{}

	// chain of handlers to be executed on TxContext
	handlers []HandlerFunc

	// index of the handler being executed
	index int

	// Logger logrus log entry for this TxContext execution
	Logger *log.Entry

	// ctx is a go context that is attached to the TxContext
	// It allows to carry deadlines, cancelation signals, ettxctx. between handlers
	//
	// This approach is not recommended by go context documentation
	// txctx.f. https://golang.org/pkg/context/#pkg-overview
	//
	// Still this recommendation against has been actively questioned
	// (txctx.f https://github.com/golang/go/issues/22602)
	// Also net/http has been following this implementation for the Request object
	// (txctx.f. https://github.com/golang/go/blob/master/src/net/http/request.go#L107)
	ctx context.Context
}

// NewTxContext creates a new TxContext
func NewTxContext() *TxContext {
	return &TxContext{
		Envelope: &envelope.Envelope{},
		Keys:     make(map[string]interface{}),
		index:    -1,
	}
}

// Reset re-initialize TxContext
func (txctx *TxContext) Reset() {
	txctx.ctx = nil
	txctx.Msg = nil
	txctx.Envelope.Reset()
	txctx.Keys = make(map[string]interface{})
	txctx.handlers = nil
	txctx.index = -1
	txctx.Logger = nil
}

// Next should be used only inside middleware
// It executes the pending handlers in the chain inside the calling handler
func (txctx *TxContext) Next() {
	txctx.index++
	for s := len(txctx.handlers); txctx.index < s; txctx.index++ {
		txctx.handlers[txctx.index](txctx)
	}
}

// Error attaches an error to TxContext
func (txctx *TxContext) Error(err error) *common.Error {
	if err == nil {
		panic("err is nil")
	}

	e, ok := err.(*common.Error)
	if !ok {
		e = &common.Error{
			Message: err.Error(),
		}
	}
	txctx.Envelope.Errors = append(txctx.Envelope.Errors, e)

	return e
}

// Abort prevents pending handlers to be executed
func (txctx *TxContext) Abort() {
	txctx.index = len(txctx.handlers)
}

// AbortWithError calls `Abort()` and `Error()``
func (txctx *TxContext) AbortWithError(err error) *common.Error {
	txctx.Abort()
	return txctx.Error(err)
}

// Prepare re-initializes TxContext, set handlers, set logger and set message
func (txctx *TxContext) Prepare(handlers []HandlerFunc, logger *log.Entry, msg interface{}) *TxContext {
	txctx.Reset()
	txctx.handlers = handlers
	txctx.Msg = msg
	txctx.Logger = logger

	return txctx
}

// Context returns the go context attached to TxContext.
// To change the go context, use WithContext.
//
// The returned context is always non-nil; it defaults to the background context.
func (txctx *TxContext) Context() context.Context {
	if txctx.ctx != nil {
		return txctx.ctx
	}
	return context.Background()
}

// WithContext attach a go context to TxContext
// The go context provided as argument must be non nil or WithContext will panic
func (txctx *TxContext) WithContext(ctx context.Context) *TxContext {
	if ctx == nil {
		panic("nil context")
	}
	txctx.ctx = ctx
	return txctx
}
