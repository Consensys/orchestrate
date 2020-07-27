// +build unit

package signer

import (
	"math/big"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	eeaHandlers "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/vault/signer/eea"
	ethereumHandlers "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/vault/signer/ethereum"
	tesseraHandlers "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/vault/signer/tessera"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/types"
	mock3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
)

const (
	chainRegistryUrl = "chainRegistryUrl"
	chainID          = 666
)

func newTxCtx(eId, sender string) *engine.TxContext {
	txctx := engine.NewTxContext()
	txctx.Logger = log.NewEntry(log.StandardLogger())
	txctx.WithContext(proxy.With(txctx.Context(), chainRegistryUrl))
	_ = txctx.Envelope.SetID(eId).
		SetChainIDUint64(chainID).
		SetGas(0).
		SetGasPrice(big.NewInt(0)).
		SetEEAMarkingTxNonce(0).
		SetNonce(0)
	_ = txctx.Envelope.SetFromString(sender)
	return txctx
}

func TestSender(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	envelopeId := utils.RandomString(12)
	txHash := ethcommon.HexToHash("0x" + utils.RandHexString(64))
	txRaw := ethcommon.HexToHash("0x" + utils.RandHexString(10))
	txSender := ethcommon.HexToAddress("0x" + utils.RandHexString(32))
	contractAddr := ethcommon.HexToAddress("0x" + utils.RandHexString(32))

	ec := mock.NewMockClient(ctrl)
	ks := mock3.NewMockKeyStore(ctrl)
	signer := TxSigner(
		ethereumHandlers.Signer(ks, ks),
		eeaHandlers.Signer(ks, ks, ec),
		tesseraHandlers.Signer(ks, ks),
	)

	t.Run("should execute raw transaction successfully", func(t *testing.T) {
		txctx := newTxCtx(envelopeId, txSender.String())
		_ = txctx.Envelope.SetJobType(tx.JobType_ETH_RAW_TX)

		ks.EXPECT().SignTx(txctx.Context(), gomock.Any(), txSender, gomock.AssignableToTypeOf(&ethtypes.Transaction{})).
			Return(txRaw.Bytes(), &txHash, nil)

		signer(txctx)

		assert.Empty(t, txctx.Envelope.GetErrors())
		assert.Equal(t, txctx.Envelope.GetRaw(), txRaw.Hex())
		assert.Equal(t, txctx.Envelope.GetTxHash(), &txHash)
	})

	t.Run("should fail to execute raw transaction successfully", func(t *testing.T) {
		expectedErr := errors.InternalError("Error")
		txctx := newTxCtx(envelopeId, txSender.String())
		_ = txctx.Envelope.SetJobType(tx.JobType_ETH_RAW_TX)

		ks.EXPECT().SignTx(txctx.Context(), gomock.Any(), txSender, gomock.AssignableToTypeOf(&ethtypes.Transaction{})).
			Return(txRaw.Bytes(), &txHash, expectedErr)

		signer(txctx)

		assert.NotEmpty(t, txctx.Envelope.GetErrors())
		assert.Equal(t, expectedErr, txctx.Envelope.GetErrors()[0])
	})

	t.Run("should execute tessera transaction successfully", func(t *testing.T) {
		txctx := newTxCtx(envelopeId, txSender.String())
		_ = txctx.Envelope.SetJobType(tx.JobType_ETH_TESSERA_PRIVATE_TX)

		ks.EXPECT().SignPrivateTesseraTx(txctx.Context(), gomock.Any(), txSender, gomock.AssignableToTypeOf(&ethtypes.Transaction{})).
			Return(txRaw.Bytes(), &txHash, nil)

		signer(txctx)

		assert.Empty(t, txctx.Envelope.GetErrors())
		assert.Equal(t, txctx.Envelope.GetRaw(), txRaw.Hex())
		assert.Equal(t, txctx.Envelope.GetTxHash(), &txHash)
	})

	t.Run("should fail to execute tessera transaction successfully", func(t *testing.T) {
		expectedErr := errors.InternalError("Error")
		txctx := newTxCtx(envelopeId, txSender.String())
		_ = txctx.Envelope.SetJobType(tx.JobType_ETH_TESSERA_PRIVATE_TX)

		ks.EXPECT().SignPrivateTesseraTx(txctx.Context(), gomock.Any(), txSender, gomock.AssignableToTypeOf(&ethtypes.Transaction{})).
			Return(txRaw.Bytes(), &txHash, expectedErr)

		signer(txctx)

		assert.NotEmpty(t, txctx.Envelope.GetErrors())
		assert.Equal(t, expectedErr, txctx.Envelope.GetErrors()[0])
	})

	t.Run("should execute eea transaction successfully", func(t *testing.T) {
		txctx := newTxCtx(envelopeId, txSender.String())
		_ = txctx.Envelope.SetJobType(tx.JobType_ETH_ORION_EEA_TX)

		ks.EXPECT().SignPrivateEEATx(txctx.Context(), gomock.Any(), txSender,
			gomock.AssignableToTypeOf(&ethtypes.Transaction{}), gomock.AssignableToTypeOf(&types.PrivateArgs{})).
			Return(txRaw.Bytes(), &txHash, nil)

		ec.EXPECT().PrivDistributeRawTransaction(txctx.Context(), chainRegistryUrl, gomock.Any()).
			Return(txHash, nil)

		ec.EXPECT().EEAPrivPrecompiledContractAddr(txctx.Context(), chainRegistryUrl).
			Return(contractAddr, nil)

		ks.EXPECT().
			SignTx(txctx.Context(), gomock.Any(), txSender, gomock.Any()).
			Return(txRaw.Bytes(), &txHash, nil)

		signer(txctx)

		assert.Empty(t, txctx.Envelope.GetErrors())
		assert.Equal(t, txctx.Envelope.GetRaw(), txRaw.Hex())
		assert.Equal(t, txctx.Envelope.GetTxHash(), &txHash)
	})

	t.Run("should fail to execute eea transaction successfully", func(t *testing.T) {
		expectedErr := errors.InternalError("Error")
		txctx := newTxCtx(envelopeId, txSender.String())
		_ = txctx.Envelope.SetJobType(tx.JobType_ETH_ORION_EEA_TX)

		ks.EXPECT().SignPrivateEEATx(txctx.Context(), gomock.Any(), txSender,
			gomock.AssignableToTypeOf(&ethtypes.Transaction{}), gomock.AssignableToTypeOf(&types.PrivateArgs{})).
			Return(txRaw.Bytes(), &txHash, expectedErr)

		signer(txctx)

		assert.NotEmpty(t, txctx.Envelope.GetErrors())
		assert.Equal(t, expectedErr, txctx.Envelope.GetErrors()[0])
	})
}
