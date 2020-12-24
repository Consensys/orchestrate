package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
)

func TestFaucetsParser(t *testing.T) {
	faucet := testutils.FakeFaucet()
	faucetModel := NewFaucetModelFromEntity(faucet)
	finalFaucet := NewFaucetFromModel(faucetModel)

	assert.Equal(t, faucet, finalFaucet)
}
