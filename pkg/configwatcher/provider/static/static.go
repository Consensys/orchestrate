package static

import (
	"context"
	"fmt"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider"
)

type Provider struct {
	msg provider.Message
}

func New(msg provider.Message) *Provider {
	return &Provider{
		msg: msg,
	}
}

func (p *Provider) Provide(ctx context.Context, msgs chan<- provider.Message) error {
	log.FromContext(ctx).WithField("provider", fmt.Sprintf("%T", p)).Infof("start providing")
	msgs <- p.msg
	return nil
}
