package sender

import (
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethclient/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
)

func TestEEAPrivateTxSender(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ec := mock.NewMockEEATransactionSender(ctrl)
	h := EEAPrivateTxSender(ec)
	rawData := ethcommon.HexToHash("0x" + utils.RandHexString(10))
	privTxHash := ethcommon.HexToHash("0x123")
	testPrivChainProxyURL := "testEEAPrivChainProxyURL"

	t.Run("send private eea transaction successfully", func(t *testing.T) {
		txctx := engine.NewTxContext()
		txctx.Logger = log.NewEntry(log.New())
		_ = txctx.Envelope.SetJobType(tx.JobType_ETH_ORION_EEA_TX).
			SetRawString(rawData.String())
		txctx.WithContext(proxy.With(txctx.Context(), testPrivChainProxyURL))

		ec.EXPECT().PrivDistributeRawTransaction(gomock.Any(), testPrivChainProxyURL, rawData.String()).
			Return(privTxHash, nil)

		h(txctx)

		assert.Len(t, txctx.Envelope.GetErrors(), 0)
		assert.Equal(t, privTxHash.String(), txctx.Envelope.GetTxHashString())
	})

	t.Run("failt to send private eea transaction when ethclient fails", func(t *testing.T) {
		txctx := engine.NewTxContext()
		err := errors.ConnectionError("error")
		txctx.Logger = log.NewEntry(log.New())
		_ = txctx.Envelope.SetJobType(tx.JobType_ETH_ORION_EEA_TX).
			SetRawString(rawData.String())
		txctx.WithContext(proxy.With(txctx.Context(), testPrivChainProxyURL))

		ec.EXPECT().PrivDistributeRawTransaction(gomock.Any(), testPrivChainProxyURL, rawData.String()).
			Return(privTxHash, err)

		h(txctx)

		assert.Len(t, txctx.Envelope.GetErrors(), 1)
		assert.Equal(t, txctx.Envelope.GetErrors()[0], err)
		assert.Equal(t, "", txctx.Envelope.GetTxHashString())
	})

	t.Run("fail to send private eea transaction without raw data", func(t *testing.T) {
		txctx := engine.NewTxContext()
		txctx.Logger = log.NewEntry(log.New())
		_ = txctx.Envelope.SetJobType(tx.JobType_ETH_ORION_EEA_TX)
		txctx.WithContext(proxy.With(txctx.Context(), testPrivChainProxyURL))

		h(txctx)

		assert.Len(t, txctx.Envelope.GetErrors(), 1)
		assert.Equal(t, "", txctx.Envelope.GetTxHashString())
	})
}
