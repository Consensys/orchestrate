// +build unit

package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, GlobalEnvelopeStoreClient(), "Global should have been set")

	var c svc.EnvelopeStoreClient
	SetGlobalEnvelopeStoreClient(c)
	assert.Nil(t, GlobalEnvelopeStoreClient(), "Global should be reset to nil")
}
