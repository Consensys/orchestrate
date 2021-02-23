// +build unit

package contracts

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"
)

func TestRegisterContract_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockDB := mocks.NewMockDB(ctrl)
	mockDBTX := mocks.NewMockTx(ctrl)
	artifactAgent := mocks.NewMockArtifactAgent(ctrl)
	repositoryAgent := mocks.NewMockRepositoryAgent(ctrl)
	methodAgent := mocks.NewMockMethodAgent(ctrl)
	eventAgent := mocks.NewMockEventAgent(ctrl)
	tagAgent := mocks.NewMockTagAgent(ctrl)

	mockDB.EXPECT().Begin().Return(mockDBTX, nil).AnyTimes()
	mockDBTX.EXPECT().Artifact().Return(artifactAgent).AnyTimes()
	mockDBTX.EXPECT().Repository().Return(repositoryAgent).AnyTimes()
	mockDBTX.EXPECT().Method().Return(methodAgent).AnyTimes()
	mockDBTX.EXPECT().Event().Return(eventAgent).AnyTimes()
	mockDBTX.EXPECT().Tag().Return(tagAgent).AnyTimes()

	usecase := NewRegisterContractUseCase(mockDB)

	//@TODO Add more advance test flows
	t.Run("should execute use case successfully", func(t *testing.T) {
		contract := testutils.FakeContract()
		repositoryAgent.EXPECT().SelectOrInsert(gomock.Any(), gomock.AssignableToTypeOf(&models.RepositoryModel{})).Return(nil)
		artifactAgent.EXPECT().SelectOrInsert(gomock.Any(), gomock.AssignableToTypeOf(&models.ArtifactModel{})).Return(nil)
		tagAgent.EXPECT().Insert(gomock.Any(), gomock.AssignableToTypeOf(&models.TagModel{}))
		methodAgent.EXPECT().InsertMultiple(gomock.Any(), gomock.AssignableToTypeOf([]*models.MethodModel{}))
		eventAgent.EXPECT().InsertMultiple(gomock.Any(), gomock.AssignableToTypeOf([]*models.EventModel{}))
		mockDBTX.EXPECT().Commit().Return(nil)
		err := usecase.Execute(ctx, contract)

		assert.NoError(t, err)
	})
}
