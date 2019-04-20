package engine

import (
	"context"
	"sync"
)

var (
	e        *Engine
	initOnce = &sync.Once{}
)

// Init intilialize global Engine
// Configuration is loaded from viper
func Init() {
	initOnce.Do(func() {
		if e != nil {
			return
		}
		config := NewConfig()
		e.SetConfig(&config)
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

// SetConfig set configuration
func SetConfig(conf *Config) {
	e.SetConfig(conf)
}

// Register register a new handler
func Register(handler HandlerFunc) {
	e.Register(handler)
}

// Run starts consuming messages from an input channel
func Run(ctx context.Context, input <-chan interface{}) {
	e.Run(ctx, input)
}

// CleanUp clean engine ressources
func CleanUp() {
	e.CleanUp()
}
