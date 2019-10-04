package metrics

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	var s *http.Server
	SetGlobalServer(s)
	assert.Nil(t, GlobalServer(), "Global should be reset to nil")

	Init(context.Background())
	assert.NotNil(t, server, "Server should have been set")
	assert.NotNil(t, server.Handler, "Server should have been set")

	go ListenAndServe()
	StopServer(context.Background())
}
