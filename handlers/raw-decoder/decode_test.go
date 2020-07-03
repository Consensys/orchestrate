// +build unit

package rawdecoder

import (
	"math/big"
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
)

func TestRawDecoder(t *testing.T) {
	testSet := []struct {
		name          string
		input         func(txctx *engine.TxContext) *engine.TxContext
		expectedTxctx func(txctx *engine.TxContext) *engine.TxContext
	}{
		{
			"decode raw without error",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.
					MustSetRawString("0xf86780808252089488a5c2d9919e46f883eb62f7b8dd9d0cc45bc2908806f05b59d3b20000801ba0cf1f0ee7b02637e3d9c334ae5689c3e1fe102faf6c21486976b271c811098ef9a06b0bd4227f2d4fe4e59c09a47ec3770c63c165c4dce40d076bae22f780bebc50").
					SetContextLabelsValue("txMode", "raw")
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.
					SetNonce(0).
					SetGas(21000).
					SetGasPrice(big.NewInt(0)).
					MustSetToString("0x88a5C2d9919e46F883EB62F7b8Dd9d0CC45bc290").
					SetValue(big.NewInt(500000000000000000)).
					MustSetTxHashString("0x8c7ca933948c0481881f77fceb997f36fbe6227a499968ce7a3c3f6a708fda31").
					MustSetFromString("0x7357589f8e367c2C31F51242fB77B350A11830F3").
					MustSetDataString("0x")
				return txctx
			},
		},
		{
			"decode tx-scheduler raw without error",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.
					MustSetRawString("0xf86780808252089488a5c2d9919e46f883eb62f7b8dd9d0cc45bc2908806f05b59d3b20000801ba0cf1f0ee7b02637e3d9c334ae5689c3e1fe102faf6c21486976b271c811098ef9a06b0bd4227f2d4fe4e59c09a47ec3770c63c165c4dce40d076bae22f780bebc50").
					SetJobType(tx.JobType_ETH_RAW_TX)
				txctx.Envelope.ContextLabels["jobUUID"] = "randomJobUUID"
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.
					SetNonce(0).
					SetGas(21000).
					SetGasPrice(big.NewInt(0)).
					SetJobType(tx.JobType_ETH_RAW_TX).
					MustSetToString("0x88a5C2d9919e46F883EB62F7b8Dd9d0CC45bc290").
					SetValue(big.NewInt(500000000000000000)).
					MustSetTxHashString("0x8c7ca933948c0481881f77fceb997f36fbe6227a499968ce7a3c3f6a708fda31").
					MustSetFromString("0x7357589f8e367c2C31F51242fB77B350A11830F3").
					MustSetDataString("0x")
				return txctx
			},
		},
		{
			"no raw to decode",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.
					SetContextLabelsValue("txMode", "raw")
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.AppendError(errors.DataError("no raw filled - could not decode").ExtendComponent(component))
				return txctx
			},
		},
		{
			"could not decode raw decode",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.
					MustSetRawString("0xa86780808252089488a5c2d9919e46f883eb62f7b8dd9d0cc45bc2908806f05b59d3b20000801ba0cf1f0ee7b02637e3d9c334ae5689c3e1fe102faf6c21486976b271c811098ef9a06b0bd4227f2d4fe4e59c09a47ec3770c63c165c4dce40d076bae22f780bebc50").
					SetContextLabelsValue("txMode", "raw")
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.AppendError(errors.DataError("could not decode raw - got rlp: expected input list for types.txdata").ExtendComponent(component))
				return txctx
			},
		},
	}

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			txctx := engine.NewTxContext()
			txctx.Logger = log.NewEntry(log.New())

			RawDecoder(test.input(txctx))

			expectedTxctx := engine.NewTxContext()
			expectedTxctx.Logger = txctx.Logger
			expectedTxctx = test.expectedTxctx(test.input(expectedTxctx))

			assert.True(t, reflect.DeepEqual(txctx, expectedTxctx), "Expected same input")
		})
	}
}
