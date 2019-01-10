package handlers

import (
	"math/rand"
	"sync"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/core/types"
)

// MockHandler is an handler used mainly for test
type MockHandler struct {
	mux     *sync.Mutex
	handled []*types.Context

	maxtime int
}

// NewMockHandler creates a new mock handler
func NewMockHandler(maxtime int) *MockHandler {
	return &MockHandler{
		&sync.Mutex{},
		[]*types.Context{},
		maxtime,
	}
}

// Handler returns handler
func (h *MockHandler) Handler() types.HandlerFunc {
	return func(ctx *types.Context) {
		// We add some randomness in time execution
		r := rand.Intn(h.maxtime)
		time.Sleep(time.Duration(r) * time.Millisecond)

		// Update handled context
		h.mux.Lock()
		defer h.mux.Unlock()
		h.handled = append(h.handled, ctx)
	}
}
