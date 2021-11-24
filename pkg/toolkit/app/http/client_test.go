// +build unit

package http

import (
	"context"
	"net/http"
	"testing"

	"github.com/consensys/orchestrate/pkg/toolkit/app/auth/utils"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/stretchr/testify/assert"
)

func TestDefaultClient(t *testing.T) {
	client := NewClient(NewDefaultConfig())
	req := &http.Request{}
	_, _ = client.Transport.RoundTrip(req)
	
	assert.Empty(t, utils.GetAuthorizationHeader(req))
	assert.Empty(t, utils.GetAPIKeyHeaderValue(req))
}


func TestApiKeyClient(t *testing.T) {
	cfg := NewDefaultConfig()
	cfg.XAPIKey = "ApiKey"

	client := NewClient(cfg)
	req := &http.Request{
		Header: http.Header{}, 
	}

	_, _ = client.Transport.RoundTrip(req)
	
	assert.Empty(t, utils.GetAuthorizationHeader(req))
	assert.Equal(t, cfg.XAPIKey, utils.GetAPIKeyHeaderValue(req))
}

func TestAuthTokenForwardClient(t *testing.T) {
	cfg := NewDefaultConfig()

	client := NewClient(cfg)
	req := &http.Request{
		Header: http.Header{}, 
	}
	
	authToken := "Bearer AuthToken"
	userInfo := &multitenancy.UserInfo{
		AuthMode: multitenancy.AuthMethodJWT,
		AuthValue: authToken,
	}
	_, _ = client.Transport.RoundTrip(req.WithContext(multitenancy.WithUserInfo(context.Background(), userInfo)))
	assert.Equal(t, authToken, utils.GetAuthorizationHeader(req))
}


func TestSkipAuthTokenForwardClient(t *testing.T) {
	cfg := NewDefaultConfig()
	cfg.AuthHeaderForward = false

	client := NewClient(cfg)
	req := &http.Request{
		Header: http.Header{}, 
	}
	
	authToken := "Bearer AuthToken"
	_, _ = client.Transport.RoundTrip(req.WithContext(utils.WithAuthorization(context.Background(), authToken)))
	assert.Empty(t, utils.GetAuthorizationHeader(req))
}
