package handlers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-sender.git/handlers/sender"
)

// Init inialize handlers
func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, sender.GlobalHandler(), "Global store should have been set")
}
