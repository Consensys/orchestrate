package steps

import (
	"fmt"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope"
)

type tracker struct {
	// Output envelopes
	output map[string]chan *envelope.Envelope

	// Envelope that can be diagnose (last one retrieve from an out channel)
	current *envelope.Envelope
}

func newTracker() *tracker {
	t := &tracker{
		output: make(map[string]chan *envelope.Envelope),
	}
	return t
}

func (t *tracker) addOutput(key string, ch chan *envelope.Envelope) {
	t.output[key] = ch
}

func (t *tracker) get(key string, timeout time.Duration) (*envelope.Envelope, error) {
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

func (t *tracker) load(key string, timeout time.Duration) error {
	e, err := t.get(key, timeout)
	if err != nil {
		return err
	}

	// Set current envelope
	t.current = e

	return nil
}
