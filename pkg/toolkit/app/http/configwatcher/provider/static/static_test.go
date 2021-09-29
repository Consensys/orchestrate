// +build unit

package static_test

import (
	"context"
	"testing"

	"github.com/consensys/orchestrate/pkg/toolkit/app/http/configwatcher/provider"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/configwatcher/provider/static"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/configwatcher/testutils"
	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	p := static.New(&testutils.Message{Conf: "test-conf"})
	msgs := make(chan provider.Message, 1)
	defer close(msgs)

	done := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		_ = p.Provide(ctx, msgs)
		close(done)
	}()

	msg, _ := (<-msgs).(*testutils.Message)
	assert.Equal(t, "test-conf", msg.Conf, "Message should have flowed properly")
	cancel()
	<-done
}
