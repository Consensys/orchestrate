// +build unit

package crafter

import (
	"context"
	"math/big"
	"testing"

	mock2 "github.com/consensys/orchestrate/pkg/toolkit/ethclient/mock"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/types/testutils"
	"github.com/consensys/orchestrate/pkg/types/tx"
	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/consensys/orchestrate/services/tx-sender/tx-sender/nonce/mocks"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCrafterTransaction_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	ec := mock2.NewMockMultiClient(ctrl)
	nm := mocks.NewMockManager(ctrl)
	chainRegistryURL := "http://chain-registry:8081"

	nextBaseFee, _ := new(big.Int).SetString("1000000000", 10)
	mediumPriority, _ := new(big.Int).SetString(mediumPriorityString, 10)

	usecase := NewCraftTransactionUseCase(ec, chainRegistryURL, nm)

	t.Run("should execute use case for LegacyTx successfully", func(t *testing.T) {
		job := testutils.FakeJob()
		job.Transaction.Nonce = ""
		job.Transaction.GasPrice = ""
		job.Transaction.Gas = ""
		job.Transaction.TransactionType = entities.LegacyTxType

		proxyURL := utils.GetProxyURL(chainRegistryURL, job.ChainUUID)
		expectedGasPrice, _ := new(big.Int).SetString("1000", 10)
		ec.EXPECT().SuggestGasPrice(gomock.Any(), proxyURL).Return(expectedGasPrice, nil)
		ec.EXPECT().EstimateGas(gomock.Any(), proxyURL, gomock.Any()).Return(uint64(1000), nil)
		nm.EXPECT().GetNonce(gomock.Any(), gomock.Any()).Return(uint64(1), nil)
		err := usecase.Execute(ctx, job)

		assert.NoError(t, err)
		assert.Equal(t, expectedGasPrice.String(), job.Transaction.GasPrice)
		assert.Equal(t, "1000", job.Transaction.Gas)
		assert.Equal(t, "1", job.Transaction.Nonce)
	})

	t.Run("should execute use case for DynamicFeeTx successfully", func(t *testing.T) {
		job := testutils.FakeJob()
		job.Transaction.Nonce = ""
		job.Transaction.Gas = ""
		job.Transaction.GasPrice = ""

		proxyURL := utils.GetProxyURL(chainRegistryURL, job.ChainUUID)
		expectedFeeHistory := testutils.FakeFeeHistory(nextBaseFee)
		ec.EXPECT().FeeHistory(gomock.Any(), proxyURL, 1, "latest").Return(expectedFeeHistory, nil)
		ec.EXPECT().EstimateGas(gomock.Any(), proxyURL, gomock.Any()).Return(uint64(1000), nil)
		nm.EXPECT().GetNonce(gomock.Any(), gomock.Any()).Return(uint64(1), nil)
		err := usecase.Execute(ctx, job)

		expectedFeeCap := new(big.Int).Add(mediumPriority, nextBaseFee)

		assert.NoError(t, err)
		assert.Equal(t, entities.DynamicFeeTxType, job.Transaction.TransactionType)
		assert.Equal(t, expectedFeeCap.String(), job.Transaction.GasFeeCap)
		assert.Equal(t, mediumPriority.String(), job.Transaction.GasTipCap)
		assert.Equal(t, "1000", job.Transaction.Gas)
		assert.Equal(t, "1", job.Transaction.Nonce)
	})

	t.Run("should execute use case for OneTimeKey successfully", func(t *testing.T) {
		job := testutils.FakeJob()
		job.Transaction.Nonce = ""
		job.Transaction.GasPrice = ""
		job.Transaction.Gas = ""
		job.InternalData.OneTimeKey = true

		proxyURL := utils.GetProxyURL(chainRegistryURL, job.ChainUUID)
		expectedGasPrice, _ := new(big.Int).SetString("1000", 10)
		ec.EXPECT().SuggestGasPrice(gomock.Any(), proxyURL).Return(expectedGasPrice, nil)
		ec.EXPECT().EstimateGas(gomock.Any(), proxyURL, gomock.Any()).Return(uint64(1000), nil)
		err := usecase.Execute(ctx, job)

		assert.NoError(t, err)
		assert.Equal(t, expectedGasPrice.String(), job.Transaction.GasPrice)
		assert.Equal(t, "1000", job.Transaction.Gas)
		assert.Equal(t, "0", job.Transaction.Nonce)
	})

	t.Run("should execute use case for EEA marking transaction successfully", func(t *testing.T) {
		job := testutils.FakeJob()
		job.Type = entities.JobType(tx.JobType_ETH_ORION_MARKING_TX.String())
		job.Transaction.Nonce = ""
		job.Transaction.GasPrice = ""
		job.Transaction.Gas = ""

		proxyURL := utils.GetProxyURL(chainRegistryURL, job.ChainUUID)
		expectedContractAddr := ethcommon.HexToAddress("0x1")
		expectedFeeHistory := testutils.FakeFeeHistory(nextBaseFee)
		ec.EXPECT().FeeHistory(gomock.Any(), proxyURL, 1, "latest").Return(expectedFeeHistory, nil)
		ec.EXPECT().EstimateGas(gomock.Any(), proxyURL, gomock.Any()).Return(uint64(1000), nil)
		ec.EXPECT().EEAPrivPrecompiledContractAddr(gomock.Any(), proxyURL).Return(expectedContractAddr, nil)
		nm.EXPECT().GetNonce(gomock.Any(), gomock.Any()).Return(uint64(1), nil)
		err := usecase.Execute(ctx, job)

		assert.NoError(t, err)
		assert.Equal(t, "1000", job.Transaction.Gas)
		assert.Equal(t, "1", job.Transaction.Nonce)
		assert.Equal(t, expectedContractAddr.String(), job.Transaction.To)
	})

	t.Run("should execute use case for EEA private transaction successfully", func(t *testing.T) {
		job := testutils.FakeJob()
		job.Type = entities.JobType(tx.JobType_ETH_ORION_EEA_TX.String())
		job.Transaction.Nonce = ""
		job.Transaction.GasPrice = ""
		job.Transaction.Gas = ""

		nm.EXPECT().GetNonce(gomock.Any(), gomock.Any()).Return(uint64(1), nil)
		err := usecase.Execute(ctx, job)

		assert.NoError(t, err)
		assert.Empty(t, job.Transaction.GasPrice)
		assert.Empty(t, job.Transaction.Gas)
		assert.Equal(t, "1", job.Transaction.Nonce)
	})

	t.Run("should execute use case for child job for DynamicTx successfully", func(t *testing.T) {
		job := testutils.FakeJob()
		job.Transaction.GasTipCap = "199999"
		job.Transaction.TransactionType = entities.DynamicFeeTxType
		job.InternalData.ParentJobUUID = job.UUID

		proxyURL := utils.GetProxyURL(chainRegistryURL, job.ChainUUID)
		expectedFeeHistory := testutils.FakeFeeHistory(nextBaseFee)
		ec.EXPECT().FeeHistory(gomock.Any(), proxyURL, 1, "latest").Return(expectedFeeHistory, nil)

		err := usecase.Execute(ctx, job)
		priority, _ := new(big.Int).SetString(job.Transaction.GasTipCap, 10)
		expectedFeeCap := new(big.Int).Add(priority, nextBaseFee)

		assert.NoError(t, err)
		assert.Equal(t, expectedFeeCap.String(), job.Transaction.GasFeeCap)
	})
}
