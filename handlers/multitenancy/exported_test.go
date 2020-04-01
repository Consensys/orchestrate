// +build unit

package multitenancy

import (
	"context"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, GlobalHandler(), "Global should have been set")

	var h engine.HandlerFunc
	SetGlobalHandler(h)
	assert.Nil(t, GlobalHandler(), "Global should be reset to nil")
}
