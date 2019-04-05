package worker

import (
	"context"

	log "github.com/sirupsen/logrus"
	common "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
	trace "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/trace"
)

// HandlerFunc is base type for an handler function processing a Context
type HandlerFunc func(ctx *Context)

// Context is the most important part of a worker.
// It allows to pass variables between handlers
type Context struct {
	// ctx is a go context that is attach to the worker Context
	// It allows to carry deadlines, cancelation signals, etc. between handlers
	// It is not recommended to do
	ctx context.Context

	// T stores all information about transaction lifecycle
	T *trace.Trace

	// Message that triggered Context execution (typically a sarama.ConsumerMessage)
	Msg interface{}

	// Keys is a key/value pair
	Keys map[string]interface{}

	// chain of handlers to be executed on context
	handlers []HandlerFunc

	// index of the handler being executed
	index int

	// Logger logrus log entry for this context execution
	Logger *log.Entry
}

// NewContext creates a new context
func NewContext() *Context {
	return &Context{
		T:     &trace.Trace{},
		Keys:  make(map[string]interface{}),
		index: -1,
	}
}

// Context return go context
func (ctx *Context) Context() context.Context {
	return ctx.ctx
}

// Reset re-initialize context
func (ctx *Context) Reset() {
	ctx.ctx = nil
	ctx.Msg = nil
	ctx.T.Reset()
	ctx.Keys = make(map[string]interface{})
	ctx.handlers = nil
	ctx.index = -1
	ctx.Logger = nil
}

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler
func (ctx *Context) Next() {
	ctx.index++
	for s := len(ctx.handlers); ctx.index < s; ctx.index++ {
		ctx.handlers[ctx.index](ctx)
	}
}

// Error attaches an error to context.
func (ctx *Context) Error(err error) *common.Error {
	if err == nil {
		panic("err is nil")
	}

	e, ok := err.(*common.Error)
	if !ok {
		e = &common.Error{
			Message: err.Error(),
		}
	}
	ctx.T.Errors = append(ctx.T.Errors, e)

	return e
}

// Abort prevents pending handlers to be executed
func (ctx *Context) Abort() {
	ctx.index = len(ctx.handlers)
}

// AbortWithError calls `Abort()` and `Error()``
func (ctx *Context) AbortWithError(err error) *common.Error {
	ctx.Abort()
	return ctx.Error(err)
}

// Prepare re-initializes context, set handlers, set logger and set message
func (ctx *Context) Prepare(handlers []HandlerFunc, logger *log.Entry, msg interface{}) {
	ctx.Reset()
	ctx.handlers = handlers
	ctx.Msg = msg
	ctx.Logger = logger
}

// WithContext attach a go context on worker Context
func WithContext(ctx context.Context, context *Context) *Context {
	context.ctx = ctx
	return context
}
