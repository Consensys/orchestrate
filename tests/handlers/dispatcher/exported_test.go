package dispatcher

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
)

func testKeyOf(txctx *engine.TxContext) (string, error) {
	return "", nil
}

func TestInit(t *testing.T) {
	SetKeyOfFuncs(testKeyOf)

	Init(context.Background())
	assert.NotNil(t, GlobalHandler(), "Global should have been set")

	var h engine.HandlerFunc
	SetGlobalHandler(h)
	assert.Nil(t, GlobalHandler(), "Global should be reset to nil")
}
