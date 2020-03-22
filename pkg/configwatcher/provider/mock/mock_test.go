// +build unit

package mock_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/testutils"
)

func TestProvider(t *testing.T) {
	p := mock.New()
	msgs := make(chan provider.Message, 1)
	defer close(msgs)

	done := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		_ = p.Provide(ctx, msgs)
		close(done)
	}()

	_ = p.ProvideMsg(ctx, &testutils.Message{Conf: "test-conf"})

	msg, _ := (<-msgs).(*testutils.Message)
	assert.Equal(t, "test-conf", msg.Conf, "Message should have flowed properly")
	cancel()
	<-done
}
