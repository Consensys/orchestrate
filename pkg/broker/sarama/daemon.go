package sarama

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/hashicorp/go-multierror"
)

type ConsumerDaemon struct {
	client   sarama.Client
	producer sarama.SyncProducer
	group    sarama.ConsumerGroup

	topics  []string
	handler sarama.ConsumerGroupHandler
}

func NewConsumerDaemon(
	client sarama.Client,
	producer sarama.SyncProducer,
	group sarama.ConsumerGroup,
	topics []string,
	handler sarama.ConsumerGroupHandler,
) *ConsumerDaemon {
	return &ConsumerDaemon{
		client:   client,
		producer: producer,
		group:    group,
		topics:   topics,
		handler:  handler,
	}
}

func (d *ConsumerDaemon) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			err := d.group.Consume(ctx, d.topics, d.handler)
			if err != nil {
				return err
			}
		}
	}
}

func (d *ConsumerDaemon) Close() error {
	gr := &multierror.Group{}
	gr.Go(d.producer.Close)
	gr.Go(d.group.Close)
	rerr := gr.Wait()

	err := d.client.Close()
	if err != nil {
		rerr = multierror.Append(rerr, err)
	}

	return rerr.ErrorOrNil()
}
