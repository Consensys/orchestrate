package app

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// App boilerplate application
type App struct {
	dying, done chan struct{}
	closeOnce   *sync.Once
}

// New creates a new boilerplate application
func New() *App {
	return &App{
		dying:     make(chan struct{}),
		done:      make(chan struct{}),
		closeOnce: &sync.Once{},
	}
}

// Start application
func (app *App) Start() {
	log.Infof("boilerplate: starts...")
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
	close(app.done)
}

// Close application
func (app *App) Close() {
	app.closeOnce.Do(func() {
		close(app.dying)
		log.Infof("boilerplate: closing...")
	})
}

// Done return a channel indicating if application has stopped running
func (app *App) Done() <-chan struct{} {
	return app.done
}
