// +build unit

package grpc

import (
	"context"
	"net"
	"testing"
	"time"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"github.com/stretchr/testify/require"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/generic"
	"google.golang.org/grpc"
)

func TestEntryPoint(t *testing.T) {
	cfg := &traefikstatic.EntryPoint{
		Address: "127.0.0.1:0",
		Transport: &traefikstatic.EntryPointsTransport{
			RespondingTimeouts: &traefikstatic.RespondingTimeouts{},
			LifeCycle:          &traefikstatic.LifeCycle{},
		},
	}

	ep := NewEntryPoint("", cfg, grpc.NewServer(), generic.NewTCP())
	done := make(chan struct{})
	go func() {
		_ = ep.ListenAndServe(context.Background())
		close(done)
	}()

	// Wait a few millisecond for server to start
	time.Sleep(50 * time.Millisecond)

	_, err := net.Dial("tcp", ep.Addr())
	require.NoError(t, err, "Dial should not error")

	_ = ep.Shutdown(context.Background())
	<-done
	_ = ep.Close()
}
