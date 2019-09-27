package app

import (
	"fmt"
	"sync/atomic"
)

// App is an application structure that expose Ready method
type App struct {
	ready   *atomic.Value
	closing chan struct{}
}

// NewApp creates new app
func NewApp() *App {
	app := &App{
		ready:   &atomic.Value{},
		closing: make(chan struct{}),
	}
	app.ready.Store(false)
	return app
}

// Ready indicates if application is ready
func (app *App) Ready() error {
	select {
	case <-app.closing:
		return fmt.Errorf("app is closing")
	default:
		if !app.ready.Load().(bool) {
			return fmt.Errorf("app is not ready")
		}
		return nil
	}
}
