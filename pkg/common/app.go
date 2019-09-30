package common

import (
	"fmt"
	"sync/atomic"
)

// App is an application structure that expose Ready method
// TODO: App pattern is a v0.1 that is functional and should be used until we mature the overall Wire pattern
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
func (app *App) IsReady() error {
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

// SetReady set readiness status to true
func (app *App) SetReady(ready bool) {
	app.ready.Store(ready)
}
