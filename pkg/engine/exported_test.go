package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, e, "Engine should have been set")

	var e *Engine
	SetGlobalEngine(e)
	assert.Nil(t, GlobalEngine(), "Global should be reset to nil")
}
