package infra

import (
	"github.com/Shopify/sarama"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
)

// PbToProducerMessage is an interface for a function that creates sarama producer message from a protobuffer
type PbToProducerMessage func(pb *tracepb.Trace) *sarama.ProducerMessage

// SaramaProducer allow to produce context on Kafka using Sarama
type SaramaProducer struct {
	p sarama.SyncProducer

	makeMsg PbToProducerMessage
}

// NewSaramaProducer creates a sarama producer
func NewSaramaProducer(p sarama.SyncProducer, f PbToProducerMessage) *SaramaProducer {
	return &SaramaProducer{p, f}
}

// Produce produce context
func (p *SaramaProducer) Produce(pb *tracepb.Trace) error {
	msg := p.makeMsg(pb)
	_, _, err := p.p.SendMessage(msg)
	return err
}
