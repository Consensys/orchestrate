// +build unit

package gaspricer

import (
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/proxy"
	"math/big"
	"testing"
)

func TestPricer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEthClient := mock.NewMockGasPricer(ctrl)
	pricer := Pricer(mockEthClient)

	t.Run("should do nothing if gasPrice is set", func(t *testing.T) {
		txctx := engine.NewTxContext()
		txctx.Reset()
		txctx.Logger = log.NewEntry(log.StandardLogger())
		txctx.WithContext(proxy.With(txctx.Context(), "URL"))
		_ = txctx.Envelope.SetGasPriceString("1000000")

		pricer(txctx)

		assert.Equal(t, txctx.Envelope.GetGasPriceString(), "1000000")
	})

	t.Run("should set gas price with no coefficient if priority not specified", func(t *testing.T) {
		txctx := engine.NewTxContext()
		txctx.Reset()
		txctx.Logger = log.NewEntry(log.StandardLogger())
		txctx.WithContext(proxy.With(txctx.Context(), "URL"))

		mockEthClient.EXPECT().SuggestGasPrice(gomock.Any(), "URL").Return(big.NewInt(10), nil)

		pricer(txctx)

		assert.Equal(t, txctx.Envelope.GetGasPriceString(), "10")
	})

	t.Run("should set gas price with coefficient 0.6 if priority is very low", func(t *testing.T) {
		txctx := engine.NewTxContext()
		txctx.Reset()
		txctx.Logger = log.NewEntry(log.StandardLogger())
		txctx.WithContext(proxy.With(txctx.Context(), "URL"))
		_ = txctx.Envelope.SetContextLabelsValue("priority", utils.PriorityVeryLow)

		mockEthClient.EXPECT().SuggestGasPrice(gomock.Any(), "URL").Return(big.NewInt(10), nil)

		pricer(txctx)

		assert.Equal(t, txctx.Envelope.GetGasPriceString(), "6")
	})

	t.Run("should set gas price with coefficient 0.8 if priority is low", func(t *testing.T) {
		txctx := engine.NewTxContext()
		txctx.Reset()
		txctx.Logger = log.NewEntry(log.StandardLogger())
		txctx.WithContext(proxy.With(txctx.Context(), "URL"))
		_ = txctx.Envelope.SetContextLabelsValue("priority", utils.PriorityLow)

		mockEthClient.EXPECT().SuggestGasPrice(gomock.Any(), "URL").Return(big.NewInt(10), nil)

		pricer(txctx)

		assert.Equal(t, txctx.Envelope.GetGasPriceString(), "8")
	})

	t.Run("should set gas price with coefficient 1 if priority is medium", func(t *testing.T) {
		txctx := engine.NewTxContext()
		txctx.Reset()
		txctx.Logger = log.NewEntry(log.StandardLogger())
		txctx.WithContext(proxy.With(txctx.Context(), "URL"))
		_ = txctx.Envelope.SetContextLabelsValue("priority", utils.PriorityMedium)

		mockEthClient.EXPECT().SuggestGasPrice(gomock.Any(), "URL").Return(big.NewInt(10), nil)

		pricer(txctx)

		assert.Equal(t, txctx.Envelope.GetGasPriceString(), "10")
	})

	t.Run("should set gas price with coefficient 1.2 if priority is high", func(t *testing.T) {
		txctx := engine.NewTxContext()
		txctx.Reset()
		txctx.Logger = log.NewEntry(log.StandardLogger())
		txctx.WithContext(proxy.With(txctx.Context(), "URL"))
		_ = txctx.Envelope.SetContextLabelsValue("priority", utils.PriorityHigh)

		mockEthClient.EXPECT().SuggestGasPrice(gomock.Any(), "URL").Return(big.NewInt(10), nil)

		pricer(txctx)

		assert.Equal(t, txctx.Envelope.GetGasPriceString(), "12")
	})

	t.Run("should set gas price with coefficient 1.4 if priority is very high", func(t *testing.T) {
		txctx := engine.NewTxContext()
		txctx.Reset()
		txctx.Logger = log.NewEntry(log.StandardLogger())
		txctx.WithContext(proxy.With(txctx.Context(), "URL"))
		_ = txctx.Envelope.SetContextLabelsValue("priority", utils.PriorityVeryHigh)

		mockEthClient.EXPECT().SuggestGasPrice(gomock.Any(), "URL").Return(big.NewInt(10), nil)

		pricer(txctx)

		assert.Equal(t, txctx.Envelope.GetGasPriceString(), "14")
	})
}
