package grpcserver

// TODO: this test does not pass make race

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/soheilhy/cmux"
// 	"github.com/spf13/viper"
// 	"github.com/stretchr/testify/assert"
// 	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/http"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/examples/helloworld/helloworld"
// )

// func TestListenAndServe(t *testing.T) {
// 	// Init GRPC server
// 	AddEnhancers(
// 		func(srv *grpc.Server) *grpc.Server {
// 			helloworld.RegisterGreeterServer(srv, &DummyGreeterServer{})
// 			return srv
// 		},
// 	)
// 	Init(context.Background())

// 	// Init HTTP server
// 	http.Init(context.Background())

// 	// ListenAndServe
// 	viper.Set("http.hostname", "127.0.0.1:0")
// 	errors := make(chan error)
// 	go func() {
// 		if err := ListenAndServe(); err != nil && err != cmux.ErrListenerClosed {
// 			errors <- err
// 		}
// 		close(errors)
// 	}()

// 	// Stop
// 	time.Sleep(100 * time.Millisecond)
// 	err := GracefulStop(context.Background())
// 	assert.Nil(t, err, "Shutdown should not error")
// 	for err := range errors {
// 		t.Fatal(err)
// 	}
// }
