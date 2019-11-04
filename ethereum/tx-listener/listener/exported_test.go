package listener

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/tx-listener/listener/base"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, GlobalListener(), "Global should have been set")

	var l TxListener
	var cfg *base.Config
	SetGlobalListener(l)
	SetGlobalConfig(cfg)
	assert.Nil(t, GlobalListener(), "Global should be reset to nil")
}
