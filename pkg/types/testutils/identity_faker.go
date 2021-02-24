package testutils

import (
	"time"

	"github.com/ConsenSys/orchestrate/pkg/multitenancy"

	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/pkg/utils"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func FakeAccount() *entities.Account {
	return &entities.Account{
		Alias:               "MyAccount",
		TenantID:            multitenancy.DefaultTenant,
		Attributes:          make(map[string]string),
		Address:             "0x5Cc634233E4a454d47aACd9fC68801482Fb02610",
		PublicKey:           ethcommon.HexToHash(utils.RandHexString(12)).String(),
		CompressedPublicKey: ethcommon.HexToHash(utils.RandHexString(12)).String(),
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
}
