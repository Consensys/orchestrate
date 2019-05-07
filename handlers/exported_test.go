package handlers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-signer.git/handlers/signer"
)

// Init inialize handlers
func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, signer.GlobalHandler(), "Global signer should have been set")
}
