// +build unit

package sender

import (
	"context"
	"testing"

	mock3 "github.com/consensys/orchestrate/pkg/sdk/client/mock"
	mock2 "github.com/consensys/orchestrate/pkg/toolkit/ethclient/mock"
	txschedulertypes "github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/types/testutils"
	"github.com/consensys/orchestrate/pkg/utils"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSendETHRaw_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ec := mock2.NewMockTransactionSender(ctrl)
	jobClient := mock3.NewMockJobClient(ctrl)
	chainRegistryURL := "chainRegistryURL:8081"
	ctx := context.Background()

	usecase := NewSendETHRawTxUseCase(ec, jobClient, chainRegistryURL)

	t.Run("should execute use case successfully", func(t *testing.T) {
		job := testutils.FakeJob()
		raw := "0xf85380839896808252088083989680808216b4a0d35c752d3498e6f5ca1630d264802a992a141ca4b6a3f439d673c75e944e5fb0a05278aaa5fabbeac362c321b54e298dedae2d31471e432c26ea36a8d49cf08f1e"
		txHash := "0x6621fbe1e2848446e38d99bfda159cdd83f555ae0ed7a4f3e1c3c79f7d6d74f3"

		job.Transaction.Raw = raw
		job.Transaction.Hash = txHash

		proxyURL := utils.GetProxyURL(chainRegistryURL, job.ChainUUID)
		ec.EXPECT().SendRawTransaction(gomock.Any(), proxyURL, raw).Return(ethcommon.HexToHash(txHash), nil)

		jobClient.EXPECT().UpdateJob(gomock.Any(), job.UUID, &txschedulertypes.UpdateJobRequest{
			Status:      entities.StatusPending,
		})

		err := usecase.Execute(ctx, job)
		assert.NoError(t, err)
		assert.Equal(t, job.Transaction.Hash, txHash)
	})
	
	t.Run("should execute use case, using resending, successfully", func(t *testing.T) {
		job := testutils.FakeJob()
		raw := "0xf85380839896808252088083989680808216b4a0d35c752d3498e6f5ca1630d264802a992a141ca4b6a3f439d673c75e944e5fb0a05278aaa5fabbeac362c321b54e298dedae2d31471e432c26ea36a8d49cf08f1e"
		txHash := "0x6621fbe1e2848446e38d99bfda159cdd83f555ae0ed7a4f3e1c3c79f7d6d74f3"

		job.Transaction.Raw = raw
		job.Transaction.Hash = txHash
		job.Status = entities.StatusPending

		proxyURL := utils.GetProxyURL(chainRegistryURL, job.ChainUUID)
		ec.EXPECT().SendRawTransaction(gomock.Any(), proxyURL, raw).Return(ethcommon.HexToHash(txHash), nil)

		jobClient.EXPECT().UpdateJob(gomock.Any(), job.UUID, &txschedulertypes.UpdateJobRequest{
			Status:      entities.StatusResending,
		})

		err := usecase.Execute(ctx, job)
		assert.NoError(t, err)
		assert.Equal(t, job.Transaction.Hash, txHash)
	})

	t.Run("should execute use case and update to warning successfully", func(t *testing.T) {
		job := testutils.FakeJob()
		raw := "0xf85380839896808252088083989680808216b4a0d35c752d3498e6f5ca1630d264802a992a141ca4b6a3f439d673c75e944e5fb0a05278aaa5fabbeac362c321b54e298dedae2d31471e432c26ea36a8d49cf08f1e"
		txHash := "0x6621fbe1e2848446e38d99bfda159cdd83f555ae0ed7a4f3e1c3c79f7d6d74f3"

		job.Transaction.Raw = raw
		job.Transaction.Hash = txHash

		hash := "0x6621fbe1e2848446e38d99bfda159cdd83f555ae0ed7a4f3e1c3c79f7d6d74f2"
		proxyURL := utils.GetProxyURL(chainRegistryURL, job.ChainUUID)
		ec.EXPECT().SendRawTransaction(gomock.Any(), proxyURL, raw).
			Return(ethcommon.HexToHash(hash), nil)

		jobClient.EXPECT().UpdateJob(gomock.Any(), job.UUID, &txschedulertypes.UpdateJobRequest{
			Status:      entities.StatusPending,
		})
		jobClient.EXPECT().UpdateJob(gomock.Any(), job.UUID, gomock.Any())

		err := usecase.Execute(ctx, job)
		assert.NoError(t, err)
		assert.Equal(t, job.Transaction.Hash, hash)
	})
}
