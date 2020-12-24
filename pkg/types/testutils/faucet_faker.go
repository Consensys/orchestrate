package testutils

import (
	"github.com/gofrs/uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
)

func FakeFaucet() *entities.Faucet {
	return &entities.Faucet{
		UUID:            uuid.Must(uuid.NewV4()).String(),
		Name:            "faucet-mainnet",
		ChainRule:       uuid.Must(uuid.NewV4()).String(),
		CreditorAccount: "0x5Cc634233E4a454d47aACd9fC68801482Fb02610",
		MaxBalance:      "6000000",
		Amount:          "100",
		Cooldown:        "10s",
		TenantID:        "_",
	}
}
