// +build unit

package usecases

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/models"
)

func TestSetCodeHash_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	account := testutils.FakeAccount()
	codeHash := "codeHash"
	codeHashModel := &models.CodehashModel{
		ChainID:  account.ChainId,
		Address:  account.Account,
		Codehash: codeHash,
	}
	mockCodeHashDataAgent := mock.NewMockCodeHashDataAgent(ctrl)
	usecase := NewSetCodeHash(mockCodeHashDataAgent)

	t.Run("should execute use case successfully", func(t *testing.T) {
		mockCodeHashDataAgent.EXPECT().Insert(context.Background(), codeHashModel).Return(nil)

		err := usecase.Execute(context.Background(), account, codeHash)

		assert.NoError(t, err)
	})

	t.Run("should fail if data agent fails", func(t *testing.T) {
		dataAgentError := fmt.Errorf("error")
		mockCodeHashDataAgent.EXPECT().Insert(context.Background(), codeHashModel).Return(dataAgentError)

		err := usecase.Execute(context.Background(), account, codeHash)

		assert.Equal(t, errors.FromError(dataAgentError).ExtendComponent(setCodeHashComponent), err)
	})
}
