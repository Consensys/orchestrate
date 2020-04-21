package usecases

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	genuuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	mockstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

func TestUpdateChain_ByUUID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	chainAgent := mockstore.NewMockChainAgent(mockCtrl)

	updateChainUC := NewUpdateChain(chainAgent)
	uuid := genuuid.NewV4().String()

	chain := &models.Chain{
		Name: "geth",
		URLs: []string{"http://geth:8545"},
	}

	chainAgent.EXPECT().UpdateChainByUUID(gomock.Any(), gomock.Eq(uuid), gomock.Eq(chain)).Times(1)

	err := updateChainUC.Execute(context.Background(), uuid, "", chain)
	assert.Nil(t, err)
}

func TestUpdateChain_ByName(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	chainAgent := mockstore.NewMockChainAgent(mockCtrl)

	updateChainUC := NewUpdateChain(chainAgent)
	chaninName := "Geth"

	chain := &models.Chain{
		URLs: []string{"http://geth:8545"},
	}

	chainAgent.EXPECT().UpdateChainByName(gomock.Any(), gomock.Eq(chaninName), gomock.Eq(chain)).Times(1)

	err := updateChainUC.Execute(context.Background(), "", chaninName, chain)
	assert.Nil(t, err)
}

func TestUpdateChain_NotAllowUUIDUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	chainAgent := mockstore.NewMockChainAgent(mockCtrl)

	updateChainUC := NewUpdateChain(chainAgent)
	uuid := genuuid.NewV4().String()

	chain := &models.Chain{
		UUID: genuuid.NewV4().String(),
		Name: "geth",
		URLs: []string{"http://geth:8545"},
	}

	err := updateChainUC.Execute(context.Background(), uuid, "", chain)
	assert.NotNil(t, err)
	assert.True(t, errors.IsConstraintViolatedError(err), "should be IsConstraintViolatedError")
}
