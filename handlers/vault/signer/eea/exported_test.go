// +build unit

package eea

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, GlobalHandler(), "Global should have been set")

	var h engine.HandlerFunc
	SetGlobalHandler(h)
	assert.Nil(t, GlobalHandler(), "Global should be reset to nil")
}
