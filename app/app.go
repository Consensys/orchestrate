package app

import (
	"context"
	"net/http"
	"sync/atomic"

	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-crafter.git/app/infra"
)

// App is main application object
type App struct {
	server *http.Server
	infra  *infra.Infra

	saramaConsumerGroup sarama.ConsumerGroup
	saramaHandler       sarama.ConsumerGroupHandler

	ctx    context.Context
	cancel func()

	ready *atomic.Value
	done  chan struct{}
}

// New creates a new application
func New(ctx context.Context) *App {
	// We set a cancellable context so we can possibly abort application form within the application
	ctx, cancel := context.WithCancel(ctx)
	app := &App{
		done:   make(chan struct{}),
		infra:  infra.NewInfra(),
		ctx:    ctx,
		cancel: cancel,
		ready:  &atomic.Value{},
	}

	// App is not yet ready
	app.ready.Store(false)

	// Initialize app
	app.init()

	// We indicate that application is ready
	app.ready.Store(true)

	return app
}

func (app *App) init() {
	// Initialize application
	initServer(app)
	app.infra.Init()
	initConsumerGroup(app)
}

// Ready indicate if app is ready
func (app *App) Ready() bool {
	select {
	case <-app.ctx.Done():
		return false
	default:
		return app.ready.Load().(bool)
	}
}

// Run application
func (app *App) Run() {
	// Start consumer group
	app.saramaConsumerGroup.Consume(
		app.ctx,
		[]string{viper.GetString("worker.in")},
		app.saramaHandler,
	)

	// We close infrastructure
	app.infra.Close()

	// We indicate that application has stopped running
	close(app.done)
}

// Done return a channel indicating that application is done running
func (app *App) Done() <-chan struct{} {
	return app.done
}
