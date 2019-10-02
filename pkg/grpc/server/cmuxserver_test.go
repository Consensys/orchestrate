package grpcserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"testing"

	"github.com/soheilhy/cmux"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/examples/helloworld/helloworld"
)

type MockGreeterServer struct {
	mux      *sync.Mutex
	requests []*helloworld.HelloRequest
}

func NewMockGreeterServer() *MockGreeterServer {
	return &MockGreeterServer{
		mux:      &sync.Mutex{},
		requests: []*helloworld.HelloRequest{},
	}
}

// SayHello implements helloworld.GreeterServer
func (s *MockGreeterServer) SayHello(ctx context.Context, req *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.requests = append(s.requests, req)
	return &helloworld.HelloReply{Message: "Hello " + req.Name}, nil
}

type MockHandler struct {
	mux      *sync.Mutex
	requests []*http.Request
}

func NewMockHandler() *MockHandler {
	return &MockHandler{
		mux:      &sync.Mutex{},
		requests: []*http.Request{},
	}
}

func (h *MockHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.mux.Lock()
	defer h.mux.Unlock()
	h.requests = append(h.requests, req)
	w.WriteHeader(http.StatusOK)
}

func TestCMuxServerListen(t *testing.T) {
	mux := http.NewServeMux()
	mockh := NewMockHandler()
	mux.Handle("/test", mockh)

	cmuxsrv := NewCMuxServer(
		grpc.NewServer(),
		&http.Server{
			Handler: mux,
		},
	)

	mocksrv := NewMockGreeterServer()
	helloworld.RegisterGreeterServer(cmuxsrv.grpc, mocksrv)

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	errors := make(chan error)
	go func() {
		if e := cmuxsrv.Serve(lis); e != nil && e != cmux.ErrListenerClosed {
			errors <- e
		}
		close(errors)
	}()

	// Test GPRC
	grpcconn, err := grpc.DialContext(
		context.Background(),
		lis.Addr().String(),
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Errorf("Failed to dial GPRC: %v", err)
	}

	client := helloworld.NewGreeterClient(grpcconn)
	_, err = client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "test"})
	assert.Nil(t, err, "SayHello should not error")
	assert.Len(t, mocksrv.requests, 1, "GRPC should have treated a request")
	_ = grpcconn.Close()

	// HTTP
	_, err = http.Get(fmt.Sprintf("http://%v/test", lis.Addr().String()))
	if err != nil {
		t.Errorf("Failed to dial HTTP: %v", err)
	}
	assert.Nil(t, err, "Get should not error")
	assert.Len(t, mockh.requests, 1, "HTTP should have treated a request")
}
