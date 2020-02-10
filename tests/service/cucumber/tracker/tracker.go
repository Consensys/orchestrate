package tracker

import (
	"fmt"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/tx"
)

type Tracker struct {
	// Output envelopes
	output map[string]chan *tx.Builder

	// Builder that can be diagnosed (last one retrieved from an out channel)
	Current *tx.Builder
}

func NewTracker() *Tracker {
	t := &Tracker{
		output: make(map[string]chan *tx.Builder),
	}
	return t
}

func (t *Tracker) AddOutput(key string, ch chan *tx.Builder) {
	t.output[key] = ch
}

func (t *Tracker) get(key string, timeout time.Duration) (*tx.Builder, error) {
	ch, ok := t.output[key]
	if !ok {
		return nil, fmt.Errorf("output %q not tracked", key)
	}

	select {
	case e := <-ch:
		return e, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("no envelope available in output %q", key)
	}
}

func (t *Tracker) Load(key string, timeout time.Duration) error {
	e, err := t.get(key, timeout)
	if err != nil {
		return err
	}

	// Set Current envelope
	t.Current = e

	return nil
}
