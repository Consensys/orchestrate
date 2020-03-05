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
	ExpectedErr                  error
	ResultErr                    error
}

// AssertRequest make sure that a TestRequest is matching expected result
func AssertRequest(t *testing.T, test *TestRequest) {
	assert.Equal(t, test.ExpectedAmount, test.ResultAmount, "Amount credited should be correct expecting %s, got %s", test.ResultAmount, test.ExpectedAmount)
	if test.ExpectedErr != nil {
		assert.Equal(t, test.ExpectedErr, test.ResultErr, "Credit should error")
	} else {
		assert.NoError(t, test.ResultErr, "Credit should not error")
	}
}
