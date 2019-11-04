package memory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, GlobalStore(), "Global secretStore should have been set")

	var store *SecretStore
	SetGlobalStore(store)
	assert.Nil(t, GlobalStore(), "Global should be reset to nil")
}
