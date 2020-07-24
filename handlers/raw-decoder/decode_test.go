// +build unit

package rawdecoder

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
)

func TestRawDecoder(t *testing.T) {
	t.Run("should decode raw field successfully", func(t *testing.T) {
		txctx := engine.NewTxContext()
		txctx.Logger = log.NewEntry(log.New())
		_ = txctx.Envelope.
			MustSetRawString("0xf86780808252089488a5c2d9919e46f883eb62f7b8dd9d0cc45bc2908806f05b59d3b20000801ba0cf1f0ee7b02637e3d9c334ae5689c3e1fe102faf6c21486976b271c811098ef9a06b0bd4227f2d4fe4e59c09a47ec3770c63c165c4dce40d076bae22f780bebc50").
			SetJobType(tx.JobType_ETH_RAW_TX)

		RawDecoder(txctx)

		assert.Equal(t, txctx.Envelope.GetNonceString(), "0")
		assert.Equal(t, txctx.Envelope.GetGasString(), "21000")
		assert.Equal(t, txctx.Envelope.GetGasPriceString(), "0")
		assert.Equal(t, txctx.Envelope.GetToString(), "0x88a5C2d9919e46F883EB62F7b8Dd9d0CC45bc290")
		assert.Equal(t, txctx.Envelope.GetValueString(), "500000000000000000")
		assert.Equal(t, txctx.Envelope.GetTxHashString(), "0x8c7ca933948c0481881f77fceb997f36fbe6227a499968ce7a3c3f6a708fda31")
		assert.Equal(t, txctx.Envelope.GetFromString(), "0x7357589f8e367c2C31F51242fB77B350A11830F3")
		assert.Equal(t, txctx.Envelope.GetData(), "0x")
	})

	t.Run("should decode raw field successfully", func(t *testing.T) {
		txctx := engine.NewTxContext()
		txctx.Logger = log.NewEntry(log.New())
		_ = txctx.Envelope.MustSetRawString("").SetJobType(tx.JobType_ETH_RAW_TX)

		RawDecoder(txctx)

		assert.Len(t, txctx.Envelope.Errors, 1)
	})
}
