// +build unit

package rest

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, server, "Server should have been set")

	var s *http.Server
	SetGlobalServer(s)
	assert.Nil(t, GlobalServer(), "Global should be reset to nil")
}
