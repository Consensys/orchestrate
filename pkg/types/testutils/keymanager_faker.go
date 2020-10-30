package testutils

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/keymanager/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

func FakeCreateETHAccountRequest() *types.CreateETHAccountRequest {
	return &types.CreateETHAccountRequest{
		Namespace: "_",
	}
}

func FakeImportETHAccountRequest() *types.ImportETHAccountRequest {
	return &types.ImportETHAccountRequest{
		Namespace:  "_",
		PrivateKey: "fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249",
	}
}

func FakeSignETHTransactionRequest() *types.SignETHTransactionRequest {
	return &types.SignETHTransactionRequest{
		Namespace: "_",
		Nonce:     0,
		To:        "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
		Amount:    "10000000000",
		GasPrice:  "10000000000",
		GasLimit:  21000,
		ChainID:   "1",
	}
}

func FakeSignQuorumPrivateTransactionRequest() *types.SignQuorumPrivateTransactionRequest {
	return &types.SignQuorumPrivateTransactionRequest{
		Namespace: "_",
		Nonce:     0,
		To:        "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
		Amount:    "10000000000",
		GasPrice:  "10000000000",
		GasLimit:  21000,
	}
}

func FakeSignEEATransactionRequest() *types.SignEEATransactionRequest {
	return &types.SignEEATransactionRequest{
		Namespace:   "_",
		Nonce:       0,
		To:          "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
		GasPrice:    "10000000000",
		GasLimit:    21000,
		ChainID:     "1",
		PrivateFrom: "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=",
		PrivateFor:  []string{"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="},
	}
}

func FakeETHAccountResponse() *types.ETHAccountResponse {
	return &types.ETHAccountResponse{
		Namespace: "_",
		Address:   ethcommon.HexToAddress("0x" + utils.RandHexString(12)).String(),
		PublicKey: ethcommon.HexToHash(utils.RandHexString(12)).String(),
	}
}
