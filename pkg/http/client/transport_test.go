package httpclient

import (
	"context"
	"net/http"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication"

	"github.com/stretchr/testify/assert"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/utils"
)

type MockTransport struct {
	roundTrips int
}

func (t *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.roundTrips++
	return nil, nil
}

func TestTransport(t *testing.T) {
	mockTransport := &MockTransport{}
	transport := NewTransport(mockTransport)

	// Test without setting authorization in context
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	_, _ = transport.RoundTrip(req)
	assert.Equal(t, 1, mockTransport.roundTrips, "Mock transport should have been called")
	auth := req.Header.Get("Authorization")
	assert.Equal(t, "", auth, "Authorization header shuld be empty")

	// Test setting authorization in context
	req, _ = http.NewRequestWithContext(
		authutils.WithAuthorization(context.Background(), "test-auth"),
		http.MethodGet, "", nil,
	)
	_, _ = transport.RoundTrip(req)
	assert.Equal(t, 2, mockTransport.roundTrips, "Mock transport should have been called")
	auth = req.Header.Get("Authorization")
	assert.Equal(t, "test-auth", auth, "Authorization header shuld be empty")

	// Test setting X-API-Key in context
	req, _ = http.NewRequestWithContext(
		authutils.WithAPIKey(context.Background(), "test-auth"),
		http.MethodGet, "", nil,
	)
	_, _ = transport.RoundTrip(req)
	assert.Equal(t, 3, mockTransport.roundTrips, "Mock transport should have been called")
	auth = req.Header.Get(authentication.APIKeyHeader)
	assert.Equal(t, "test-auth", auth, "Authorization header shuld be empty")
}
