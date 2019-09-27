package chanregistry

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
)

// ChanRegistry keeps channel for each scenarios and topics
type ChanRegistry interface {
	NewEnvelopeChan(scenarioID, topic string) chan *envelope.Envelope
	GetEnvelopeChan(scenarioID, topic string) chan *envelope.Envelope
	CloseEnvelopeChan(scenarioID, topic string) error
}

type EnvelopeChanRegistry struct {
	// ScenarioID -> topic -> envelope chan
	EnvelopeChan map[string]map[string]chan *envelope.Envelope
}

// NewChanRegistry creates a new instance of a ChanRegistry
func NewChanRegistry() *EnvelopeChanRegistry {
	return &EnvelopeChanRegistry{
		EnvelopeChan: make(map[string]map[string]chan *envelope.Envelope),
	}
}

// NewEnvelopeChan creates a new envelope channel. Returns chan if already exist.
func (c *EnvelopeChanRegistry) NewEnvelopeChan(scenarioID, topic string) chan *envelope.Envelope {

	if c.EnvelopeChan[scenarioID] == nil {
		c.EnvelopeChan[scenarioID] = make(map[string]chan *envelope.Envelope)
	}
	if c.EnvelopeChan[scenarioID][topic] == nil {
		c.EnvelopeChan[scenarioID][topic] = make(chan *envelope.Envelope, 30)
	}
	return c.EnvelopeChan[scenarioID][topic]
}

// GetEnvelopeChan get a channel from the registry
func (c *EnvelopeChanRegistry) GetEnvelopeChan(scenarioID, topic string) chan *envelope.Envelope {
	if c.EnvelopeChan[scenarioID] == nil {
		return nil
	}
	return c.EnvelopeChan[scenarioID][topic]
}

// CloseEnvelopeChan closes a channel and delete it in the registry
func (c *EnvelopeChanRegistry) CloseEnvelopeChan(scenarioID, topic string) error {
	close(c.EnvelopeChan[scenarioID][topic])
	delete(c.EnvelopeChan[scenarioID], topic)
	return nil
}
