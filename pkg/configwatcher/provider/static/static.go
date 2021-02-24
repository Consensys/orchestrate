package static

import (
	"context"
	"fmt"

	"github.com/ConsenSys/orchestrate/pkg/configwatcher/provider"
	"github.com/ConsenSys/orchestrate/pkg/log"
)

type Provider struct {
	msg    provider.Message
	logger *log.Logger
}

func New(msg provider.Message) *Provider {
	return &Provider{
		msg:    msg,
		logger: log.NewLogger().SetComponent("configwatcher"),
	}
}

func (p *Provider) Provide(ctx context.Context, msgs chan<- provider.Message) error {
	p.logger.WithField("provider", fmt.Sprintf("%T", p)).
		Debug("start providing")
	msgs <- p.msg
	return nil
}
