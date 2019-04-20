package testutils

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git/types"
)

// TestCreditData useful data to test CreditFunc
type TestCreditData struct {
	Req            *types.Request
	ResultAmount   *big.Int
	ResultOK       bool
	ResultErr      error
	ExpectedAmount *big.Int
	ExpectedOK     bool
	ExpectedErr    error
}

// AssertCreditData make sure that a TestData is matching expected result
func AssertCreditData(t *testing.T, data *TestCreditData) {
	assert.Equal(t, data.ResultOK, data.ExpectedOK, "Credit status incorrect")
	assert.Equal(t, 0, data.ResultAmount.Cmp(data.ExpectedAmount), "Amound credited should be correct")
	if data.ExpectedErr == nil {
		assert.Nil(t, data.ResultErr, "Credit should not error")
	} else {
		assert.NotNil(t, data.ResultErr, "Credit should error")
	}
}
