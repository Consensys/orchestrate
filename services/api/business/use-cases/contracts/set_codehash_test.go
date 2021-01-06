// +build unit

package contracts

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/mocks"
	models2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"
)

func TestSetCodeHash_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	codeHash := "codeHash"
	codeHashModel := &models2.CodehashModel{
		ChainID:  chainID,
		Address:  contractAddress,
		Codehash: codeHash,
	}
	codeHashAgent := mocks.NewMockCodeHashAgent(ctrl)
	usecase := NewSetCodeHashUseCase(codeHashAgent)

	t.Run("should execute use case successfully", func(t *testing.T) {
		codeHashAgent.EXPECT().Insert(ctx, codeHashModel).Return(nil)

		err := usecase.Execute(ctx, chainID, contractAddress, codeHash)

		assert.NoError(t, err)
	})

	t.Run("should fail if data agent fails", func(t *testing.T) {
		dataAgentError := fmt.Errorf("error")
		codeHashAgent.EXPECT().Insert(ctx, codeHashModel).Return(dataAgentError)

		err := usecase.Execute(ctx, chainID, contractAddress, codeHash)

		assert.Equal(t, errors.FromError(dataAgentError).ExtendComponent(setCodeHashComponent), err)
	})
}
