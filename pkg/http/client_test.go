// +build unit

package http

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ConsenSys/orchestrate/pkg/auth/utils"
)

func TestDefaultClient(t *testing.T) {
	client := NewClient(NewDefaultConfig())
	req := &http.Request{}
	_, _ = client.Transport.RoundTrip(req)
	
	assert.Empty(t, req.Header.Get(utils.AuthorizationHeader))
	assert.Empty(t, req.Header.Get(utils.APIKeyHeader))
}


func TestApiKeyClient(t *testing.T) {
	cfg := NewDefaultConfig()
	cfg.MultiTenancy = true
	cfg.APIKey = "ApiKey"

	client := NewClient(cfg)
	req := &http.Request{
		Header: http.Header{}, 
	}

	_, _ = client.Transport.RoundTrip(req)
	
	assert.Empty(t, req.Header.Get(utils.AuthorizationHeader))
	assert.Equal(t, cfg.APIKey, req.Header.Get(utils.APIKeyHeader))
}

func TestAuthTokenForwardClient(t *testing.T) {
	cfg := NewDefaultConfig()
	cfg.MultiTenancy = true

	client := NewClient(cfg)
	req := &http.Request{
		Header: http.Header{}, 
	}
	
	authToken := "Bearer AuthToken"
	_, _ = client.Transport.RoundTrip(req.WithContext(utils.WithAuthorization(context.Background(), authToken)))
	assert.Equal(t, authToken, req.Header.Get(utils.AuthorizationHeader))
}


func TestMultitenancyClient(t *testing.T) {
	cfg := NewDefaultConfig()
	cfg.MultiTenancy = true
	cfg.APIKey = "ApiKey"

	client := NewClient(cfg)
	req := &http.Request{
		Header: http.Header{}, 
	}
	
	authToken := "Bearer AuthToken"
	_, _ = client.Transport.RoundTrip(req.WithContext(utils.WithAuthorization(context.Background(), authToken)))
	assert.Equal(t, authToken, req.Header.Get(utils.AuthorizationHeader))
	assert.Empty(t, req.Header.Get(utils.APIKeyHeader))
}

func TestSkipAuthTokenForwardClient(t *testing.T) {
	cfg := NewDefaultConfig()
	cfg.MultiTenancy = true
	cfg.AuthHeaderForward = false

	client := NewClient(cfg)
	req := &http.Request{
		Header: http.Header{}, 
	}
	
	authToken := "Bearer AuthToken"
	_, _ = client.Transport.RoundTrip(req.WithContext(utils.WithAuthorization(context.Background(), authToken)))
	assert.Empty(t, req.Header.Get(utils.AuthorizationHeader))
}
