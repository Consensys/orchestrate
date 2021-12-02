package testutils

import (
	"math/big"

	"github.com/consensys/orchestrate/pkg/types/api"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gofrs/uuid"
)

func FakeRegisterFaucetRequest() *api.RegisterFaucetRequest {
	return &api.RegisterFaucetRequest{
		Name:            "faucet-mainnet",
		ChainRule:       uuid.Must(uuid.NewV4()).String(),
		CreditorAccount: ethcommon.HexToAddress("0x6230592812dE2E256D1512504c3E8A3C49975f07"),
		MaxBalance:      hexutil.Big(*big.NewInt(6000000)),
		Amount:          hexutil.Big(*big.NewInt(100)),
		Cooldown:        "10s",
	}
}

func FakeUpdateFaucetRequest() *api.UpdateFaucetRequest {
	return &api.UpdateFaucetRequest{
		Name:            "faucet-mainnet",
		ChainRule:       uuid.Must(uuid.NewV4()).String(),
		CreditorAccount: ethcommon.HexToAddress("0x6230592812dE2E256D1512504c3E8A3C49975f07"),
		MaxBalance:      hexutil.Big(*big.NewInt(6000000)),
		Amount:          hexutil.Big(*big.NewInt(100)),
		Cooldown:        "10s",
	}
}
