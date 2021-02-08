// +build unit

package chainregistry

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client/mock"
)

func TestInit(t *testing.T) {
	ctrl := gomock.NewController(t)
	Init(mock.NewMockChainClient(ctrl))
	assert.NotNil(t, GlobalManager(), "Global should have been set")

	var mngr *Manager
	SetGlobalManager(mngr)
	assert.Nil(t, GlobalManager(), "Global should be reset to nil")
}
