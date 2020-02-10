package tessera

import (
	"fmt"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/tx"

	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/tessera"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
)

type MockTesseraClient struct {
	t *testing.T
}

func (tc *MockTesseraClient) AddClient(chainID string, tesseraEndpoint tessera.EnclaveEndpoint) {}

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

func makeSignerContext(_ int) *engine.TxContext {
	txctx := engine.NewTxContext()
	txctx.Reset()
	txctx.Logger = log.NewEntry(log.StandardLogger())

	_ = txctx.Builder.
		SetChainIDUint64(10).
		SetData([]byte{0}).
		SetMethod(tx.Method_ETH_SENDRAWPRIVATETRANSACTION)
	return txctx
}

func TestTxHashSetter(t *testing.T) {
	setter := txHashSetter(&MockTesseraClient{t})

	ROUNDS := 4
	for i := 0; i < ROUNDS; i++ {
		txctx := makeSignerContext(i)
		setter(txctx)
		assert.Emptyf(t, txctx.Builder.GetErrors(), fmt.Sprintf("TxHash should not be nil"))
		assert.NotNilf(t, txctx.Builder.GetData(), fmt.Sprintf("TxHash should not be nil"))
	}
}
