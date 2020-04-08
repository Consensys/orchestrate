// +build unit

package storer

import (
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	clientmock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/client/mock"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
)

func TestRawTxStore(t *testing.T) {
	testSet := []struct {
		name           string
		input          func(txctx *engine.TxContext) *engine.TxContext
		expectedStatus svc.Status
	}{
		{
			"Store",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.SetChainID(big.NewInt(1)).SetID("test")
				return txctx
			},
			svc.Status_PENDING,
		},
		{
			"Store envelope without Metadata UUID",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.SetChainID(big.NewInt(1)).SetID("test")
				err := errors.InternalError("error").ExtendComponent(component)
				txctx.Envelope.Errors = append(txctx.Envelope.Errors, err)
				return txctx
			},
			svc.Status_ERROR,
		},
	}

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			txctx := engine.NewTxContext()
			txctx.Logger = log.NewEntry(log.New())
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			registry := clientmock.NewMockEnvelopeStoreClient(mockCtrl)
			registry.EXPECT().Store(gomock.Any(), gomock.AssignableToTypeOf(&svc.StoreRequest{}))
			registry.EXPECT().SetStatus(gomock.Any(), &svc.SetStatusRequest{
				Id:     "test",
				Status: test.expectedStatus,
			})
			registry.EXPECT().LoadByID(gomock.Any(), &svc.LoadByIDRequest{
				Id: "test",
			}).Return(&svc.StoreResponse{
				StatusInfo: &svc.StatusInfo{Status: test.expectedStatus},
			}, nil)

			h := RawTxStore(registry)
			h(test.input(txctx))
			e, _ := registry.LoadByID(txctx.Context(), &svc.LoadByIDRequest{Id: txctx.Envelope.GetID()})
			assert.Equal(t, test.expectedStatus, e.StatusInfo.Status, "Expected same status")
		})
	}

}
