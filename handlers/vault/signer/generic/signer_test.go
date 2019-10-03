package generic

import (
	"fmt"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/multi-vault/keystore"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/types/ethereum"
)

func mockSignerFunc(keystore.KeyStore, *engine.TxContext, ethcommon.Address, *ethtypes.Transaction) ([]byte, *ethcommon.Hash, error) {
	return []byte{}, &ethcommon.Hash{}, nil
}

var alreadySignedTx = "0x00"

func makeSignerContext(i int) *engine.TxContext {
	txctx := engine.NewTxContext()
	txctx.Reset()
	txctx.Logger = log.NewEntry(log.StandardLogger())

	switch i % 8 {
	case 0:
		h := ethcommon.HexToHash("0x12345678")
		txctx.Envelope.Chain = chain.FromInt(10)
		txctx.Envelope.Tx = &ethereum.Transaction{
			Raw:  ethereum.HexToData(alreadySignedTx),
			Hash: ethereum.NewHash(h.Bytes()),
		}
	case 1:
		h := ethcommon.HexToHash("0x12345678")
		txctx.Envelope.Chain = chain.FromInt(0)
		txctx.Envelope.Tx = &ethereum.Transaction{
			Raw:  ethereum.HexToData(alreadySignedTx),
			Hash: ethereum.NewHash(h.Bytes()),
		}
	case 2:
		txctx.Envelope.Chain = chain.FromInt(0)
		txctx.Envelope.Tx = &ethereum.Transaction{}
	case 3:
		txctx.Envelope.Chain = chain.FromInt(10)
		txctx.Envelope.Tx = &ethereum.Transaction{}
	case 4:
		txctx.Envelope.Chain = chain.FromInt(10)
		txctx.Envelope.Tx = &ethereum.Transaction{
			TxData: &ethereum.TxData{
				Data: &ethereum.Data{
					Raw: []byte{0},
				},
			},
		}
		txctx.Envelope.Protocol = &chain.Protocol{
			Type: chain.ProtocolType_QUORUM_TESSERA,
		}
	case 5:
		txctx.Envelope.Chain = chain.FromInt(10)
		txctx.Envelope.Tx = &ethereum.Transaction{}
		txctx.Envelope.Protocol = &chain.Protocol{
			Type: chain.ProtocolType_QUORUM_TESSERA,
		}
	case 6:
		txctx.Envelope.Chain = chain.FromInt(0)
		txctx.Envelope.Tx = &ethereum.Transaction{
			TxData: &ethereum.TxData{
				Data: &ethereum.Data{
					Raw: []byte{0},
				},
			},
		}
		txctx.Envelope.Protocol = &chain.Protocol{
			Type: chain.ProtocolType_QUORUM_TESSERA,
		}
	case 7:
		txctx.Envelope.Chain = chain.FromInt(10)
		txctx.Envelope.Tx = &ethereum.Transaction{}
		txctx.Envelope.Protocol = &chain.Protocol{
			Type: chain.ProtocolType_PANTHEON_ORION,
		}
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
		assert.NotNilf(t, txctx.Envelope.Tx.GetRaw(), fmt.Sprintf("TxRawSignature should not be nil"))
		assert.NotNilf(t, txctx.Envelope.Tx.GetHash(), fmt.Sprintf("TxHash should not be nil"))

	}
}
