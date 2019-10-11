package creditor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, ctrl, "Controller should have been set")

	var ctrl *Controller
	SetGlobalController(ctrl)
	assert.Nil(t, GlobalController(), "Global should be reset to nil")
}
