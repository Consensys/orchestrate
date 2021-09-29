// +build unit

package chainregistry

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/consensys/orchestrate/pkg/sdk/client/mock"
)

func TestInit(t *testing.T) {
	ctrl := gomock.NewController(t)
	Init(mock.NewMockChainClient(ctrl))
	assert.NotNil(t, GlobalManager(), "Global should have been set")

	var mngr *Manager
	SetGlobalManager(mngr)
	assert.Nil(t, GlobalManager(), "Global should be reset to nil")
}
