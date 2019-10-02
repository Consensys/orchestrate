package testutils

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/types"
)

// TestRequest useful data to test CreditFunc
type TestRequest struct {
	Req                               *types.Request
	ResultAmount, ExpectedAmount      *big.Int
	ResultOK, ExpectedOK, ExpectedErr bool
	ResultErr                         error
}

// AssertRequest make sure that a TestRequest is matching expected result
func AssertRequest(t *testing.T, test *TestRequest) {
	assert.Equal(t, test.ResultOK, test.ExpectedOK, "Credit status incorrect")
	assert.Equal(t, 0, test.ResultAmount.Cmp(test.ExpectedAmount), "Amount credited should be correct")
	if test.ExpectedErr {
		assert.NotNil(t, test.ResultErr, "Credit should error")
	} else {
		assert.Nil(t, test.ResultErr, "Credit should not error")
	}
}
