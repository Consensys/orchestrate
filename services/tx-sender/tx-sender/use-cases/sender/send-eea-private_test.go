// +build unit

package sender

import (
	"context"
	"testing"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/sdk/client/mock"
	mock2 "github.com/ConsenSys/orchestrate/pkg/toolkit/ethclient/mock"
	txschedulertypes "github.com/ConsenSys/orchestrate/pkg/types/api"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/pkg/types/testutils"
	"github.com/ConsenSys/orchestrate/pkg/utils"
	mocks2 "github.com/ConsenSys/orchestrate/services/tx-sender/tx-sender/nonce/mocks"
	"github.com/ConsenSys/orchestrate/services/tx-sender/tx-sender/use-cases/mocks"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func TestSendEEAPrivate_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	signTx := mocks.NewMockSignEEATransactionUseCase(ctrl)
	ec := mock2.NewMockEEATransactionSender(ctrl)
	crafter :=  mocks.NewMockCraftTransactionUseCase(ctrl)
	jobClient := mock.NewMockJobClient(ctrl)
	nonceManager := mocks2.NewMockManager(ctrl)
	chainRegistryURL := "chainRegistryURL:8081"
	ctx := context.Background()

	usecase := NewSendEEAPrivateTxUseCase(signTx, crafter, ec, jobClient, chainRegistryURL, nonceManager)
	
	t.Run("should execute use case successfully", func(t *testing.T) {
		job := testutils.FakeJob()
		raw := "rawData"
		txHash := "0x0000000000000000000000000000000000000000000000000000000000000abc"

		crafter.EXPECT().Execute(gomock.Any(), job).Return(nil)
		signTx.EXPECT().Execute(gomock.Any(), job).Return(raw, txHash, nil)
		job.Transaction.Raw = raw
		job.Transaction.Hash = txHash
		
		proxyURL := utils.GetProxyURL(chainRegistryURL, job.ChainUUID)
		ec.EXPECT().PrivDistributeRawTransaction(gomock.Any(), proxyURL, job.Transaction.Raw).Return(ethcommon.HexToHash(txHash), nil)
		nonceManager.EXPECT().IncrementNonce(gomock.Any(), job).Return(nil)
		jobClient.EXPECT().UpdateJob(gomock.Any(), job.UUID, &txschedulertypes.UpdateJobRequest{
			Status:      entities.StatusStored,
			Transaction: job.Transaction,
		})
		
		err := usecase.Execute(ctx, job)
		assert.NoError(t, err)
		assert.Equal(t, job.Transaction.Hash, txHash)
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
		ec.EXPECT().PrivDistributeRawTransaction(gomock.Any(), proxyURL, job.Transaction.Raw).Return(ethcommon.HexToHash(""), expectedErr)
		nonceManager.EXPECT().CleanNonce(gomock.Any(), job, expectedErr).Return(nil)
		
		err := usecase.Execute(ctx, job)
		assert.Equal(t, err, expectedErr)
	})
}
