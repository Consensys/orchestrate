package storer

import (
	"math/big"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope"
	envstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope-store"
	clientmock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope-store/client/mocks"
)

func TestRawTxStore(t *testing.T) {
	testSet := []struct {
		name           string
		input          func(txctx *engine.TxContext) *engine.TxContext
		expectedStatus envstore.Status
	}{
		{
			"Store",
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Envelope.Chain = (&chain.Chain{}).SetID(big.NewInt(1))
				txctx.Envelope.Metadata = &envelope.Metadata{Id: "test"}
				return txctx
			},
			envstore.Status_PENDING,
		},
		{
			"Store envelope without Metadata ID",
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Envelope.Chain = (&chain.Chain{}).SetID(big.NewInt(1))
				txctx.Envelope.Metadata = &envelope.Metadata{Id: "test"}
				err := errors.InternalError("error").ExtendComponent(component)
				txctx.Envelope.Errors = append(txctx.Envelope.Errors, err)
				return txctx
			},
			envstore.Status_ERROR,
		},
	}

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			txctx := engine.NewTxContext()
			txctx.Logger = log.NewEntry(log.New())

			registry := clientmock.New()
			h := RawTxStore(registry)
			h(test.input(txctx))

			e, _ := registry.LoadByID(txctx.Context(), &envstore.LoadByIDRequest{Id: txctx.Envelope.GetMetadata().GetId()})
			assert.Equal(t, test.expectedStatus, e.StatusInfo.Status, "Expected same status")
		})
	}

}
