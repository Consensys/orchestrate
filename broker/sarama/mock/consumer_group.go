package mock

import (
	"context"
	"fmt"
	"sync"

	"github.com/Shopify/sarama"
	uuid "github.com/satori/go.uuid"
)

// ConsumerGroup is a mock implementation of a sarama.ConsumerGroup
type ConsumerGroup struct {
	errors chan error

	closeOnce *sync.Once
	msgs      map[string]map[int32][]*sarama.ConsumerMessage
	group     string
}

// NewConsumerGroup creates a new ConsumerGroup mock
func NewConsumerGroup(group string, msgs map[string]map[int32][]*sarama.ConsumerMessage) *ConsumerGroup {
	return &ConsumerGroup{
		group:     group,
		errors:    make(chan error),
		closeOnce: &sync.Once{},
		msgs:      msgs,
	}
}

// Consume joins a cluster of consumers for a given list of topics and
// starts a blocking ConsumerGroupSession through the ConsumerGroupHandler.
func (g *ConsumerGroup) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	// Compute claims
	claims := make(map[string][]int32)
	for _, topic := range topics {
		if partitions, ok := g.msgs[topic]; ok {
			claims[topic] = []int32{}
			for partition := range partitions {
				claims[topic] = append(claims[topic], partition)
			}
		}
	}

	// Create Consumer Group
	s := NewConsumerGroupSession(ctx, g.group, claims)

	// Call SetUp hook
	handler.Setup(s)

	wg := &sync.WaitGroup{}
	for topic, partitions := range s.Claims() {
		for _, partition := range partitions {
			c := NewConsumerGroupClaim(topic, partition, s.LastMarkedOffset(topic, partition).Offset)

			// Feed mock claim with messages
			go func(c *ConsumerGroupClaim) {
				for _, msg := range g.msgs[c.Topic()][c.Partition()] {
					c.ExpectMessage(msg)
				}
			}(c)

			wg.Add(1)
			go func(c *ConsumerGroupClaim) {
				// ConsumeClaim in dedicated go routine
				handler.ConsumeClaim(s, c)
				wg.Done()
			}(c)
		}
	}
	wg.Wait()

	// Call CleanUp loop
	handler.Cleanup(s)
	return nil
}

// Errors returns a read channel of errors that occurred during the consumer life-cycle.
func (g *ConsumerGroup) Errors() <-chan error {
	return g.errors
}

// Close stops the ConsumerGroup and detaches any running sessions.
func (g *ConsumerGroup) Close() error {
	g.closeOnce.Do(func() {
		close(g.errors)
	})
	return nil
}

var (
	mux          = &sync.Mutex{}
	generationID = make(map[string]int32)
)

func nextGenerationID(group string) int32 {
	mux.Lock()
	if _, ok := generationID[group]; !ok {
		generationID[group] = 1
	} else {
		generationID[group]++
	}
	defer mux.Unlock()
	return generationID[group]
}

// MarkedOffset represent a marked offset
type MarkedOffset struct {
	Offset   int64
	Metadata string
}

// ConsumerGroupSession is a mock implementation of a sarama.ConsumerGroupSession
type ConsumerGroupSession struct {
	generationID int32
	memberID     string
	claims       map[string][]int32

	mux          *sync.Mutex
	MarkedOffset map[string]map[int32][]*MarkedOffset

	ctx context.Context
}

// NewConsumerGroupSession creates a new ConsumerGroupSession
func NewConsumerGroupSession(ctx context.Context, group string, claims map[string][]int32) *ConsumerGroupSession {
	return &ConsumerGroupSession{
		generationID: nextGenerationID(group),
		memberID:     fmt.Sprintf("%v-%v", "mock", uuid.NewV4().String()),
		mux:          &sync.Mutex{},
		MarkedOffset: make(map[string]map[int32][]*MarkedOffset),
		claims:       claims,
		ctx:          ctx,
	}
}

// Claims returns information about the claimed partitions by topic.
func (s *ConsumerGroupSession) Claims() map[string][]int32 {
	return s.claims
}

// MemberID returns the cluster member ID.
func (s *ConsumerGroupSession) MemberID() string {
	return s.memberID
}

// GenerationID returns the current generation ID.
func (s *ConsumerGroupSession) GenerationID() int32 {
	return s.generationID
}

// Context returns the session context.
func (s *ConsumerGroupSession) Context() context.Context {
	return s.ctx
}

func (s *ConsumerGroupSession) getOffsets(topic string, partition int32) []*MarkedOffset {
	if _, ok := s.MarkedOffset[topic]; !ok {
		s.MarkedOffset[topic] = make(map[int32][]*MarkedOffset)
	}

	if _, ok := s.MarkedOffset[topic][partition]; !ok {
		s.MarkedOffset[topic][partition] = make([]*MarkedOffset, 0)
	}

	return s.MarkedOffset[topic][partition]
}

// MarkOffset marks the provided offset, alongside a metadata string
func (s *ConsumerGroupSession) MarkOffset(topic string, partition int32, offset int64, metadata string) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if mOffset := s.lastMarkedOffset(topic, partition); mOffset != nil && mOffset.Offset >= offset {
		return
	}

	offsets := s.getOffsets(topic, partition)
	s.MarkedOffset[topic][partition] = append(offsets, &MarkedOffset{offset, metadata})
}

// MarkMessage marks a message as consumed.
func (s *ConsumerGroupSession) MarkMessage(msg *sarama.ConsumerMessage, metadata string) {
	s.MarkOffset(msg.Topic, msg.Partition, msg.Offset+1, metadata)
}

func (s *ConsumerGroupSession) lastMarkedOffset(topic string, partition int32) *MarkedOffset {
	if len(s.getOffsets(topic, partition)) == 0 {
		return &MarkedOffset{}
	}
	return s.getOffsets(topic, partition)[len(s.getOffsets(topic, partition))-1]
}

// LastMarkedOffset return last marked offset for a pair topic/partition
func (s *ConsumerGroupSession) LastMarkedOffset(topic string, partition int32) *MarkedOffset {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.lastMarkedOffset(topic, partition)
}

// ResetOffset resets to the provided offset, alongside a metadata string
func (s *ConsumerGroupSession) ResetOffset(topic string, partition int32, offset int64, metadata string) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if mOffset := s.lastMarkedOffset(topic, partition); mOffset != nil && mOffset.Offset < offset {
		return
	}

	for i, mOffset := range s.getOffsets(topic, partition) {
		if mOffset.Offset >= offset {
			s.MarkedOffset[topic][partition] = s.MarkedOffset[topic][partition][:i]
			break
		}
	}

	s.MarkedOffset[topic][partition] = append(s.MarkedOffset[topic][partition], &MarkedOffset{offset, metadata})
}

// ConsumerGroupClaim is a mock implementation of a sarama.ConsumerGroupClaim
type ConsumerGroupClaim struct {
	topic         string
	partition     int32
	initialOffset int64

	msgs          chan *sarama.ConsumerMessage
	mux           *sync.Mutex
	highWaterMark int64
}

// NewConsumerGroupClaim creates a new ConsumerGroupClaim
func NewConsumerGroupClaim(topic string, partition int32, initialOffset int64) *ConsumerGroupClaim {
	return &ConsumerGroupClaim{
		topic:         topic,
		partition:     partition,
		initialOffset: initialOffset,
		msgs:          make(chan *sarama.ConsumerMessage),
		mux:           &sync.Mutex{},
		highWaterMark: 0,
	}
}

// Topic returns the consumed topic name.
func (c ConsumerGroupClaim) Topic() string {
	return c.topic
}

// Partition returns the consumed partition.
func (c ConsumerGroupClaim) Partition() int32 {
	return c.partition
}

// InitialOffset returns the initial offset that was used as a starting point for this claim.
func (c ConsumerGroupClaim) InitialOffset() int64 {
	return c.initialOffset
}

// Messages returns the read channel for the messages that are returned by
func (c ConsumerGroupClaim) Messages() <-chan *sarama.ConsumerMessage {
	return c.msgs
}

// HighWaterMarkOffset returns the high water mark offset of the partition,
// i.e. the offset that will be used for the next message that will be produced.
// You can use this to determine how far behind the processing is.
func (c ConsumerGroupClaim) HighWaterMarkOffset() int64 {
	return c.highWaterMark
}

// ExpectMessage add a message that will be consumed
func (c ConsumerGroupClaim) ExpectMessage(msg *sarama.ConsumerMessage) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.msgs <- msg
}
