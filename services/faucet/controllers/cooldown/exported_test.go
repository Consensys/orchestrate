// +build unit

package cooldown

import (
	"context"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/faucet"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, ctrl, "Controller should have been set")

	var ctrl *Controller
	SetGlobalController(ctrl)
	assert.Nil(t, GlobalController(), "Global should be reset to nil")
}

func TestControl(t *testing.T) {
	Init(context.Background())
	var f faucet.CreditFunc

	assert.NotNil(t, Control(f))
}
