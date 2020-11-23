package sarama

import (
	"context"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine"
)

type consumerGroupCtxKeyType string

const sessionCtxKey consumerGroupCtxKeyType = "session"
const claimCtxKey consumerGroupCtxKeyType = "claim"

// WithConsumerGroupSessionAndClaim attach a sarama ConsumerGroupSession & and sarama.ConsumerGroupClaim on context
func WithConsumerGroupSessionAndClaim(ctx context.Context, s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) context.Context {
	return context.WithValue(context.WithValue(ctx, sessionCtxKey, s), claimCtxKey, c)
}

// GetConsumerGroupSessionAndClaim return sarama ConsumerGroupSession & sarama ConsumerGroupClaim attached on context
func GetConsumerGroupSessionAndClaim(ctx context.Context) (sarama.ConsumerGroupSession, sarama.ConsumerGroupClaim) {
	s, ok := ctx.Value(sessionCtxKey).(sarama.ConsumerGroupSession)
	if !ok {
		return nil, nil
	}

	c, ok := ctx.Value(claimCtxKey).(sarama.ConsumerGroupClaim)
	if !ok {
		return nil, nil
	}
	return s, c
}

// EngineConsumerGroupHandler implements ConsumerGroupHandler interface
// (c.f https://godoc.org/github.com/Shopify/sarama#ConsumerGroupHandler)
//
// It uses a pkg Engine to Consumer messages
type EngineConsumerGroupHandler struct {
	engine *engine.Engine
}

// NewEngineConsumerGroupHandler creates a new EngineConsumerGroupHandler
func NewEngineConsumerGroupHandler(e *engine.Engine) *EngineConsumerGroupHandler {
	return &EngineConsumerGroupHandler{
		engine: e,
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (h *EngineConsumerGroupHandler) Setup(s sarama.ConsumerGroupSession) error {
	log.WithFields(log.Fields{
		"kafka.generation_id": s.GenerationID(),
		"kafka.member_id":     s.MemberID(),
	}).Infof("sarama: ready to consume claims %v", s.Claims())

	return nil
}

// ConsumeClaim starts a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed it finishes its processing and exits loop
//
// Make sure that you have registered the chain of HandlerFunc on context before ConsumeClaim is called
func (h *EngineConsumerGroupHandler) ConsumeClaim(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) error {
	logger := log.WithFields(log.Fields{
		"kafka.topic":     c.Topic(),
		"kafka.partition": c.Partition(),
	})

	logger.WithFields(log.Fields{
		"offset": c.InitialOffset(),
	}).Infof("sarama: start consuming claim")

	// Attach ConsumerGroupSession to context
	ctx := WithConsumerGroupSessionAndClaim(s.Context(), s, c)
	h.engine.Run(ctx, Pipe(ctx, c.Messages()))

	logger.Infof("sarama: stopped consuming claim")

	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
// but before the offsets are committed for the very last time.
func (h *EngineConsumerGroupHandler) Cleanup(s sarama.ConsumerGroupSession) error {
	h.engine.CleanUp()
	log.Infof("sarama: all claims consumed")
	return nil
}
