package generic

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/args"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/ethereum"
)

func TestTransactionFromTxContext(t *testing.T) {
	testSet := []struct {
		name           string
		input          func(txctx *engine.TxContext) *engine.TxContext
		expectedOutput *ethtypes.Transaction
	}{
		{
			"constructor transaction",
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Envelope.Tx = &ethereum.Transaction{TxData: &ethereum.TxData{}}
				txctx.Envelope.Args = &envelope.Args{Call: &args.Call{Method: &abi.Method{Signature: "constructor"}}}
				txctx.Envelope.GetTx().GetTxData().
					SetNonce(10).
					SetValue(big.NewInt(9)).
					SetGas(8).
					SetGasPrice(big.NewInt(7)).
					SetData([]byte{6})
				return txctx
			},
			ethtypes.NewContractCreation(10, big.NewInt(9), 8, big.NewInt(7), []byte{6}),
		},
		{
			"transaction",
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Envelope.Tx = &ethereum.Transaction{TxData: &ethereum.TxData{}}
				txctx.Envelope.Args = &envelope.Args{Call: &args.Call{Method: &abi.Method{Signature: "test()"}}}
				txctx.Envelope.GetTx().GetTxData().
					SetNonce(10).
					SetValue(big.NewInt(9)).
					SetGas(8).
					SetGasPrice(big.NewInt(7)).
					SetData([]byte{6}).
					SetTo(common.HexToAddress("0x1"))
				return txctx
			},
			ethtypes.NewTransaction(10, common.HexToAddress("0x1"), big.NewInt(9), 8, big.NewInt(7), []byte{6}),
		},
	}

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			txctx := engine.NewTxContext()
			txctx.Logger = log.NewEntry(log.New())

			tx := TransactionFromTxContext(test.input(txctx))

			assert.True(t, reflect.DeepEqual(tx, test.expectedOutput), "Expected same input")
		})
	}

}
