// +build unit

package sender

import (
	"fmt"
	"reflect"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethclient/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
)

const (
	testMarkingPrivateFor      = "privateFor"
	testMarkingPrivateForError = "error"
	testMarkingChainProxyURL   = "chainURL"
)

var (
	testMarkingTxHash = ethcommon.HexToHash("0xABC0123")
	testMarkingError  = errors.HTTPConnectionError("error")
)

func TestTesseraRawMarkingTxSender(t *testing.T) {
	testSet := []struct {
		name          string
		input         func(txctx *engine.TxContext) *engine.TxContext
		expectedTxctx func(txctx *engine.TxContext) *engine.TxContext
	}{
		{
			"send quorum raw private transaction without error",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.
					SetJobType(tx.JobType_ETH_TESSERA_MARKING_TX).
					SetRaw([]byte{11, 22}).
					SetPrivateFor([]string{testMarkingPrivateFor})

				return txctx.WithContext(proxy.With(txctx.Context(), testMarkingChainProxyURL))
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.SetTxHash(testMarkingTxHash)
				return txctx
			},
		},
		{
			"error without data filled",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.
					SetJobType(tx.JobType_ETH_TESSERA_MARKING_TX).
					SetPrivateFrom(testMarkingPrivateFor)
				return txctx.WithContext(proxy.With(txctx.Context(), testMarkingChainProxyURL))
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				err := errors.DataError("no raw or privateFor filled").SetComponent(component)
				_ = txctx.Error(err)
				return txctx
			},
		},
		{
			"error in tessera client",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.
					SetJobType(tx.JobType_ETH_TESSERA_MARKING_TX).
					SetRaw([]byte{11, 22}).
					SetPrivateFor([]string{testMarkingPrivateForError})
				return txctx.WithContext(proxy.With(txctx.Context(), testMarkingChainProxyURL))
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Error(testMarkingError)
				return txctx
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockQuorumClient(mockCtrl)
	mockClient.EXPECT().SendQuorumRawPrivateTransaction(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq([]string{testMarkingPrivateFor})).
		Return(testMarkingTxHash, nil).AnyTimes()
	mockClient.EXPECT().SendQuorumRawPrivateTransaction(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq([]string{testMarkingPrivateForError})).
		Return(ethcommon.Hash{}, testMarkingError).AnyTimes()

	for _, test := range testSet {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			txctx := engine.NewTxContext()
			txctx.Logger = log.NewEntry(log.New())

			h := TesseraRawMarkingTxSender(mockClient)
			h(test.input(txctx))

			expectedTxctx := engine.NewTxContext()
			expectedTxctx.Logger = txctx.Logger
			expectedTxctx = test.expectedTxctx(test.input(expectedTxctx))

			assert.True(t, reflect.DeepEqual(txctx, expectedTxctx), fmt.Sprintf("%s: expected same input", test.name))
		})
	}
}
