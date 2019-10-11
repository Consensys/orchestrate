package cooldown

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, ctrl, "Controller should have been set")

	var cfg *Config
	SetGlobalConfig(cfg)
	assert.Nil(t, GlobalConfig(), "Global should be reset to nil")

	var ctrl *Controller
	SetGlobalController(ctrl)
	assert.Nil(t, GlobalController(), "Global should be reset to nil")
}
