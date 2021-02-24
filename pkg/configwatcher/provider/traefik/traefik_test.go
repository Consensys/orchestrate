// +build unit

package traefik_test

import (
	"context"
	"testing"

	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/safe"
	"github.com/stretchr/testify/assert"
	"github.com/ConsenSys/orchestrate/pkg/configwatcher/provider"
	"github.com/ConsenSys/orchestrate/pkg/configwatcher/provider/traefik"
	"github.com/ConsenSys/orchestrate/pkg/http/config/dynamic"
)

type MockTraefikProvider struct {
	confChan chan traefikdynamic.Message
}

func (p *MockTraefikProvider) Init() error {
	return nil
}

func (p *MockTraefikProvider) Provide(confChan chan<- traefikdynamic.Message, pool *safe.Pool) error {
	confChan <- <-p.confChan
	return nil
}

func (p *MockTraefikProvider) ProvideMsg(msg traefikdynamic.Message) {
	p.confChan <- msg
}

func TestTraefik(t *testing.T) {
	mock := &MockTraefikProvider{confChan: make(chan traefikdynamic.Message, 1)}

	ctx, cancel := context.WithCancel(context.Background())
	prvdr := traefik.New(mock, safe.NewPool(ctx))

	msgs := make(chan provider.Message, 1)
	done := make(chan struct{})
	go func() {
		_ = prvdr.Provide(ctx, msgs)
		close(done)
	}()

	traefikMsg := traefikdynamic.Message{ProviderName: "test-provider"}
	mock.ProvideMsg(traefikMsg)

	msg := (<-msgs).(*dynamic.Message)
	assert.Equal(t, "test-provider", msg.ProviderName(), "Message should have flowed correctly")

	cancel()
	<-done
}
