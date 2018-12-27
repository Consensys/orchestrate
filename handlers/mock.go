package handlers

import (
	"math/rand"
	"sync"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
)

// MockHandler is an handler used mainly for test
type MockHandler struct {
	mux     *sync.Mutex
	handled []*infra.Context

	maxtime int
}

// NewMockHandler creates a new mock handler
func NewMockHandler(maxtime int) *MockHandler {
	return &MockHandler{
		&sync.Mutex{},
		[]*infra.Context{},
		maxtime,
	}
}

// Handler returns handler
func (h *MockHandler) Handler() infra.HandlerFunc {
	return func(ctx *infra.Context) {
		// We add some randomness in time execution
		r := rand.Intn(h.maxtime)
		time.Sleep(time.Duration(r) * time.Millisecond)

		// Update handled context
		h.mux.Lock()
		defer h.mux.Unlock()
		h.handled = append(h.handled, ctx)
	}
}
