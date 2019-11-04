package enricher

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, GlobalHandler(), "Global should have been set")

	var handler engine.HandlerFunc
	SetGlobalHandler(handler)
	assert.Nil(t, GlobalHandler(), "Global should be reset to nil")
}
