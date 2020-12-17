package testutils

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gofrs/uuid"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/chainregistry"
)

// TestRequest useful data to test CreditFunc
type TestRequest struct {
	Req                          *chainregistry.Request
	ResultAmount, ExpectedAmount *big.Int
	ExpectedErr                  error
	ResultErr                    error
}

func FakeFaucet() *chainregistry.Faucet {
	return &chainregistry.Faucet{
		UUID:       uuid.Must(uuid.NewV4()).String(),
		Amount:     big.NewInt(10),
		MaxBalance: big.NewInt(10),
		Creditor:   ethcommon.HexToAddress("0x12278c8C089ef98b4045f0b649b61Ed4316B1a50"),
	}
}
