package testutils

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/faucet.git/types"
)

// TestRequest useful data to test CreditFunc
type TestRequest struct {
	Req                          *types.Request
	ResultAmount, ExpectedAmount *big.Int
	ResultOK, ExpectedOK         bool
	ResultErr, ExpectedErr       error
}

// AssertRequest make sure that a TestRequest is matching expected result
func AssertRequest(t *testing.T, test *TestRequest) {
	assert.Equal(t, test.ResultOK, test.ExpectedOK, "Credit status incorrect")
	assert.Equal(t, 0, test.ResultAmount.Cmp(test.ExpectedAmount), "Amound credited should be correct")
	if test.ExpectedErr == nil {
		assert.Nil(t, test.ResultErr, "Credit should not error")
	} else {
		assert.NotNil(t, test.ResultErr, "Credit should error")
	}
}
