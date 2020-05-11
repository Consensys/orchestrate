// +build unit

package tessera

import (
	"fmt"
	"reflect"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	ksmock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore/mock"
)

type output struct {
	sig  []byte
	hash *ethcommon.Hash
	err  error
}

var addressNoError = ethcommon.HexToAddress("0x1")
var hashNoError = ethcommon.HexToHash("0x1")
var sigNoError = []byte{1}

var addressError = ethcommon.HexToAddress("0x2")
var hashError = ethcommon.HexToHash("0x2")
var sigError = []byte{2}

func TestSignTx(t *testing.T) {
	testSet := []struct {
		name           string
		txctx          func(txctx *engine.TxContext) *engine.TxContext
		sender         ethcommon.Address
		expectedOutput output
	}{
		{
			"signTx without error",
			func(txctx *engine.TxContext) *engine.TxContext {
				return txctx
			},
			addressNoError,
			output{
				sig:  sigNoError,
				hash: &hashNoError,
				err:  nil,
			},
		},
		{
			"signTx with error",
			func(txctx *engine.TxContext) *engine.TxContext {
				return txctx
			},
			addressError,
			output{
				sig:  sigError,
				hash: &hashError,
				err:  errors.FromError(fmt.Errorf("error")).ExtendComponent(component),
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	k := ksmock.NewMockKeyStore(mockCtrl)
	k.EXPECT().
		SignPrivateTesseraTx(gomock.Any(), gomock.Any(), addressNoError, gomock.Any()).
		Return(sigNoError, &hashNoError, nil).
		AnyTimes()
	k.EXPECT().
		SignPrivateTesseraTx(gomock.Any(), gomock.Any(), addressError, gomock.Any()).
		Return(sigError, &hashError, fmt.Errorf("error")).
		AnyTimes()

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			txctx := engine.NewTxContext()
			txctx.Logger = log.NewEntry(log.New())
			sig, hash, err := signTx(k, test.txctx(txctx), test.sender, &ethtypes.Transaction{})

			assert.True(t, reflect.DeepEqual(test.expectedOutput.sig, sig), "Expected same sig")
			assert.True(t, reflect.DeepEqual(test.expectedOutput.hash, hash), "Expected same hash")
			assert.Equal(t, test.expectedOutput.err, err, "Expected same error")
		})
	}
}
