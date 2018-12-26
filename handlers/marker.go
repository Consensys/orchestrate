package handlers

import (
	"sync"

	"github.com/Shopify/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
)

// OffsetMarker is an interface to mark offsets of message that have been consumed
type OffsetMarker interface {
	Mark(ctx *infra.Context)
}

// SimpleSaramaOffsetMarker marks message each time it encounters a message with higher offset than the last one marked
type SimpleSaramaOffsetMarker struct {
	s sarama.ConsumerGroupSession

	mux        *sync.Mutex
	lastOffset int64
}

// NewSimpleSaramaOffsetMarker creates a new simple marker
func NewSimpleSaramaOffsetMarker(s sarama.ConsumerGroupSession) *SimpleSaramaOffsetMarker {
	return &SimpleSaramaOffsetMarker{s, &sync.Mutex{}, -1}
}

// Mark mark an offset each time it encounter a message with higher offset than the last one marked
func (c *SimpleSaramaOffsetMarker) Mark(ctx *infra.Context) {
	c.mux.Lock()
	msg := ctx.Msg.(*sarama.ConsumerMessage)
	if msg.Offset > c.lastOffset {
		// Current message has a larger offset than the last one marked
		c.lastOffset = msg.Offset
		c.mux.Unlock()
		// Mark message
		c.s.MarkMessage(msg, "")
		return
	}
	c.mux.Unlock()
}

// Marker creates an handler that mark offsets
func Marker(offset OffsetMarker) infra.HandlerFunc {
	return func(ctx *infra.Context) {
		offset.Mark(ctx)
	}
}
