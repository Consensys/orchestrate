// +build unit

package aggregator_test

import (
	"context"
	"testing"

	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/configwatcher/provider"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/configwatcher/provider/aggregator"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/configwatcher/provider/mock"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/configwatcher/testutils"
	"github.com/stretchr/testify/assert"
)

func TestAggregator(t *testing.T) {
	p1 := mock.New()
	p2 := mock.New()

	p := aggregator.New()
	p.AddProvider(p1)
	p.AddProvider(p2)

	msgs := make(chan provider.Message, 2)
	defer close(msgs)

	done := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		_ = p.Provide(ctx, msgs)
		close(done)
	}()

	_ = p1.ProvideMsg(ctx, &testutils.Message{Conf: "test-conf1"})
	_ = p2.ProvideMsg(ctx, &testutils.Message{Conf: "test-conf2"})

	msg, _ := (<-msgs).(*testutils.Message)
	assert.Equal(t, "test-conf1", msg.Conf, "#1 Message should have flowed properly")
	msg, _ = (<-msgs).(*testutils.Message)
	assert.Equal(t, "test-conf2", msg.Conf, "#2 Message should have flowed properly")

	cancel()
	<-done
}
