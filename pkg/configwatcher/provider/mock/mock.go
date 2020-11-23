package mock

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/configwatcher/provider"
)

type Provider struct {
	msgs chan provider.Message
}

func New() *Provider {
	return &Provider{
		msgs: make(chan provider.Message),
	}
}

func (p *Provider) Provide(ctx context.Context, msgs chan<- provider.Message) error {
	select {
	case <-ctx.Done():
		close(p.msgs)
	case msg := <-p.msgs:
		msgs <- msg
	}
	return nil
}

func (p *Provider) ProvideMsg(ctx context.Context, msg provider.Message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case p.msgs <- msg:
		return nil
	}
}
