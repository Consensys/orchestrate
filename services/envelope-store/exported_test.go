package envelopestore

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	evlpstore "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/services/envelope-store"
)

func TestInit(t *testing.T) {
	viper.Set(typeViperKey, "mock")
	Init()
	assert.NotNil(t, GlobalEnvelopeStoreServer(), "Global should have been set")

	var s evlpstore.EnvelopeStoreServer
	SetGlobalEnvelopeStoreServer(s)
	assert.Nil(t, GlobalEnvelopeStoreServer(), "Global should be reset to nil")
}