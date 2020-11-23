package consumer

import (
	"github.com/Shopify/sarama"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine"
)

type EmbeddingConsumerGroupHandler struct {
	engine  *broker.EngineConsumerGroupHandler
	isReady chan bool
}

func NewEmbeddingConsumerGroupHandler(e *engine.Engine) *EmbeddingConsumerGroupHandler {
	return &EmbeddingConsumerGroupHandler{
		engine:  broker.NewEngineConsumerGroupHandler(e),
		isReady: make(chan bool, 1),
	}
}

func (h *EmbeddingConsumerGroupHandler) Setup(s sarama.ConsumerGroupSession) error {
	err := h.engine.Setup(s)
	h.isReady <- true
	return err
}

func (h *EmbeddingConsumerGroupHandler) IsReady() chan bool {
	return h.isReady
}

func (h *EmbeddingConsumerGroupHandler) ConsumeClaim(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) error {
	return h.engine.ConsumeClaim(s, c)
}

func (h *EmbeddingConsumerGroupHandler) Cleanup(s sarama.ConsumerGroupSession) error {
	return h.engine.Cleanup(s)
}
