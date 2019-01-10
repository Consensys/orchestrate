package infra

import (
	"fmt"
	"sync"

	"github.com/Shopify/sarama"
)

// SimpleSaramaOffsetMarker marks message each time it encounters a message with higher offset than the last one marked
type SimpleSaramaOffsetMarker struct {
	s sarama.ConsumerGroupSession

	mux        *sync.Mutex
	lastOffset int64
}

// Mark mark an offset each time it encounter a message with higher offset than the last one marked
func (c *SimpleSaramaOffsetMarker) Mark(msg interface{}) error {
	c.mux.Lock()
	cast, ok := msg.(*sarama.ConsumerMessage)
	if !ok {
		// Format is incorrect
		c.mux.Unlock()
		return fmt.Errorf("Message does not match expected format")
	}
	if cast.Offset > c.lastOffset {
		// Current message has a larger offset than the last one marked
		c.lastOffset = cast.Offset
		c.mux.Unlock()
		// Mark message
		c.s.MarkMessage(cast, "")
		return nil
	}

	c.mux.Unlock()
	return nil
}

// NewSimpleSaramaOffsetMarker creates a new simple marker
func NewSimpleSaramaOffsetMarker(s sarama.ConsumerGroupSession) *SimpleSaramaOffsetMarker {
	return &SimpleSaramaOffsetMarker{s, &sync.Mutex{}, -1}
}
