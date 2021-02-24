package testutils

import (
	"github.com/ConsenSys/orchestrate/pkg/types/api"
	"github.com/gofrs/uuid"
)

func FakeRegisterFaucetRequest() *api.RegisterFaucetRequest {
	return &api.RegisterFaucetRequest{
		Name:            "faucet-mainnet",
		ChainRule:       uuid.Must(uuid.NewV4()).String(),
		CreditorAccount: "0x6230592812dE2E256D1512504c3E8A3C49975f07",
		MaxBalance:      "60000",
		Amount:          "100",
		Cooldown:        "10s",
	}
}

func FakeUpdateFaucetRequest() *api.UpdateFaucetRequest {
	return &api.UpdateFaucetRequest{
		Name:            "faucet-mainnet",
		ChainRule:       uuid.Must(uuid.NewV4()).String(),
		CreditorAccount: "0x6230592812dE2E256D1512504c3E8A3C49975f07",
		MaxBalance:      "6000",
		Amount:          "100",
		Cooldown:        "10s",
	}
}
