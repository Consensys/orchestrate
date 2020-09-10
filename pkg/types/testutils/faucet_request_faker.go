package testutils

import (
	"math/big"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/chain-registry"
)

// TestRequest useful data to test CreditFunc
type TestRequest struct {
	Req                          *chainregistry.Request
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

func FakeFaucet() *chainregistry.Faucet {
	return &chainregistry.Faucet{
		UUID:       uuid.Must(uuid.NewV4()).String(),
		Amount:     big.NewInt(10),
		MaxBalance: big.NewInt(10),
		Creditor:   ethcommon.HexToAddress("0xacd"),
	}
}
