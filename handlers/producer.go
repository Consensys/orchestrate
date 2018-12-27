package handlers

import (
	"github.com/Shopify/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
)

// ContextProducer produces a context in another service typically a Kafka queue
type ContextProducer interface {
	Produce(ctx *infra.Context) error
}

// CtxToProducerMessage is an interface for a function that creates sarama message to produce from a Context
type CtxToProducerMessage func(ctx *infra.Context) *sarama.ProducerMessage

// SaramaProducer allow to produce context on Kafka using Sarama
type SaramaProducer struct {
	p sarama.SyncProducer

	makeMsg CtxToProducerMessage
}

// NewSaramaProducer creates a sarama producer
func NewSaramaProducer(p sarama.SyncProducer, f CtxToProducerMessage) *SaramaProducer {
	return &SaramaProducer{p, f}
}

// Produce produce context
func (p *SaramaProducer) Produce(ctx *infra.Context) error {
	msg := p.makeMsg(ctx)
	_, _, err := p.p.SendMessage(msg)
	// TODO: for later tracking we could register partition/offset info on context
	return err
}

// Producer creates a producer handler
func Producer(p ContextProducer) infra.HandlerFunc {
	return func(ctx *infra.Context) {
		p.Produce(ctx)
	}
}
