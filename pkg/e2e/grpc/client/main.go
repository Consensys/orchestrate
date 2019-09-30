package main

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	grpcclient "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/grpc/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/tracing/opentracing/jaeger"
	"google.golang.org/grpc/examples/helloworld/helloworld"
)

func sayHello(client helloworld.GreeterClient, name string, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := client.SayHello(
		ctx,
		&helloworld.HelloRequest{Name: name},
	)
	if err != nil {
		log.WithError(err).Errorf("e2e: SayHello to %q failed", name)
	} else {
		log.Infof("e2e: SayHello to %q succeeded (%q)", name, resp.GetMessage())
	}
}

func main() {
	// Set log level to debug
	log.SetLevel(log.DebugLevel)

	// Initialize Jaegger
	jaeger.Init(context.Background())

	conn, err := grpcclient.DialContextWithDefaultOptions(
		context.Background(),
		"localhost:8080",
	)
	if err != nil {
		log.WithError(err).Fatalf("e2e: GRPC client could not connect to server")
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.WithError(err).Warn("could not close gRPC connection")
		}
	}()

	client := helloworld.NewGreeterClient(conn)
	// Should succeed
	for i := 0; i < 10; i++ {
		sayHello(client, "test-name", 100*time.Millisecond)
	}

	// Should error
	sayHello(client, "", 100*time.Millisecond)

	// Should error
	sayHello(client, "test-name", 20*time.Millisecond)

	// Should error
	sayHello(client, "test-panic", 100*time.Millisecond)
}
