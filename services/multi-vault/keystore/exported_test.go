// +build unit

package keystore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, GlobalKeyStore(), "Controller should have been set")

	var ks keystore.KeyStore
	SetGlobalKeyStore(ks)
	assert.Nil(t, GlobalKeyStore(), "Global should be reset to nil")
}
