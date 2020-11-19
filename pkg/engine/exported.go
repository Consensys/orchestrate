package engine

import (
	"context"
	"sync"
)

const component = "engine"

var (
	e        *Engine
	initOnce = &sync.Once{}
)

// Init initialize global Engine
// Configuration is loaded from viper
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if e != nil {
			return
		}
		conf := NewConfig()
		e = NewEngine(&conf)
	})
}

// SetGlobalEngine set global engine
func SetGlobalEngine(engine *Engine) {
	e = engine
}

// GlobalEngine returns global engine
func GlobalEngine() *Engine {
	return e
}

// Register register a new handler
func Register(handler HandlerFunc) {
	e.Register(handler)
}

// Register register a new handler to execute at the end
func RegisterWrapper(handler HandlerFunc) {
	e.RegisterWrapper(handler)
}

// Run starts consuming messages from an input channel
func Run(ctx context.Context, input <-chan Msg) {
	e.Run(ctx, input)
}
