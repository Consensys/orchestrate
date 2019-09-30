package tessera

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/tessera"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/types/ethereum"
)

type MockTesseraClient struct {
	t *testing.T
}

func (tc *MockTesseraClient) AddClient(chainID string, tesseraEndpoint tessera.EnclaveEndpoint) {

}

func (tc *MockTesseraClient) StoreRaw(chainID string, rawTx []byte, privateFrom string) (txHash []byte, err error) {
	if chainID == "0" {
		return []byte(``), fmt.Errorf("mock: store raw failed")
	}
	return hexutil.MustDecode("0xabcdef"), nil
}

func (tc *MockTesseraClient) GetStatus(chainID string) (status string, err error) {
	if chainID == "0" {
		return "", fmt.Errorf("mock: get status failed")
	}
	return "", nil
}

func makeSignerContext(i int) *engine.TxContext {
	txctx := engine.NewTxContext()
	txctx.Reset()
	txctx.Logger = log.NewEntry(log.StandardLogger())

	switch i % 3 {
	case 0:
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
	case 1:
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
	case 2:
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
	}
	return txctx
}

func TestTxHashSetter(t *testing.T) {
	setter := txHashSetter(&MockTesseraClient{t})

	ROUNDS := 4
	for i := 0; i < ROUNDS; i++ {
		txctx := makeSignerContext(i)
		setter(txctx)
		assert.Emptyf(t, txctx.Envelope.Error(), fmt.Sprintf("TxHash should not be nil"))
		assert.NotNilf(t, txctx.Envelope.Tx.GetTxData().GetData(), fmt.Sprintf("TxHash should not be nil"))
	}
}
