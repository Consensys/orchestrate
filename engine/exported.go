package engine

import (
	"context"
)

func init() {
	e = NewEngine(nil)
}

var e *Engine

// Init intilialize global Engine
// Configuration is loaded from viper
func Init() {
	config := NewConfig()
	e.SetConfig(&config)
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

// Use register a new handler
func Use(handler HandlerFunc) {
	e.Use(handler)
}

// Run starts consuming messages from an input channel
func Run(ctx context.Context, input <-chan interface{}) {
	e.Run(ctx, input)
}

// CleanUp clean engine ressources
func CleanUp() {
	e.CleanUp()
}
