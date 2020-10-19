// +build unit

package faucet

import (
	"math/big"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/chainregistry"
	txschedulertypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/txscheduler"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client/mock"

	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	mockregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client/mock"
)

const (
	testChainUUID = "asdav-asdasd-asdasd"
	testChainName = "testChain"
)

var candidate = &types.Faucet{
	UUID:       "testUUID",
	MaxBalance: big.NewInt(10),
	Amount:     big.NewInt(10),
	Creditor:   ethcommon.HexToAddress("0xab"),
}

var (
	testSenderAddr = ethcommon.HexToAddress("0xac")
	faucetNotFoundErr = errors.NotFoundError("not found faucet candidate")
)

func newTestTxEnvelope(chainUUID, chainName string, sender ethcommon.Address) *engine.TxContext {
	txctx := engine.NewTxContext()
	txctx.Logger = log.NewEntry(log.New())
	_ = txctx.Envelope.SetChainUUID(chainUUID).SetChainName(chainName).SetFrom(sender)
	return txctx
}

func TestMaxBalanceControl_Execute(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	chainRegistryClient := mockregistry.NewMockChainRegistryClient(mockCtrl)
	txSchedulerClient := mock.NewMockTransactionSchedulerClient(mockCtrl)
	h := Faucet(chainRegistryClient, txSchedulerClient)

	t.Run("should trigger a new faucet transaction, with chainUUID, successfully", func(t *testing.T) {
		txctx := newTestTxEnvelope(testChainUUID, "", testSenderAddr)
		chainRegistryClient.EXPECT().GetChainByUUID(gomock.Any(), testChainUUID).Return(&models.Chain{
			UUID: testChainUUID,
			Name: testChainName,
		}, nil)

		chainRegistryClient.EXPECT().GetFaucetCandidate(gomock.Any(), testSenderAddr, testChainUUID).Return(candidate, nil)
		txSchedulerClient.EXPECT().SendTransferTransaction(gomock.Any(), &txschedulertypes.TransferRequest{
			ChainName: testChainName,
			Params: txschedulertypes.TransferParams{
				From:  candidate.Creditor.Hex(),
				To:    txctx.Envelope.MustGetFromAddress().String(),
				Value: candidate.Amount.String(),
			},
			Labels: types.FaucetToJobLabels(candidate),
		})

		h(txctx)

		assert.Len(t, txctx.Envelope.GetErrors(), 0)
	})
	
	t.Run("should trigger a new faucet transaction, with chanName, successfully", func(t *testing.T) {
		txctx := newTestTxEnvelope("", testChainName, testSenderAddr)
		chainRegistryClient.EXPECT().GetChainByName(gomock.Any(), testChainName).Return(&models.Chain{
			UUID: testChainUUID,
			Name: testChainName,
		}, nil)

		chainRegistryClient.EXPECT().GetFaucetCandidate(gomock.Any(), testSenderAddr, testChainUUID).Return(candidate, nil)
		txSchedulerClient.EXPECT().SendTransferTransaction(gomock.Any(), &txschedulertypes.TransferRequest{
			ChainName: testChainName,
			Params: txschedulertypes.TransferParams{
				From:  candidate.Creditor.Hex(),
				To:    txctx.Envelope.MustGetFromAddress().String(),
				Value: candidate.Amount.String(),
			},
			Labels: types.FaucetToJobLabels(candidate),
		})

		h(txctx)

		assert.Len(t, txctx.Envelope.GetErrors(), 0)
	})
	
	t.Run("should fail in case if fails to fetch faucet candidates", func(t *testing.T) {
		expectedErr := errors.ConnectionError("cannot retrieve faucet candidates")
		txctx := newTestTxEnvelope(testChainUUID, "", testSenderAddr)
		chainRegistryClient.EXPECT().GetChainByUUID(gomock.Any(), testChainUUID).Return(&models.Chain{
			UUID: testChainUUID,
			Name: testChainName,
		}, nil)

		chainRegistryClient.EXPECT().GetFaucetCandidate(gomock.Any(), testSenderAddr, testChainUUID).Return(nil, expectedErr)

		h(txctx)

		assert.Len(t, txctx.Envelope.GetErrors(), 1)
		err := txctx.Envelope.GetErrors()[0]
		assert.Equal(t, err, expectedErr.ExtendComponent(component))
	})
	
	t.Run("should ignore in case there is not available candidates", func(t *testing.T) {
		txctx := newTestTxEnvelope(testChainUUID, "", testSenderAddr)
		chainRegistryClient.EXPECT().GetChainByUUID(gomock.Any(), testChainUUID).Return(&models.Chain{
			UUID: testChainUUID,
			Name: testChainName,
		}, nil)

		chainRegistryClient.EXPECT().GetFaucetCandidate(gomock.Any(), testSenderAddr, testChainUUID).Return(nil, faucetNotFoundErr)

		h(txctx)

		assert.Len(t, txctx.Envelope.GetErrors(), 0)
	})
	
	t.Run("should fail to send faucet transaction", func(t *testing.T) {
		expectedErr := errors.ConnectionError("cannot reach tx-scheduler service")
		txctx := newTestTxEnvelope(testChainUUID, "", testSenderAddr)
		chainRegistryClient.EXPECT().GetChainByUUID(gomock.Any(), testChainUUID).Return(&models.Chain{
			UUID: testChainUUID,
			Name: testChainName,
		}, nil)

		chainRegistryClient.EXPECT().GetFaucetCandidate(gomock.Any(), testSenderAddr, testChainUUID).Return(candidate, nil)
		txSchedulerClient.EXPECT().SendTransferTransaction(gomock.Any(), &txschedulertypes.TransferRequest{
			ChainName: testChainName,
			Params: txschedulertypes.TransferParams{
				From:  candidate.Creditor.Hex(),
				To:    txctx.Envelope.MustGetFromAddress().String(),
				Value: candidate.Amount.String(),
			},
			Labels: types.FaucetToJobLabels(candidate),
		}).Return(nil, expectedErr)

		h(txctx)

		assert.Len(t, txctx.Envelope.GetErrors(), 1)
		err := txctx.Envelope.GetErrors()[0]
		assert.Equal(t, err, expectedErr.ExtendComponent(component))
	})
}
