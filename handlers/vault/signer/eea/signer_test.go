// +build unit

package eea

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/keystore/mocks"
)

type output struct {
	sig  []byte
	hash *common.Hash
	err  error
}

var addressNoError = common.HexToAddress("0x1")
var hashNoError = common.HexToHash("0x1")
var sigNoError = []byte{1}

var addressError = common.HexToAddress("0x2")
var hashError = common.HexToHash("0x2")
var sigError = []byte{2}

func TestSignTx(t *testing.T) {
	testSet := []struct {
		name           string
		txctx          func(txctx *engine.TxContext) *engine.TxContext
		sender         common.Address
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
	k := mocks.NewMockKeyStore(mockCtrl)
	k.EXPECT().
		SignPrivateEEATx(gomock.Any(), gomock.Any(), addressNoError, gomock.Any(), gomock.Any()).
		Return(sigNoError, &hashNoError, nil).
		AnyTimes()
	k.EXPECT().
		SignPrivateEEATx(gomock.Any(), gomock.Any(), addressError, gomock.Any(), gomock.Any()).
		Return(sigError, &hashError, fmt.Errorf("error")).
		AnyTimes()

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			txctx := engine.NewTxContext()
			sig, hash, err := signTx(k, test.txctx(txctx), test.sender, &ethtypes.Transaction{})

			assert.True(t, reflect.DeepEqual(test.expectedOutput.sig, sig), "Expected same sig")
			assert.True(t, reflect.DeepEqual(test.expectedOutput.hash, hash), "Expected same hash")
			assert.Equal(t, test.expectedOutput.err, err, "Expected same error")
		})
	}
}
