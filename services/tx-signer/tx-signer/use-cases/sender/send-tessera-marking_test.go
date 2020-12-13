// +build unit

package sender

import (
	"context"
	"fmt"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	mock2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	txschedulertypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/txscheduler"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	mock3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/client/mock"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-signer/tx-signer/nonce/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-signer/tx-signer/use-cases/mocks"

	"github.com/golang/mock/gomock"
)

func TestSendTesseraMarking_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	signTx := mocks.NewMockSignETHTransactionUseCase(ctrl)
	ec := mock2.NewMockQuorumTransactionSender(ctrl)
	txSchedulerClient := mock3.NewMockTransactionSchedulerClient(ctrl)
	chainRegistryURL := "chainRegistryURL:8081"
	nonceChecker := mocks2.NewMockChecker(ctrl)
	ctx := context.Background()

	usecase := NewSendTesseraMarkingTxUseCase(signTx, ec, txSchedulerClient, chainRegistryURL, nonceChecker)
	
	t.Run("should execute use case successfully", func(t *testing.T) {
		job := testutils.FakeJob()
		job.Transaction.PrivateFor = []string{"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="} 
		raw := "rawData"
		txHash := "0x0000000000000000000000000000000000000000000000000000000000000abc"

		nonceChecker.EXPECT().Check(ctx, job).Return(nil)
		signTx.EXPECT().Execute(ctx, job).Return(raw, txHash, nil)
		job.Transaction.Raw = raw
		job.Transaction.Hash = txHash
		
		proxyURL := fmt.Sprintf("%s/%s", chainRegistryURL, job.ChainUUID)
		ec.EXPECT().SendQuorumRawPrivateTransaction(ctx, proxyURL, job.Transaction.Raw, job.Transaction.PrivateFor).
			Return(ethcommon.HexToHash(txHash), nil)
		nonceChecker.EXPECT().OnSuccess(ctx, job).Return(nil)
		txSchedulerClient.EXPECT().UpdateJob(ctx, job.UUID, &txschedulertypes.UpdateJobRequest{
			Status:      utils.StatusPending,
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
	
		nonceChecker.EXPECT().Check(ctx, job).Return(nil)
		signTx.EXPECT().Execute(ctx, job).Return(raw, txHash, nil)
		job.Transaction.Raw = raw
		job.Transaction.Hash = txHash
		
		proxyURL := fmt.Sprintf("%s/%s", chainRegistryURL, job.ChainUUID)
		ec.EXPECT().SendQuorumRawPrivateTransaction(ctx, proxyURL, job.Transaction.Raw, job.Transaction.PrivateFor).
			Return(ethcommon.HexToHash(txHash2), nil)
		nonceChecker.EXPECT().OnSuccess(ctx, job).Return(nil)
		txSchedulerClient.EXPECT().UpdateJob(ctx, job.UUID, &txschedulertypes.UpdateJobRequest{
			Status:      utils.StatusPending,
			Transaction: job.Transaction,
		})
		
		txSchedulerClient.EXPECT().UpdateJob(ctx, job.UUID, gomock.Any())
		
		err := usecase.Execute(ctx, job)
		assert.NoError(t, err)
		assert.Equal(t, job.Transaction.Hash, txHash2)
	})
	
	t.Run("should fail to execute use case if nonce checker fails", func(t *testing.T) {
		job := testutils.FakeJob()
	
		expectedErr := errors.NonceTooLowWarning("invalid nonce")
		nonceChecker.EXPECT().Check(ctx, job).Return(expectedErr)
		
		err := usecase.Execute(ctx, job)
		assert.Equal(t, err, expectedErr)
	})
	
	t.Run("should fail to execute use case if signer fails", func(t *testing.T) {
		job := testutils.FakeJob()
		raw := "rawData"
		txHash := "0x0000000000000000000000000000000000000000000000000000000000000abc"
	
		nonceChecker.EXPECT().Check(ctx, job).Return(nil)
		
		expectedErr := errors.InternalError("internal error")
		signTx.EXPECT().Execute(ctx, job).Return(raw, txHash, expectedErr)
		
		err := usecase.Execute(ctx, job)
		assert.Equal(t, err, expectedErr)
	})
	
	t.Run("should fail to execute use case if update job fails", func(t *testing.T) {
		job := testutils.FakeJob()
		raw := "rawData"
		txHash := "0x0000000000000000000000000000000000000000000000000000000000000abc"
	
		nonceChecker.EXPECT().Check(ctx, job).Return(nil)
		signTx.EXPECT().Execute(ctx, job).Return(raw, txHash, nil)
		job.Transaction.Raw = raw
		job.Transaction.Hash = txHash
	
		expectedErr := errors.InternalError("internal error")
		txSchedulerClient.EXPECT().UpdateJob(ctx, job.UUID, &txschedulertypes.UpdateJobRequest{
			Status:      utils.StatusPending,
			Transaction: job.Transaction,
		}).Return(nil, expectedErr)
		
		err := usecase.Execute(ctx, job)
		assert.Equal(t, err, expectedErr)
	})
	
	t.Run("should fail to execute use case if send transaction fails", func(t *testing.T) {
		job := testutils.FakeJob()
		raw := "rawData"
		txHash := "0x0000000000000000000000000000000000000000000000000000000000000abc"
	
		nonceChecker.EXPECT().Check(ctx, job).Return(nil)
		signTx.EXPECT().Execute(ctx, job).Return(raw, txHash, nil)
		job.Transaction.Raw = raw
		job.Transaction.Hash = txHash
	
		expectedErr := errors.InternalError("internal error")
		proxyURL := fmt.Sprintf("%s/%s", chainRegistryURL, job.ChainUUID)
		txSchedulerClient.EXPECT().UpdateJob(ctx, job.UUID, &txschedulertypes.UpdateJobRequest{
			Status:      utils.StatusPending,
			Transaction: job.Transaction,
		})
		ec.EXPECT().SendQuorumRawPrivateTransaction(ctx, proxyURL, job.Transaction.Raw, job.Transaction.PrivateFor).
			Return(ethcommon.HexToHash(txHash), expectedErr)
		nonceChecker.EXPECT().OnFailure(ctx, job, expectedErr).Return(nil)
		
		err := usecase.Execute(ctx, job)
		assert.Equal(t, err, expectedErr)
	})
}
