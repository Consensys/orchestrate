package secretstore

import (
	"context"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	viper.Set(secretStoreViperKey, "test")

	Init(context.Background())
	assert.NotNil(t, secretStore, "Global secretStore should have been set")
}
