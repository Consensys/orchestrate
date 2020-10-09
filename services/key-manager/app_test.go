// +build unit

package keymanager

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := NewConfig(viper.New())
	cfg.Store.Type = "hashicorp-vault"

	_, err := NewTxSigner(cfg)
	assert.NoError(t, err, "Creating App should not error")
}
