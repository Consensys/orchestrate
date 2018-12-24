package types

import (
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
)

// HandlerFunc is base type for a function processing a Trace
type HandlerFunc func(ctx *Context)

// Context allows us to transmit information through middlewares
type Context struct {
	// T stores information about transaction lifecycle in high level types
	T *Trace
	// Sarama message that triggered Context execution
	Msg *sarama.Message
	// Protobuffer
	pb *tracepb.Trace

	// Keys is a key/value pair
	Keys map[string]interface{}

	// Handlers to be executed on context
	handlers []HandlerFunc
	// Handler being executed
	index int
}

// NewContext creates a new context
func NewContext() *Context {
	t := NewTrace()
	return &Context{
		pb:    &tracepb.Trace{},
		T:     &t,
		Keys:  make(map[string]interface{}),
		index: -1,
	}
}

// Reset re-initialize context
func (ctx *Context) Reset() {
	ctx.Msg = nil
	ctx.pb.Reset()
	ctx.T.Reset()
	ctx.Keys = make(map[string]interface{})
	ctx.handlers = nil
	ctx.index = -1
}

// Next should be used in middleware
// It executes pending handlers
func (ctx *Context) Next() {
	ctx.index++
	for s := len(ctx.handlers); ctx.index < s; ctx.index++ {
		ctx.handlers[ctx.index](ctx)
	}
}

// Error attaches an error to context.
func (ctx *Context) Error(err error) *Error {
	if err == nil {
		panic("err is nil")
	}

	e, ok := err.(*Error)
	if !ok {
		e = &Error{
			Err:  err,
			Type: ErrorTypeUnknown,
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
func (ctx *Context) AbortWithError(err error) *Error {
	ctx.Abort()
	return ctx.Error(err)
}

// Init initialize a context
func (ctx *Context) Init(handlers []HandlerFunc) {
	ctx.Reset()
	ctx.handlers = handlers
}

// Prepare re-initializes context, set handlers and loads sarama message
func (ctx *Context) Prepare(handlers []HandlerFunc, msg *sarama.Message) {
	ctx.Init(handlers)
	ctx.loadMessage(msg)
}

// LoadMessage unmarshal sarama message into protobuffer
func (ctx *Context) loadMessage(msg *sarama.Message) {
	ctx.Msg = msg
	err := proto.Unmarshal(ctx.Msg.Value, ctx.pb)
	if err != nil {
		// Indicate error for a possible middleware to recover it
		e := &Error{
			Err:  err,
			Type: ErrorTypeLoad,
		}
		ctx.Error(e)
	}
}
