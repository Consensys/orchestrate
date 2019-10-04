package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	evlpstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope-store"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, GlobalEnvelopeStoreClient(), "Global should have been set")

	var c evlpstore.EnvelopeStoreClient
	SetGlobalEnvelopeStoreClient(c)
	assert.Nil(t, GlobalEnvelopeStoreClient(), "Global should be reset to nil")
}
