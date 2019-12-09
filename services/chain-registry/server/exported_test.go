package server

import (
	"context"
	"testing"

	"github.com/containous/traefik/v2/pkg/config/static"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	SetGlobalStaticConfig(&static.Configuration{})
	Init(context.Background())
	assert.NotNil(t, GlobalServer(), "Global server should have been set")
}
