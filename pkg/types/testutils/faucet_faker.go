package testutils

import (
	"math/big"

	"github.com/consensys/orchestrate/pkg/types/entities"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gofrs/uuid"
)

func FakeFaucet() *entities.Faucet {
	return &entities.Faucet{
		UUID:            uuid.Must(uuid.NewV4()).String(),
		Name:            "faucet-mainnet",
		ChainRule:       uuid.Must(uuid.NewV4()).String(),
		CreditorAccount: ethcommon.HexToAddress("0x5Cc634233E4a454d47aACd9fC68801482Fb02610"),
		MaxBalance:      hexutil.Big(*big.NewInt(6000000)),
		Amount:          hexutil.Big(*big.NewInt(100)),
		Cooldown:        "10s",
		TenantID:        "_",
	}
}
