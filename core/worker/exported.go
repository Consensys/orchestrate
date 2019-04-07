package worker

import (
	"context"
)

func init() {
	w = NewWorker(nil)
}

var w *Worker

// Init intilialize global worker
// Configuration is loaded from viper
func Init() {
	config := NewConfig()
	w.SetConfig(&config)
}

// SetGlobalWorker set global worker
func SetGlobalWorker(worker *Worker) {
	w = worker
}

// GlobalWorker returns global worker
func GlobalWorker() *Worker {
	return w
}

// SetConfig set configuration
func SetConfig(conf *Config) {
	w.SetConfig(conf)
}

// Use register a new handler
func Use(handler HandlerFunc) {
	w.Use(handler)
}

// Run starts consuming messages from an input channel
func Run(ctx context.Context, input <-chan interface{}) {
	w.Run(ctx, input)
}

// CleanUp clean worker ressources
func CleanUp() {
	w.CleanUp()
}
