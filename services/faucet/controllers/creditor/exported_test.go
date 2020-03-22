// +build unit

package creditor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/faucet"
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
