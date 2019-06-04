package chanregistry

import (
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
)

type ChanRegistry struct {
	EnvelopeChan map[string]map[string]chan *envelope.Envelope
}

func NewChanRegistry() *ChanRegistry {
	return &ChanRegistry{
		EnvelopeChan: make(map[string]map[string]chan *envelope.Envelope),
	}
}

func (c *ChanRegistry) NewEnvelopeChan(testID, topic string) chan *envelope.Envelope {

	if c.EnvelopeChan[testID] == nil {
		c.EnvelopeChan[testID] = make(map[string]chan *envelope.Envelope)
	}
	if c.EnvelopeChan[testID][topic] == nil {
		c.EnvelopeChan[testID][topic] = make(chan *envelope.Envelope)
	}
	return c.EnvelopeChan[testID][topic]
}

func (c *ChanRegistry) SetEnvelopeChan(testID, topic string, e chan *envelope.Envelope) error {

	if c.EnvelopeChan[testID] == nil {
		c.EnvelopeChan[testID] = make(map[string]chan *envelope.Envelope)
	}

	if c.EnvelopeChan[testID][topic] != nil {
		return fmt.Errorf("envelope registry: Channel already set for %s - %s", testID, topic)
	}
	c.EnvelopeChan[testID][topic] = e
	return nil
}

func (c *ChanRegistry) GetEnvelopeChan(testID, topic string) chan *envelope.Envelope {
	return c.EnvelopeChan[testID][topic]
}

func (c *ChanRegistry) CloseEnvelopeChan(testID, topic string) error {
	close(c.EnvelopeChan[testID][topic])
	delete(c.EnvelopeChan[testID], topic)
	return nil
}
