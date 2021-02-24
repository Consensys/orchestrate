package parsers

import (
	"testing"

	"github.com/ConsenSys/orchestrate/pkg/types/testutils"
	"github.com/stretchr/testify/assert"
)

func TestFaucetsParser(t *testing.T) {
	faucet := testutils.FakeFaucet()
	faucetModel := NewFaucetModelFromEntity(faucet)
	finalFaucet := NewFaucetFromModel(faucetModel)

	assert.Equal(t, faucet, finalFaucet)
}
