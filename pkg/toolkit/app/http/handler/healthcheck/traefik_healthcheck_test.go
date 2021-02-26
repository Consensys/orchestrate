// +build unit

package healthcheck

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTraefikHealthCheck(t *testing.T) {
	b := NewTraefikBuilder()
	ctx, cancel := context.WithCancel(context.Background())
	h, err := b.Build(ctx, "", &dynamic.HealthCheck{}, nil)
	require.NoError(t, err)

	// Test Live
	rw := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "http://host.com/live", nil)

	h.ServeHTTP(rw, req)
	resp := rw.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "/live should valid status")

	// Test Ready
	rw = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "http://host.com/ready", nil)

	h.ServeHTTP(rw, req)
	resp = rw.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "/ready should valid status")

	// Test Ready after canceling context
	cancel()
	time.Sleep(50 * time.Millisecond)
	rw = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "http://host.com/ready", nil)

	h.ServeHTTP(rw, req)
	resp = rw.Result()
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode, "/ready should valid status")
}
