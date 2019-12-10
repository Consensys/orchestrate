package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/pg"
)

func TestInit(t *testing.T) {
	Init()
	assert.NotNil(t, GlobalStoreRegistry(), "Global should have been set")

	var chainRegistry *pg.ChainRegistry
	SetGlobalStoreRegistry(chainRegistry)
	assert.Nil(t, GlobalStoreRegistry(), "Global should be reset to nil")
}
