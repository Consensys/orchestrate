package generic

import (
	"fmt"
	"math/big"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/tx"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/keystore"
)

func mockSignerFunc(keystore.KeyStore, *engine.TxContext, ethcommon.Address, *ethtypes.Transaction) ([]byte, *ethcommon.Hash, error) {
	return []byte{}, &ethcommon.Hash{}, nil
}

var alreadySignedTx = "0x00"

func makeSignerContext(i int) *engine.TxContext {
	txctx := engine.NewTxContext()
	txctx.Reset()
	txctx.Logger = log.NewEntry(log.StandardLogger())

	switch i % 5 {
	case 0:
		_ = txctx.Envelope.
			SetTxHash(ethcommon.HexToHash("0x12345678")).
			SetChainIDUint64(10).
			SetNonce(10).
			SetValue(big.NewInt(100)).
			MustSetToString("0x1").
			MustSetFromString("0x2").
			SetGas(11).
			SetGasPrice(big.NewInt(12))
	case 1:
		_ = txctx.Envelope.
			SetTxHash(ethcommon.HexToHash("0x12345678")).
			SetChainIDUint64(0).
			SetNonce(10).
			SetValue(big.NewInt(100)).
			MustSetToString("0x1").
			MustSetFromString("0x2").
			SetGas(11).
			SetGasPrice(big.NewInt(12))
	case 2:
		_ = txctx.Envelope.SetChainIDUint64(0).
			SetChainIDUint64(0).
			SetNonce(10).
			SetValue(big.NewInt(100)).
			MustSetToString("0x1").
			MustSetFromString("0x2").
			SetGas(11).
			SetGasPrice(big.NewInt(12))
	case 3:
		_ = txctx.Envelope.SetChainIDUint64(10).
			SetChainIDUint64(10).
			SetNonce(10).
			SetValue(big.NewInt(100)).
			MustSetToString("0x1").
			MustSetFromString("0x2").
			SetGas(11).
			SetGasPrice(big.NewInt(12))
	case 4:
		_ = txctx.Envelope.
			SetChainIDUint64(10).
			MustSetDataString("0").
			SetMethod(tx.Method_ETH_SENDRAWPRIVATETRANSACTION).
			SetNonce(10).
			SetValue(big.NewInt(100)).
			MustSetToString("0x1").
			MustSetFromString("0x2").
			SetGas(11).
			SetGasPrice(big.NewInt(12))
	case 5:
		_ = txctx.Envelope.
			SetChainIDUint64(10).
			SetMethod(tx.Method_ETH_SENDRAWPRIVATETRANSACTION).
			SetNonce(10).
			MustSetRawString(alreadySignedTx).
			SetValue(big.NewInt(100)).
			MustSetToString("0x1").
			MustSetFromString("0x2").
			SetGas(11).
			SetGasPrice(big.NewInt(12))
	case 6:
		_ = txctx.Envelope.
			SetChainIDUint64(0).
			MustSetRawString("0").
			SetMethod(tx.Method_ETH_SENDRAWPRIVATETRANSACTION).
			SetNonce(10).
			SetValue(big.NewInt(100)).
			MustSetToString("0x1").
			MustSetFromString("0x2").
			SetGas(11).
			SetGasPrice(big.NewInt(12))
	case 7:
		_ = txctx.Envelope.
			SetChainIDUint64(0).
			SetMethod(tx.Method_EEA_SENDPRIVATETRANSACTION).
			SetNonce(10).
			SetValue(big.NewInt(100)).
			MustSetToString("0x1").
			MustSetFromString("0x2").
			SetGas(11).
			SetGasPrice(big.NewInt(12))
	}
	return txctx
}

func TestGeneric(t *testing.T) {
	// Just checking the signer is properly generated
	handler := GenerateSignerHandler(
		mockSignerFunc,
		keystore.GlobalKeyStore(),
		"A success message",
		"An error message",
	)

	ROUNDS := 100
	for i := 0; i < ROUNDS; i++ {
		txctx := makeSignerContext(i)
		handler(txctx)
		assert.NotNilf(t, txctx.Envelope.GetRaw(), fmt.Sprintf("TxRawSignature should not be nil"))
		assert.NotNilf(t, txctx.Envelope.GetTxHash(), fmt.Sprintf("TxHash should not be nil"))

	}
}
