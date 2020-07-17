// +build unit

package scheduler

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	clientutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/client-utils"
	types2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client/mock"
	"math/big"
	"testing"
)

func TestCredit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTxSchedulerClient := mock.NewMockTransactionSchedulerClient(ctrl)
	ctx := context.Background()
	request := &types.Request{
		ScheduleUUID: "scheduleUUID",
		ParentTxID:   "parentJobUUID",
		ChainUUID:    "chainUUID",
		Beneficiary:  common.HexToAddress("0x1"),
		FaucetsCandidates: map[string]types.Faucet{
			"faucet0": {
				Amount:   big.NewInt(10),
				Creditor: common.HexToAddress("0x2"),
			},
		},
		ElectedFaucet: "faucet0",
	}

	faucet := NewFaucet(mockTxSchedulerClient)

	t.Run("should credit account successfully with authToken", func(t *testing.T) {
		authToken := "myAuthToken"
		rctx := utils.WithAuthorization(ctx, authToken)

		expectedContext := context.WithValue(rctx, clientutils.RequestHeaderKey, map[string]string{
			multitenancy.AuthorizationMetadata: authToken,
		})
		expectedCreateJobRequest := &types2.CreateJobRequest{
			ScheduleUUID: request.ScheduleUUID,
			ChainUUID:    request.ChainUUID,
			Type:         types2.EthereumTransaction,
			Labels: map[string]string{
				"parentJobUUID": request.ParentTxID,
			},
			Transaction: &types2.ETHTransaction{
				From:  request.FaucetsCandidates[request.ElectedFaucet].Creditor.String(),
				To:    request.Beneficiary.String(),
				Value: request.FaucetsCandidates[request.ElectedFaucet].Amount.String(),
			},
		}

		mockTxSchedulerClient.EXPECT().CreateJob(expectedContext, expectedCreateJobRequest).Return(&types2.JobResponse{
			UUID: "uuid",
		}, nil)
		mockTxSchedulerClient.EXPECT().StartJob(expectedContext, "uuid").Return(nil)

		amount, err := faucet.Credit(rctx, request)

		assert.NoError(t, err)
		assert.Equal(t, request.FaucetsCandidates[request.ElectedFaucet].Amount, amount)
	})

	t.Run("should credit account successfully with API key", func(t *testing.T) {
		apiKey := "myAPIKey"
		rctx := utils.WithAPIKey(ctx, apiKey)

		expectedContext := context.WithValue(rctx, clientutils.RequestHeaderKey, map[string]string{
			utils.APIKeyHeader: apiKey,
		})
		expectedCreateJobRequest := &types2.CreateJobRequest{
			ScheduleUUID: request.ScheduleUUID,
			ChainUUID:    request.ChainUUID,
			Type:         types2.EthereumTransaction,
			Labels: map[string]string{
				"parentJobUUID": request.ParentTxID,
			},
			Transaction: &types2.ETHTransaction{
				From:  request.FaucetsCandidates[request.ElectedFaucet].Creditor.String(),
				To:    request.Beneficiary.String(),
				Value: request.FaucetsCandidates[request.ElectedFaucet].Amount.String(),
			},
		}

		mockTxSchedulerClient.EXPECT().CreateJob(expectedContext, expectedCreateJobRequest).Return(&types2.JobResponse{
			UUID: "uuid",
		}, nil)
		mockTxSchedulerClient.EXPECT().StartJob(expectedContext, "uuid").Return(nil)

		amount, err := faucet.Credit(rctx, request)

		assert.NoError(t, err)
		assert.Equal(t, request.FaucetsCandidates[request.ElectedFaucet].Amount, amount)
	})

	t.Run("should credit account via transfer if chainName is specified", func(t *testing.T) {
		requestChainName := &types.Request{
			ScheduleUUID: "scheduleUUID",
			ParentTxID:   "parentJobUUID",
			ChainUUID:    "chainUUID",
			Beneficiary:  common.HexToAddress("0x1"),
			FaucetsCandidates: map[string]types.Faucet{
				"faucet0": {
					Amount:   big.NewInt(10),
					Creditor: common.HexToAddress("0x2"),
				},
			},
			ElectedFaucet: "faucet0",
			ChainName:     "chainName",
		}
		expectedTransferRequest := &types2.TransferRequest{
			BaseTransactionRequest: types2.BaseTransactionRequest{
				ChainName: requestChainName.ChainName,
				Labels: map[string]string{
					"parentJobUUID": requestChainName.ParentTxID,
				},
			},
			Params: types2.TransferParams{
				Value: requestChainName.FaucetsCandidates[request.ElectedFaucet].Amount.String(),
				From:  requestChainName.FaucetsCandidates[request.ElectedFaucet].Creditor.String(),
				To:    requestChainName.Beneficiary.String(),
			},
		}

		mockTxSchedulerClient.EXPECT().SendTransferTransaction(ctx, expectedTransferRequest).Return(&types2.TransactionResponse{
			UUID: "uuid",
		}, nil)

		amount, err := faucet.Credit(ctx, requestChainName)

		assert.NoError(t, err)
		assert.Equal(t, requestChainName.FaucetsCandidates[request.ElectedFaucet].Amount, amount)
	})

	t.Run("should fail with ServiceConnectionError if creating job fails", func(t *testing.T) {
		mockTxSchedulerClient.EXPECT().CreateJob(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))

		_, err := faucet.Credit(ctx, request)

		assert.True(t, errors.IsServiceConnectionError(err))
	})

	t.Run("should fail with ServiceConnectionError if starting job fails", func(t *testing.T) {
		mockTxSchedulerClient.EXPECT().CreateJob(gomock.Any(), gomock.Any()).Return(&types2.JobResponse{
			UUID: "uuid",
		}, nil)
		mockTxSchedulerClient.EXPECT().StartJob(gomock.Any(), "uuid").Return(fmt.Errorf("error"))

		_, err := faucet.Credit(ctx, request)

		assert.True(t, errors.IsServiceConnectionError(err))
	})
}
