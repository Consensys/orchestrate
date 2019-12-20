package api

import (
	"context"
	"testing"

	"github.com/containous/traefik/v2/pkg/config/runtime"

	"github.com/containous/traefik/v2/pkg/config/static"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init(context.Background())
}

func TestNewBuilder(t *testing.T) {

	config := &static.Configuration{
		API: &static.API{Dashboard: true},
	}
	builder := NewBuilder(config)
	assert.NotNil(t, builder, "New builder should initialize an handler")

	configuration := &runtime.Configuration{}
	handler := builder(configuration)
	assert.NotNil(t, handler, "builder should initialize http.Handler")
}
