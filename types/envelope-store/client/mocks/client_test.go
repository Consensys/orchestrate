package mocks

import (
	"context"
	"math/big"
	"testing"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/tx"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	evlpstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope-store"
)

func TestEnvelopeStoreClient(t *testing.T) {
	client := New()
	var c interface{} = client
	_, ok := c.(evlpstore.EnvelopeStoreClient)
	assert.True(t, ok, "Should match EnvelopeStoreClient interface")

	_, err := client.Store(
		context.Background(),
		&evlpstore.StoreRequest{
			Envelope: tx.NewEnvelope().
				SetChainIDUint64(888).
				SetID("a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11").
				MustSetTxHashString("0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210").
				TxEnvelopeAsRequest(),
		},
	)
	assert.Nil(t, err, "Store should not error")

	_, err = client.LoadByID(
		context.Background(),
		&evlpstore.LoadByIDRequest{
			Id: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
		},
	)
	assert.Nil(t, err, "LoadByID should not error")

	_, err = client.LoadByTxHash(
		context.Background(),
		&evlpstore.LoadByTxHashRequest{
			ChainId: big.NewInt(888).String(),
			TxHash:  "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210",
		},
	)
	assert.Nil(t, err, "LoadByTxHash should not error")

	_, err = client.SetStatus(
		context.Background(),
		&evlpstore.SetStatusRequest{
			Id:     "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
			Status: evlpstore.Status_PENDING,
		},
	)
	assert.Nil(t, err, "SetStatus should not error")

	_, err = client.LoadPending(
		context.Background(),
		&evlpstore.LoadPendingRequest{
			Duration: utils.DurationToPDuration(300 * time.Millisecond),
		},
	)
	assert.Nil(t, err, "LoadPending should not error")
}
