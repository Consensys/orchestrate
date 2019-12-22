package kafka

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, GlobalHook(), "Global should have been set")

	SetGlobalHook(nil)
	assert.Nil(t, GlobalHook(), "Global should be reset to nil")
}
