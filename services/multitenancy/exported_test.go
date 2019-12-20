package multitenancy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, keyBuilder, "auth Manager should have been set")

	var k *KeyBuilder
	SetKeyBuilder(k)
	assert.Nil(t, GlobalKeyBuilder(), "Global should be reset to nil")
}
