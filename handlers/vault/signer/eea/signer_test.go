// +build unit

package eea

import (
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
)

const component = "handler.signer.eea"

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
		SetNonce(0).
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
	// precompiledContractAddr := ethcommon.HexToAddress("0x" + utils.RandHexString(32))

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	k := mock.NewMockKeyStore(mockCtrl)
	signTx := generateSignTx()

	t.Run("should execute eea signer successfully", func(t *testing.T) {
		txctx := newTxCtx(envelopeId, txHash.String(), txRaw.String(), txSender.String())

		k.EXPECT().
			SignPrivateEEATx(txctx.Context(), gomock.Any(),  gomock.Any(), tx,
				gomock.AssignableToTypeOf(&types.PrivateArgs{})).
			Return(txRaw.Bytes(), nil, nil)

		sig, hash, err := signTx(k, txctx, txSender, tx)

		assert.Nil(t, err)
		assert.Equal(t, hash, &ethcommon.Hash{})
		assert.Equal(t, sig, txRaw.Bytes())
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
}
