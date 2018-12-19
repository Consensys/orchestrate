package infra

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
)

// SaramaProcessor interface to process a kafka message
type SaramaProcessor interface {
	ProcessMessage(message *sarama.ConsumerMessage)
	ProcessError(err error)
}

// Listener consumes and processes messages from a Kafka partition partiion
type Listener struct {
	w         *SaramaWorker
	pc        sarama.PartitionConsumer
	processor SaramaProcessor
	stop      chan struct{}
}

//NewListener creates a new listener on sarama partition
func NewListener(w *SaramaWorker, pc sarama.PartitionConsumer, processor SaramaProcessor) *Listener {
	return &Listener{
		w,
		pc,
		processor,
		make(chan struct{}, 1),
	}
}

// Listen start listener
func (l *Listener) Listen() {
	fmt.Printf("Listener %p: main loop starts...\n", l)
	for {
		select {
		case <-l.stop:
			// In case of graceful interuption we close the listener
			l.Close()
			break
		case err := <-l.pc.Errors():
			l.processor.ProcessError(err)
		case msg := <-l.pc.Messages():
			l.processor.ProcessMessage(msg)
		}
	}
}

// Stop gracefully interupts listener (interuption appends asynchronously)
func (l *Listener) Stop() {
	// Send message into channel so we have no problem of concurrency
	l.stop <- struct{}{}
}

// Close close listener
func (l *Listener) Close() {
	fmt.Printf("Listener %p: received stop signal...\n", l)

	// Notify worker that listener has stopped
	if l.w != nil {
		l.w.Reap(l)
	}
}

// SaramaWorker is a Kafka worker based on sarama
type SaramaWorker struct {
	consumer  sarama.Consumer
	signals   chan os.Signal
	listeners map[*Listener]bool
	new, dead chan *Listener
}

var (
	maxListener = 5
)

// NewSaramaWorker creates a new sarama worker
func NewSaramaWorker(c sarama.Consumer) *SaramaWorker {
	return &SaramaWorker{
		c,
		make(chan os.Signal, 3),
		make(map[*Listener]bool),
		make(chan *Listener, maxListener),
		make(chan *Listener, maxListener),
	}
}

// InitSignals redirect signals
func (w *SaramaWorker) InitSignals() {
	signal.Notify(w.signals)
}

// ProcessSignal process signals
func (w *SaramaWorker) ProcessSignal(signal os.Signal) {
	switch signal {
	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
		// Gracefully stop
		fmt.Printf("Worker: gracefully stops...\n")
		w.Stop()
	default:
		// Exit
		fmt.Printf("Worker: unknown signal exits...\n")
		w.Exit(1)
	}
}

// Reap queue a listener to be removed
func (w *SaramaWorker) Reap(l *Listener) {
	// Send message into channel so we have no problem of concurrency
	w.dead <- l
}

func (w *SaramaWorker) addListener(l *Listener) error {
	if _, ok := w.listeners[l]; ok {
		return fmt.Errorf("Listener already exists")
	}

	w.listeners[l] = true
	fmt.Printf("Worker: added listener %p\n", l)
	return nil
}

func (w *SaramaWorker) removeListener(l *Listener) {
	delete(w.listeners, l)
	fmt.Printf("Worker: removed listener %p\n", l)
}

// Subscribe allows to consume and process messages from a partition and return a function to stop listening to the event
func (w *SaramaWorker) Subscribe(topic string, partition int32, offset int64, processor SaramaProcessor) error {
	pc, err := w.consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		return err
	}

	l := NewListener(w, pc, processor)
	w.new <- l
	return nil
}

// Run start listeners
func (w *SaramaWorker) Run() {
	w.InitSignals()
	fmt.Printf("Worker: main loop starts\n")
	for {
		select {
		case signal := <-w.signals:
			fmt.Printf("Worker: got signal %v\n", signal)
			w.ProcessSignal(signal)
		case l := <-w.new:
			fmt.Printf("Worker: got new listener %p\n", l)
			err := w.addListener(l)
			if err != nil {
			} else {
				fmt.Printf("Worker: start listener %p\n", l)
				// Start listener in a dedicated go routine
				go l.Listen()
			}
		case l := <-w.dead:
			fmt.Printf("Worker: get dead listener %p\n", l)
			w.removeListener(l)
			if len(w.listeners) == 0 {
				// We exit as all listeners are stoped
				w.Exit(0)
			}
		}
	}
}

// Exit exit process
func (w *SaramaWorker) Exit(status int) {
	fmt.Printf("Worker: exiting with status %v\n", status)
	os.Exit(status)
}

// Stop gracefully stops worker
func (w *SaramaWorker) Stop() {
	// Gracefully stops every listener
	for l := range w.listeners {
		l.Stop()
	}
}
