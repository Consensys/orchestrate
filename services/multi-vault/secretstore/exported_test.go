package secretstore

import (
	"context"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/multi-vault/secretstore/services"
)

func TestInit(t *testing.T) {
	viper.Set(secretStoreViperKey, "test")

	Init(context.Background())
	assert.NotNil(t, secretStore, "Global secretStore should have been set")

	var secretStore services.SecretStore
	SetGlobalSecretStore(secretStore)
	assert.Nil(t, GlobalSecretStore(), "Global should be reset to nil")

}