package keystore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, GlobalKeyStore(), "Controller should have been set")

	var ks KeyStore
	SetGlobalKeyStore(ks)
	assert.Nil(t, GlobalKeyStore(), "Global should be reset to nil")
}
