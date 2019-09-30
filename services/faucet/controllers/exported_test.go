package controllers

import (
	"context"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	viper.Set("faucet.type", "mock")
	Init(context.Background())
	assert.NotNil(t, ctrl, "Control should have been set")
}
