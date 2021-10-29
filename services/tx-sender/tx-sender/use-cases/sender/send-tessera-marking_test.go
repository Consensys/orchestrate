// +build unit

package sender

import (
	"context"
	"testing"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/sdk/client/mock"
	mock2 "github.com/consensys/orchestrate/pkg/toolkit/ethclient/mock"
	"github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/types/testutils"
	"github.com/consensys/orchestrate/pkg/utils"
	mocks2 "github.com/consensys/orchestrate/services/tx-sender/tx-sender/nonce/mocks"
	"github.com/consensys/orchestrate/services/tx-sender/tx-sender/use-cases/mocks"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func TestSendTesseraMarking_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	signTx := mocks.NewMockSignETHTransactionUseCase(ctrl)
	ec := mock2.NewMockQuorumTransactionSender(ctrl)
	crafter := mocks.NewMockCraftTransactionUseCase(ctrl)
	jobClient := mock.NewMockJobClient(ctrl)
	nonceChecker := mocks2.NewMockManager(ctrl)
	chainRegistryURL := "chainRegistryURL:8081"
	ctx := context.Background()

	usecase := NewSendTesseraMarkingTxUseCase(signTx, crafter, ec, jobClient, chainRegistryURL, nonceChecker)

	t.Run("should execute use case successfully", func(t *testing.T) {
		job := testutils.FakeJob()
		job.Transaction.PrivateFor = []string{"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="}
		raw := "rawData"
		txHash := "0x0000000000000000000000000000000000000000000000000000000000000abc"

		crafter.EXPECT().Execute(gomock.Any(), job).Return(nil)
		signTx.EXPECT().Execute(gomock.Any(), job).Return(raw, txHash, nil)
		job.Transaction.Raw = raw
		job.Transaction.Hash = txHash

		proxyURL := utils.GetProxyURL(chainRegistryURL, job.ChainUUID)
		ec.EXPECT().SendQuorumRawPrivateTransaction(gomock.Any(), proxyURL, job.Transaction.Raw, job.Transaction.PrivateFor, nil, 0).
			Return(ethcommon.HexToHash(txHash), nil)
		nonceChecker.EXPECT().IncrementNonce(gomock.Any(), job).Return(nil)
		jobClient.EXPECT().UpdateJob(gomock.Any(), job.UUID, &api.UpdateJobRequest{
			Status:      entities.StatusPending,
			Transaction: job.Transaction,
		})

		err := usecase.Execute(ctx, job)
		assert.NoError(t, err)
		assert.Equal(t, job.Transaction.Hash, txHash)
	})
	
	t.Run("should execute use case, using resending, successfully", func(t *testing.T) {
		job := testutils.FakeJob()
		job.Status = entities.StatusPending
		job.Transaction.PrivateFor = []string{"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="}
		raw := "rawData"
		txHash := "0x0000000000000000000000000000000000000000000000000000000000000abc"

		crafter.EXPECT().Execute(gomock.Any(), job).Return(nil)
		job.Transaction.Raw = raw
		job.Transaction.Hash = txHash

		proxyURL := utils.GetProxyURL(chainRegistryURL, job.ChainUUID)
		ec.EXPECT().SendQuorumRawPrivateTransaction(gomock.Any(), proxyURL, job.Transaction.Raw, job.Transaction.PrivateFor, nil, 0).
			Return(ethcommon.HexToHash(txHash), nil)
		nonceChecker.EXPECT().IncrementNonce(gomock.Any(), job).Return(nil)
		jobClient.EXPECT().UpdateJob(gomock.Any(), job.UUID, &api.UpdateJobRequest{
			Status:      entities.StatusResending,
			Transaction: job.Transaction,
		})

		err := usecase.Execute(ctx, job)
		assert.NoError(t, err)
		assert.Equal(t, job.Transaction.Hash, txHash)
	})

	t.Run("should execute use case, update warning, successfully", func(t *testing.T) {
		job := testutils.FakeJob()
		raw := "rawData"
		txHash := "0x0000000000000000000000000000000000000000000000000000000000000abc"
		txHash2 := "0x0000000000000000000000000000000000000000000000000000000000000aba"

		crafter.EXPECT().Execute(gomock.Any(), job).Return(nil)
		signTx.EXPECT().Execute(gomock.Any(), job).Return(raw, txHash, nil)
		job.Transaction.Raw = raw
		job.Transaction.Hash = txHash

		proxyURL := utils.GetProxyURL(chainRegistryURL, job.ChainUUID)
		ec.EXPECT().SendQuorumRawPrivateTransaction(gomock.Any(), proxyURL, job.Transaction.Raw, job.Transaction.PrivateFor, nil, 0).
			Return(ethcommon.HexToHash(txHash2), nil)
		nonceChecker.EXPECT().IncrementNonce(gomock.Any(), job).Return(nil)
		jobClient.EXPECT().UpdateJob(gomock.Any(), job.UUID, &api.UpdateJobRequest{
			Status:      entities.StatusPending,
			Transaction: job.Transaction,
		})

		jobClient.EXPECT().UpdateJob(gomock.Any(), job.UUID, gomock.Any())

		err := usecase.Execute(ctx, job)
		assert.NoError(t, err)
		assert.Equal(t, job.Transaction.Hash, txHash2)
	})

	t.Run("should fail to execute use case if nonce checker fails", func(t *testing.T) {
		job := testutils.FakeJob()

		expectedErr := errors.NonceTooLowWarning("invalid nonce")
		crafter.EXPECT().Execute(gomock.Any(), job).Return(expectedErr)

		err := usecase.Execute(ctx, job)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to execute use case if signer fails", func(t *testing.T) {
		job := testutils.FakeJob()
		raw := "rawData"
		txHash := "0x0000000000000000000000000000000000000000000000000000000000000abc"

		crafter.EXPECT().Execute(gomock.Any(), job).Return(nil)

		expectedErr := errors.InternalError("internal error")
		signTx.EXPECT().Execute(gomock.Any(), job).Return(raw, txHash, expectedErr)

		err := usecase.Execute(ctx, job)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to execute use case if update job fails", func(t *testing.T) {
		job := testutils.FakeJob()
		raw := "rawData"
		txHash := "0x0000000000000000000000000000000000000000000000000000000000000abc"

		crafter.EXPECT().Execute(gomock.Any(), job).Return(nil)
		signTx.EXPECT().Execute(gomock.Any(), job).Return(raw, txHash, nil)
		job.Transaction.Raw = raw
		job.Transaction.Hash = txHash

		expectedErr := errors.InternalError("internal error")
		jobClient.EXPECT().UpdateJob(gomock.Any(), job.UUID, &api.UpdateJobRequest{
			Status:      entities.StatusPending,
			Transaction: job.Transaction,
		}).Return(nil, expectedErr)

		err := usecase.Execute(ctx, job)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to execute use case if send transaction fails", func(t *testing.T) {
		job := testutils.FakeJob()
		raw := "rawData"
		txHash := "0x0000000000000000000000000000000000000000000000000000000000000abc"

		crafter.EXPECT().Execute(gomock.Any(), job).Return(nil)
		signTx.EXPECT().Execute(gomock.Any(), job).Return(raw, txHash, nil)
		job.Transaction.Raw = raw
		job.Transaction.Hash = txHash

		expectedErr := errors.InternalError("internal error")
		proxyURL := utils.GetProxyURL(chainRegistryURL, job.ChainUUID)
		jobClient.EXPECT().UpdateJob(gomock.Any(), job.UUID, &api.UpdateJobRequest{
			Status:      entities.StatusPending,
			Transaction: job.Transaction,
		})
		ec.EXPECT().SendQuorumRawPrivateTransaction(gomock.Any(), proxyURL, job.Transaction.Raw, job.Transaction.PrivateFor, nil, 0).
			Return(ethcommon.HexToHash(txHash), expectedErr)
		nonceChecker.EXPECT().CleanNonce(gomock.Any(), job, expectedErr).Return(nil)

		err := usecase.Execute(ctx, job)
		assert.Equal(t, err, expectedErr)
	})
}
