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

	// ctx is a go context that is attached to the worker Context
	// It allows to carry deadlines, cancelation signals, etc. between handlers
	//
	// This approach is not recommended by go context documentation
	// c.f. https://golang.org/pkg/context/#pkg-overview
	//
	// Still this recommendation against has been actively questioned
	// (c.f https://github.com/golang/go/issues/22602)
	// Also net/http has been following this implementation for the Request object
	// (c.f. https://github.com/golang/go/blob/master/src/net/http/request.go#L107)
	ctx context.Context
}

// NewContext creates a new context
func NewContext() *Context {
	return &Context{
		T:     &trace.Trace{},
		Keys:  make(map[string]interface{}),
		index: -1,
	}
}

// Reset re-initialize context
func (c *Context) Reset() {
	c.ctx = nil
	c.Msg = nil
	c.T.Reset()
	c.Keys = make(map[string]interface{})
	c.handlers = nil
	c.index = -1
	c.Logger = nil
}

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler
func (c *Context) Next() {
	c.index++
	for s := len(c.handlers); c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// Error attaches an error to context.
func (c *Context) Error(err error) *common.Error {
	if err == nil {
		panic("err is nil")
	}

	e, ok := err.(*common.Error)
	if !ok {
		e = &common.Error{
			Message: err.Error(),
		}
	}
	c.T.Errors = append(c.T.Errors, e)

	return e
}

// Abort prevents pending handlers to be executed
func (c *Context) Abort() {
	c.index = len(c.handlers)
}

// AbortWithError calls `Abort()` and `Error()``
func (c *Context) AbortWithError(err error) *common.Error {
	c.Abort()
	return c.Error(err)
}

// Prepare re-initializes context, set handlers, set logger and set message
func (c *Context) Prepare(handlers []HandlerFunc, logger *log.Entry, msg interface{}) {
	c.Reset()
	c.handlers = handlers
	c.Msg = msg
	c.Logger = logger
}

// Context returns the go context attached to the worker Context.
// To change the context, use WithContext.
//
// The returned context is always non-nil; it defaults to the background context.
func (c *Context) Context() context.Context {
	if c.ctx != nil {
		return c.ctx
	}
	return context.Background()
}

// WithContext attach a go context to a worker Context
// The go context provided as argument must be non nil or WithContext will panic
func (c *Context) WithContext(ctx context.Context) *Context {
	if ctx == nil {
		panic("nil context")
	}
	c.ctx = ctx
	return c
}
