// +build unit

package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/proto"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, GlobalClient(), "Global should have been set")

	var c svc.ContractRegistryClient
	SetGlobalClient(c)
	assert.Nil(t, GlobalClient(), "Global should be reset to nil")
}
