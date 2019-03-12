package utils

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// SignalListener listen to signals and trigger callbacks
type SignalListener struct {
	signals chan os.Signal

	closed    chan struct{}
	closeOnce *sync.Once

	cb func(signal os.Signal)
}

// NewSignalListener creates a new SignalListener
func NewSignalListener(cb func(os.Signal)) *SignalListener {
	l := &SignalListener{
		signals:   make(chan os.Signal, 3),
		closed:    make(chan struct{}),
		closeOnce: &sync.Once{},
		cb:        cb,
	}

	go l.listen()

	return l
}

// Close signal listener
func (l *SignalListener) Close() {
	l.closeOnce.Do(func() {
		close(l.closed)
	})
}

// Listen start Listening to signals
func (l *SignalListener) listen() {
	// Redirect signals
	signal.Notify(l.signals)
signalLoop:
	for {
		select {
		case signal := <-l.signals:
			l.processSignal(signal)
		case <-l.closed:
			break signalLoop
		}
	}
}

func (l *SignalListener) processSignal(signal os.Signal) {
	switch signal {
	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
		log.Warnf("signal: %q intercepted", signal.String())
		l.cb(signal)
	default:
		log.Fatalf("signal: unknown signal %q intercepted, exit now", signal.String())
	}
}
