// +build unit

package eea

import (
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	mock2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
)

type output struct {
	sig  []byte
	hash *ethcommon.Hash
	err  error
}

const (
	chainRegistryUrl = "chainRegistryUrl"
	chainID          = 666
)

func newTxCtx(eId, txHash, txRaw, sender string) *engine.TxContext {
	txctx := engine.NewTxContext()
	txctx.Logger = log.NewEntry(log.StandardLogger())
	txctx.WithContext(proxy.With(txctx.Context(), chainRegistryUrl))
	_ = txctx.Envelope.SetID(eId).
		SetTxHash(ethcommon.HexToHash(txHash)).
		SetChainIDUint64(chainID).
		SetEEAMarkingTxNonce(0).
		SetRawString(txRaw)
	_ = txctx.Envelope.SetFromString(sender)
	return txctx
}

func TestSender_EnvelopeStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := ethtypes.NewContractCreation(0, nil, 0, nil, []byte{})
	envelopeId := utils.RandomString(12)
	txHash := ethcommon.HexToHash("0x" + utils.RandHexString(64))
	enlaveKey := ethcommon.HexToHash("0x" + utils.RandHexString(64))
	txRaw := ethcommon.HexToHash("0x" + utils.RandHexString(10))
	privTxRaw := "0x" + utils.RandHexString(10)
	txSender := ethcommon.HexToAddress("0x" + utils.RandHexString(32))
	precompiledContractAddr := ethcommon.HexToAddress("0x" + utils.RandHexString(32))

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	k := mock.NewMockKeyStore(mockCtrl)
	ec := mock2.NewMockClient(mockCtrl)
	signTx := generateSignTx(ec)

	t.Run("should execute eea signer successfully", func(t *testing.T) {
		txctx := newTxCtx(envelopeId, txHash.String(), txRaw.String(), txSender.String())
		k.EXPECT().
			SignPrivateEEATx(txctx.Context(), gomock.Any(), txSender, tx,
				gomock.AssignableToTypeOf(&types.PrivateArgs{})).
			Return([]byte(privTxRaw), &enlaveKey, nil)

		ec.EXPECT().PrivDistributeRawTransaction(txctx.Context(), chainRegistryUrl, hexutil.Encode([]byte(privTxRaw))).
			Return(enlaveKey, nil)

		ec.EXPECT().EEAPrivPrecompiledContractAddr(txctx.Context(), chainRegistryUrl).
			Return(precompiledContractAddr, nil)

		markingTxNonce, _ := txctx.Envelope.GetEEAMarkingNonce()
		markingTx := ethtypes.NewTransaction(
			markingTxNonce,
			precompiledContractAddr,
			tx.Value(),
			tx.Gas(),
			tx.GasPrice(),
			enlaveKey.Bytes(),
		)

		k.EXPECT().
			SignTx(txctx.Context(), gomock.Any(), txSender, markingTx).
			Return(txRaw.Bytes(), &txHash, nil)

		sig, hash, err := signTx(k, txctx, txSender, tx)

		assert.Nil(t, err)
		assert.Equal(t, sig, txRaw.Bytes())
		assert.Equal(t, &txHash, hash)
	})

	t.Run("should fail to execute eea signer if SignPrivateEEATx fails", func(t *testing.T) {
		expectedErr := errors.InternalError("Error")
		txctx := newTxCtx(envelopeId, txHash.String(), txRaw.String(), txSender.String())

		k.EXPECT().
			SignPrivateEEATx(txctx.Context(), gomock.Any(), txSender, tx,
				gomock.AssignableToTypeOf(&types.PrivateArgs{})).
			Return([]byte(privTxRaw), &enlaveKey, expectedErr)

		_, _, err := signTx(k, txctx, txSender, tx)

		assert.NotNil(t, err)
		assert.Equal(t, errors.FromError(err).ExtendComponent(component), expectedErr)
	})

	t.Run("should fail to execute eea signer if PrivDistributeRawTransaction fails", func(t *testing.T) {
		expectedErr := errors.InternalError("Error")
		txctx := newTxCtx(envelopeId, txHash.String(), txRaw.String(), txSender.String())

		k.EXPECT().
			SignPrivateEEATx(txctx.Context(), gomock.Any(), txSender, tx,
				gomock.AssignableToTypeOf(&types.PrivateArgs{})).
			Return([]byte(privTxRaw), &enlaveKey, nil)

		ec.EXPECT().PrivDistributeRawTransaction(txctx.Context(), chainRegistryUrl, hexutil.Encode([]byte(privTxRaw))).
			Return(enlaveKey, expectedErr)

		_, _, err := signTx(k, txctx, txSender, tx)

		assert.NotNil(t, err)
		assert.Equal(t, errors.CryptoOperationError(expectedErr.Error()).ExtendComponent(component), err)
	})

	t.Run("should fail to execute eea signer if EEAPrivPrecompiledContractAddr fails", func(t *testing.T) {
		expectedErr := errors.InternalError("Error")
		txctx := newTxCtx(envelopeId, txHash.String(), txRaw.String(), txSender.String())

		k.EXPECT().
			SignPrivateEEATx(txctx.Context(), gomock.Any(), txSender, tx,
				gomock.AssignableToTypeOf(&types.PrivateArgs{})).
			Return([]byte(privTxRaw), &enlaveKey, nil)

		ec.EXPECT().PrivDistributeRawTransaction(txctx.Context(), chainRegistryUrl, hexutil.Encode([]byte(privTxRaw))).
			Return(enlaveKey, nil)

		ec.EXPECT().EEAPrivPrecompiledContractAddr(txctx.Context(), chainRegistryUrl).
			Return(precompiledContractAddr, expectedErr)

		_, _, err := signTx(k, txctx, txSender, tx)

		assert.NotNil(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(component), err)
	})

	t.Run("should fail to execute eea signer if EEAPrivPrecompiledContractAddr fails", func(t *testing.T) {
		expectedErr := errors.InternalError("Error")
		txctx := newTxCtx(envelopeId, txHash.String(), txRaw.String(), txSender.String())

		k.EXPECT().
			SignPrivateEEATx(txctx.Context(), gomock.Any(), txSender, tx,
				gomock.AssignableToTypeOf(&types.PrivateArgs{})).
			Return([]byte(privTxRaw), &enlaveKey, nil)

		ec.EXPECT().PrivDistributeRawTransaction(txctx.Context(), chainRegistryUrl, hexutil.Encode([]byte(privTxRaw))).
			Return(enlaveKey, nil)

		ec.EXPECT().EEAPrivPrecompiledContractAddr(txctx.Context(), chainRegistryUrl).
			Return(precompiledContractAddr, nil)

		k.EXPECT().
			SignTx(txctx.Context(), gomock.Any(), txSender, gomock.Any()).
			Return(txRaw.Bytes(), &txHash, expectedErr)

		_, _, err := signTx(k, txctx, txSender, tx)

		assert.NotNil(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(component), err)
	})
}
