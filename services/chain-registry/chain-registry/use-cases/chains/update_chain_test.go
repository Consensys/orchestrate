package chains

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
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
	chainUUID := uuid.Must(uuid.NewV4()).String()

	chain := &models.Chain{
		Name: "geth",
		URLs: []string{"http://geth:8545"},
	}

	chainAgent.EXPECT().UpdateChain(gomock.Any(), gomock.Eq(chainUUID), []string{}, gomock.Eq(chain)).Times(1)

	err := updateChainUC.Execute(context.Background(), chainUUID, "", []string{}, chain)
	assert.NoError(t, err)
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

	chainAgent.EXPECT().UpdateChainByName(gomock.Any(), gomock.Eq(chaninName), []string{}, gomock.Eq(chain)).Times(1)

	err := updateChainUC.Execute(context.Background(), "", chaninName, []string{}, chain)
	assert.NoError(t, err)
}

func TestUpdateChain_NotAllowUUIDUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	chainAgent := mockstore.NewMockChainAgent(mockCtrl)

	updateChainUC := NewUpdateChain(chainAgent)
	chainUUID := uuid.Must(uuid.NewV4()).String()

	chain := &models.Chain{
		UUID: uuid.Must(uuid.NewV4()).String(),
		Name: "geth",
		URLs: []string{"http://geth:8545"},
	}

	err := updateChainUC.Execute(context.Background(), chainUUID, "", []string{}, chain)
	assert.Error(t, err)
	assert.True(t, errors.IsConstraintViolatedError(err), "should be IsConstraintViolatedError")
}
