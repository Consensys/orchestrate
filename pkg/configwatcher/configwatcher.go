package configwatcher

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/configwatcher/provider"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
)

//go:generate mockgen -source=configwatcher.go -destination=mock/mock.go -package=mock

type Listener func(context.Context, interface{}) error

type Watcher interface {
	AddListener(Listener)
	Run(context.Context) error
	Close() error
}

type watcher struct {
	cfg *Config

	provider provider.Provider

	merge func(map[string]interface{}) interface{}

	listeners []func(context.Context, interface{}) error

	wg             *sync.WaitGroup
	inputMsgs      chan provider.Message
	preloadMsgs    chan provider.Message
	dispatchedMsgs map[string]chan provider.Message

	currentConfigs map[string]interface{}

	throttle func(ctx context.Context, throttleDuration time.Duration, in <-chan interface{}, out chan<- interface{})
}

func New(
	cfg *Config,
	prvdr provider.Provider,
	merge func(map[string]interface{}) interface{},
	listeners []func(context.Context, interface{}) error,
) Watcher {
	return &watcher{
		cfg:            cfg,
		provider:       prvdr,
		merge:          merge,
		listeners:      listeners,
		wg:             &sync.WaitGroup{},
		inputMsgs:      make(chan provider.Message, 100),
		preloadMsgs:    make(chan provider.Message, 100),
		dispatchedMsgs: make(map[string]chan provider.Message),
		currentConfigs: make(map[string]interface{}),
		throttle:       Throttle,
	}
}

func (w *watcher) AddListener(listener Listener) {
	w.listeners = append(w.listeners, listener)
}

func (w *watcher) Run(ctx context.Context) (err error) {
	utils.InParallel(
		func() { w.watchInputMsgs(ctx) },
		func() { w.watchPreloadedMsgs(ctx) },
		func() { err = w.provider.Provide(ctx, w.inputMsgs) },
	)
	w.wg.Wait()
	return
}

func (w *watcher) watchInputMsgs(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-w.inputMsgs:
			if !ok {
				return
			}
			w.preloadMsg(ctx, msg)
		}
	}
}

func (w *watcher) watchPreloadedMsgs(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-w.preloadMsgs:
			if !ok {
				return
			}
			w.loadMsg(ctx, msg)
		}
	}
}

func (w *watcher) loadMsg(ctx context.Context, msg provider.Message) {
	// Make sure that configuration has been updated
	currentConfig, ok := w.currentConfigs[msg.ProviderName()]
	if ok && reflect.DeepEqual(currentConfig, msg.Configuration()) {
		return
	}

	// We got a new configuration so we update current config
	w.currentConfigs[msg.ProviderName()] = msg.Configuration()
	log.FromContext(ctx).
		WithField(log.ProviderName, msg.ProviderName()).
		Infof("got new configuration")

	// Call listeners
	conf := w.merge(w.currentConfigs)

	for _, listener := range w.listeners {
		err := listener(ctx, conf)
		if err != nil {
			log.FromContext(ctx).
				WithField(log.ProviderName, msg.ProviderName()).
				WithError(err).
				Warning("config listener error")
		}
	}
}

func (w *watcher) preloadMsg(ctx context.Context, msg provider.Message) {
	if reflect.ValueOf(msg.Configuration()).IsZero() {
		log.FromContext(ctx).Infof("Skipping empty Configuration for provider %s", msg.ProviderName)
		return
	}

	msgs, ok := w.dispatchedMsgs[msg.ProviderName()]
	if !ok {
		msgs = make(chan provider.Message)
		w.dispatchedMsgs[msg.ProviderName()] = msgs
		w.wg.Add(1)
		go func() {
			out := pipeOut(w.preloadMsgs)
			w.throttle(ctx, w.cfg.ProvidersThrottleDuration, pipeIn(msgs), out)
			close(out)
			w.wg.Done()
		}()
	}
	msgs <- msg
}

func (w *watcher) Close() error {
	close(w.inputMsgs)
	close(w.preloadMsgs)
	for _, v := range w.dispatchedMsgs {
		close(v)
	}
	return nil
}

func pipeIn(in <-chan provider.Message) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		for msg := range in {
			out <- msg
		}
		close(out)
	}()
	return out
}

func pipeOut(out chan<- provider.Message) chan<- interface{} {
	in := make(chan interface{})
	go func() {
		for msg := range in {
			out <- msg.(provider.Message)
		}
	}()

	return in
}
