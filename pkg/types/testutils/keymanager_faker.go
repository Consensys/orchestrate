package testutils

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
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
		Data:      "0x",
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
		Data:      "0x",
	}
}

func FakeSignEEATransactionRequest() *types.SignEEATransactionRequest {
	return &types.SignEEATransactionRequest{
		Namespace:   "_",
		Nonce:       0,
		To:          "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
		ChainID:     "1",
		PrivateFrom: "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=",
		PrivateFor:  []string{"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="},
		Data:        "0x",
	}
}

func FakeETHAccountResponse() *types.ETHAccountResponse {
	return &types.ETHAccountResponse{
		Namespace: "_",
		Address:   ethcommon.HexToAddress("0x" + utils.RandHexString(12)).String(),
		PublicKey: ethcommon.HexToHash(utils.RandHexString(12)).String(),
	}
}

func FakeSignTypedDataRequest() *types.SignTypedDataRequest {
	return &types.SignTypedDataRequest{
		Namespace: "_",
		DomainSeparator: types.DomainSeparator{
			Name:              "orchestrate",
			Version:           "v2.6.0",
			ChainID:           1,
			VerifyingContract: "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
			Salt:              "mySalt",
		},
		Types: map[string][]types.Type{
			"Mail": {
				{Name: "sender", Type: "address"},
				{Name: "recipient", Type: "address"},
				{Name: "content", Type: "string"},
			},
		},
		Message: map[string]interface{}{
			"sender":    "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
			"recipient": "0xFE3B557E8Fb62b89F4916B721be55cEb828dBd73",
			"content":   "my content",
		},
		MessageType: "Mail",
	}
}
