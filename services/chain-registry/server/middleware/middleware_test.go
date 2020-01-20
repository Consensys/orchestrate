package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/containous/alice"
	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/config/runtime"
	"github.com/stretchr/testify/assert"
)

type MockHandler struct {
	served int
	next   http.Handler
}

func (h *MockHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	h.served++
	if h.next != nil {
		h.next.ServeHTTP(rw, req)
	}
}

func TestBuilder(t *testing.T) {
	h1 := &MockHandler{}
	customMiddlewares := map[string]alice.Constructor{
		"custom-middleware": func(next http.Handler) (http.Handler, error) {
			h1.next = next
			return h1, nil
		},
	}

	rtConf := runtime.NewConfig(dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Middlewares: map[string]*dynamic.Middleware{
				"traefik-middleware": {
					Headers: &dynamic.Headers{
						CustomRequestHeaders: map[string]string{"traefik-middleware": "value-traefik-middleware"},
					},
				},
			},
		},
	})

	// Build chain and give it and handler
	builder := NewBuilder(rtConf.Middlewares, nil, customMiddlewares)
	chain := builder.BuildChain(
		context.Background(),
		[]string{"custom-middleware", "traefik-middleware"},
	)
	h2 := &MockHandler{}
	h, _ := chain.Then(h2)

	// Test serve
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	// Check that custom middleware has been served
	assert.Equal(t, 1, h1.served, "Custom Middleware 1 should have been served")

	// Check that trefik middleware has been served
	assert.Equal(t, "value-traefik-middleware", req.Header.Get("traefik-middleware"), "Traefik Middleware should have been served")

	// Check that given handler has been served
	assert.Equal(t, 1, h2.served, "Given Handler should have been served")
}
