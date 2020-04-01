// +build unit

package secretstore

import (
	"context"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/secretstore/services"
)

func TestInit(t *testing.T) {
	viper.Set(secretStoreViperKey, "in-memory")

	Init(context.Background())
	assert.NotNil(t, secretStore, "Global secretStore should have been set")

	var s services.SecretStore
	SetGlobalSecretStore(s)
	assert.Nil(t, GlobalSecretStore(), "Global should be reset to nil")

}
