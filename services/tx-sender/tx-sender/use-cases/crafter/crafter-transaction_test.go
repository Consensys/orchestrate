// +build unit

package crafter

import (
	"context"
	"math/big"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	mock2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sender/tx-sender/nonce/mocks"

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

	usecase := NewCraftTransactionUseCase(ec, chainRegistryURL, nm)

	t.Run("should execute use case successfully", func(t *testing.T) {
		job := testutils.FakeJob()
		job.Transaction.Nonce = ""
		job.Transaction.GasPrice = ""
		job.Transaction.Gas = ""
		
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
	
	t.Run("should execute use case for Tessera private successfully", func(t *testing.T) {
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
		expectedGasPrice, _ := new(big.Int).SetString("1000", 10)
		expectedContractAddr := ethcommon.HexToAddress("0x1")
		ec.EXPECT().SuggestGasPrice(gomock.Any(), proxyURL).Return(expectedGasPrice, nil)
		ec.EXPECT().EstimateGas(gomock.Any(), proxyURL, gomock.Any()).Return(uint64(1000), nil)
		ec.EXPECT().EEAPrivPrecompiledContractAddr(gomock.Any(), proxyURL).Return(expectedContractAddr, nil)
		nm.EXPECT().GetNonce(gomock.Any(), gomock.Any()).Return(uint64(1), nil)
		err := usecase.Execute(ctx, job)

		assert.NoError(t, err)
		assert.Equal(t, expectedGasPrice.String(), job.Transaction.GasPrice)
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
	
	t.Run("should execute use case for child job successfully", func(t *testing.T) {
		job := testutils.FakeJob()
		job.Transaction.Nonce = ""
		job.Transaction.GasPrice = ""
		job.Transaction.Gas = ""
		job.InternalData.ParentJobUUID = job.UUID

		err := usecase.Execute(ctx, job)

		assert.NoError(t, err)
		assert.Empty(t, job.Transaction.GasPrice)
		assert.Empty(t, job.Transaction.Gas)
		assert.Empty(t, job.Transaction.Nonce)
	})
}
