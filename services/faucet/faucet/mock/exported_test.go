package mock

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, fct, "Faucet should have been set")

	var f *Faucet
	SetGlobalFaucet(f)
	assert.Nil(t, GlobalFaucet(), "Global should be reset to nil")
}
