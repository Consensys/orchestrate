package proxy

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

type outputGetURL struct {
	proxy string
	err   error
}

func TestGetURL(t *testing.T) {
	testSet := []struct {
		name           string
		txctx          func(txctx *engine.TxContext) *engine.TxContext
		expectedOutput outputGetURL
	}{
		{
			"Check if ChainURLCtxKey has already been injected",
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.WithContext(With(txctx.Context(), "test"))
				return txctx
			},
			outputGetURL{
				proxy: "test",
			},
		},
		{
			"Check if ChainURLCtxKey has already been injected with error",
			func(txctx *engine.TxContext) *engine.TxContext {
				return txctx
			},
			outputGetURL{
				err: errors.InternalError("chain proxy url not found"),
			},
		},
	}

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			txctx := engine.NewTxContext()
			txctx.Logger = log.NewEntry(log.New())

			proxy, err := GetURL(test.txctx(txctx))

			assert.Equal(t, test.expectedOutput.err, err, "should get the correct error")
			assert.Equal(t, test.expectedOutput.proxy, proxy, "should get the correct proxy url")
		})
	}
}
