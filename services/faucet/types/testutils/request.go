package testutils

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types"
)

// TestRequest useful data to test CreditFunc
type TestRequest struct {
	Req                          *types.Request
	ResultAmount, ExpectedAmount *big.Int
	ExpectedErr                  bool
	ResultErr                    error
}

// AssertRequest make sure that a TestRequest is matching expected result
func AssertRequest(t *testing.T, test *TestRequest) {
	assert.True(t, test.ResultAmount.Cmp(test.ExpectedAmount) == 0, "Amount credited should be correct expecting %s, got %s, %s", test.ResultAmount, test.ExpectedAmount)
	if test.ExpectedErr {
		assert.NotNil(t, test.ResultErr, "Credit should error")
	} else {
		assert.Nil(t, test.ResultErr, "Credit should not error")
	}
}
