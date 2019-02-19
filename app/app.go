package app

import (
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)

// App boilerplate application
type App struct {
	dying, done chan struct{}
	closeOnce   *sync.Once

	server *http.Server
	ready  *atomic.Value
}

// New creates a new application
func New() *App {
	app := &App{
		dying:     make(chan struct{}),
		done:      make(chan struct{}),
		closeOnce: &sync.Once{},
		ready:     &atomic.Value{},
	}

	// We are not yet ready
	app.ready.Store(false)

	// Initialize app
	app.init()

	return app
}

func (app *App) init() {
	// Initialize server on application
	initServer(app)
}

// Run application
func (app *App) Run() {
	// Run main loop
	log.Infof("boilerplate: starts running...")
	go app.run()
}

func (app *App) run() {
	// Indicate that app is ready
	app.ready.Store(true)

	ticker := time.NewTicker(time.Second)
appLoop:
	for {
		select {
		case t := <-ticker.C:
			log.Infof("boilerplate: %v", t)
		case <-app.dying:
			break appLoop
		}
	}
	ticker.Stop()
	close(app.done)
}

// Closed return whether the application has been closed
func (app *App) Closed() bool {
	select {
	case <-app.dying:
		return true
	default:
		return false
	}
}

// Ready indicate if app is ready
func (app *App) Ready() bool {
	return app.ready.Load().(bool)
}

// Close application
func (app *App) Close() {
	app.closeOnce.Do(func() {
		// Indicate that app is no more ready
		app.ready.Store(false)
		close(app.dying)
		log.Infof("boilerplate: closing...")
	})
}

// Done return a channel indicating if application has stopped running
func (app *App) Done() <-chan struct{} {
	return app.done
}
