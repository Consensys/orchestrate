// +build unit

package configwatcher

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/configwatcher/provider/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/configwatcher/testutils"
)

func TestLoadMessage(t *testing.T) {
	assertNotRunning(t, func(t *testing.T, ctx context.Context, w *watcher, _ *mock.Provider, lis *testutils.Listener) {
		// Test Load Message
		w.loadMsg(ctx, &testutils.Message{Name: "test-provider", Conf: "test-conf"})
		assert.Equal(t, []string{"test-conf@test-provider"}, <-lis.Calls, "#1 Listener should have been called with correct value")

		// Test Load with a DeepEqual Message
		w.loadMsg(ctx, &testutils.Message{Name: "test-provider", Conf: "test-conf"})
		select {
		case <-lis.Calls:
			t.Errorf("#2 Listener should not have been called")
		default:
		}
	})
}

func TestPreLoadMessage(t *testing.T) {
	assertNotRunning(t, func(t *testing.T, ctx context.Context, w *watcher, _ *mock.Provider, _ *testutils.Listener) {
		// Test Load Message
		w.preloadMsg(ctx, &testutils.Message{Name: "test-provider", Conf: "test-conf"})
		msg := (<-w.preloadMsgs).(*testutils.Message)
		assert.Equal(t, "test-conf", msg.Conf, "#1 Message should have properly flowed")
		assert.Len(t, w.dispatchedMsgs, 1, "#1 Dispatcher should have been created")

		// Test Load Second message
		w.preloadMsg(ctx, &testutils.Message{Name: "test-provider", Conf: "test-conf2"})
		msg = (<-w.preloadMsgs).(*testutils.Message)
		assert.Equal(t, "test-conf2", msg.Conf, "#2 Message should have properly flowed")
		assert.Len(t, w.dispatchedMsgs, 1, "#2 Dispatcher should have been reused")
	})
}

func TestRun(t *testing.T) {
	assertRunning(t, func(t *testing.T, ctx context.Context, w *watcher, prvdr *mock.Provider, lis *testutils.Listener) {
		_ = prvdr.ProvideMsg(ctx, &testutils.Message{Name: "test-provider", Conf: "test-conf"})
		assert.Equal(t, []string{"test-conf@test-provider"}, <-lis.Calls, "#1 Listener should have been called with correct value")
	})
}

func assertNotRunning(t *testing.T, test func(*testing.T, context.Context, *watcher, *mock.Provider, *testutils.Listener)) {
	w, prvdr, lis := prepareWatcher()

	ctx, cancel := context.WithCancel(context.Background())
	test(t, ctx, w, prvdr, lis)
	cancel()
	w.wg.Wait()
	w.Close()
	close(lis.Calls)
}

func assertRunning(t *testing.T, test func(*testing.T, context.Context, *watcher, *mock.Provider, *testutils.Listener)) {
	w, prvdr, lis := prepareWatcher()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	go func() {
		_ = w.Run(ctx)
		close(done)
	}()

	test(t, ctx, w, prvdr, lis)

	cancel()
	<-done
	w.Close()
	close(lis.Calls)
}

func prepareWatcher() (*watcher, *mock.Provider, *testutils.Listener) {
	prvdr := mock.New()
	lis := &testutils.Listener{Calls: make(chan []string, 100)}
	listeners := []func(ctx context.Context, mergedConf interface{}) error{lis.Listen}
	w := New(&Config{0}, prvdr, testutils.MergeConfiguration, listeners).(*watcher)
	w.throttle = func(ctx context.Context, throttleDuration time.Duration, in <-chan interface{}, out chan<- interface{}) {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-in:
				out <- msg
			}
		}
	}
	return w, prvdr, lis
}
