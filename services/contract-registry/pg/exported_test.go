package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init()
	assert.NotNil(t, GlobalContractRegistry(), "Global should have been set")

	var contractRegistry *ContractRegistry
	SetGlobalContractRegistry(contractRegistry)
	assert.Nil(t, GlobalContractRegistry(), "Global should be reset to nil")
}
