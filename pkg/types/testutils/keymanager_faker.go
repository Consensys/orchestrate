package testutils

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/keymanager/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

func FakeCreateETHAccountRequest() *types.CreateETHAccountRequest {
	return &types.CreateETHAccountRequest{
		KeyType:   utils.Secp256k1,
		Namespace: "_",
	}
}

func FakeETHAccountResponse() *types.ETHAccountResponse {
	return &types.ETHAccountResponse{
		KeyType:   utils.Secp256k1,
		Namespace: "_",
		Address:   ethcommon.HexToAddress("0x" + utils.RandHexString(12)).String(),
		PublicKey: ethcommon.HexToHash("0x" + utils.RandHexString(12)).String(),
	}
}
