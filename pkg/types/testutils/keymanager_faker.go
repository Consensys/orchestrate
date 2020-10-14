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

func FakeImportETHAccountRequest() *types.ImportETHAccountRequest {
	return &types.ImportETHAccountRequest{
		KeyType:    utils.Secp256k1,
		Namespace:  "_",
		PrivateKey: "fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249",
	}
}

func FakeETHAccountResponse() *types.ETHAccountResponse {
	return &types.ETHAccountResponse{
		KeyType:   utils.Secp256k1,
		Namespace: "_",
		Address:   ethcommon.HexToAddress("0x" + utils.RandHexString(12)).String(),
		PublicKey: ethcommon.HexToHash(utils.RandHexString(12)).String(),
	}
}
