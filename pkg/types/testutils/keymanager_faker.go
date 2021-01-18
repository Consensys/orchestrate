package testutils

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethTypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager/ethereum"
	zksTypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager/zk-snarks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
)

func FakeCreateETHAccountRequest() *ethTypes.CreateETHAccountRequest {
	return &ethTypes.CreateETHAccountRequest{
		Namespace: "_",
	}
}

func FakeCreateZKSAccountRequest() *zksTypes.CreateZKSAccountRequest {
	return &zksTypes.CreateZKSAccountRequest{
		Namespace: "_",
	}
}

func FakeImportETHAccountRequest() *ethTypes.ImportETHAccountRequest {
	return &ethTypes.ImportETHAccountRequest{
		Namespace:  "_",
		PrivateKey: "fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249",
	}
}

func FakeSignETHTransactionRequest() *ethTypes.SignETHTransactionRequest {
	return &ethTypes.SignETHTransactionRequest{
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

func FakeSignQuorumPrivateTransactionRequest() *ethTypes.SignQuorumPrivateTransactionRequest {
	return &ethTypes.SignQuorumPrivateTransactionRequest{
		Namespace: "_",
		Nonce:     0,
		To:        "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
		Amount:    "10000000000",
		GasPrice:  "10000000000",
		GasLimit:  21000,
		Data:      "0x",
	}
}

func FakeSignEEATransactionRequest() *ethTypes.SignEEATransactionRequest {
	return &ethTypes.SignEEATransactionRequest{
		Namespace:   "_",
		Nonce:       0,
		To:          "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
		ChainID:     "1",
		PrivateFrom: "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=",
		PrivateFor:  []string{"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="},
		Data:        "0x",
	}
}

func FakeETHAccountResponse() *ethTypes.ETHAccountResponse {
	return &ethTypes.ETHAccountResponse{
		Namespace: "_",
		Address:   ethcommon.HexToAddress("0x" + utils.RandHexString(12)).String(),
		PublicKey: ethcommon.HexToHash(utils.RandHexString(12)).String(),
	}
}

func FakeSignTypedDataRequest() *ethTypes.SignTypedDataRequest {
	return &ethTypes.SignTypedDataRequest{
		Namespace: "_",
		DomainSeparator: ethTypes.DomainSeparator{
			Name:              "orchestrate",
			Version:           "v2.6.0",
			ChainID:           1,
			VerifyingContract: "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
			Salt:              "mySalt",
		},
		Types: map[string][]ethTypes.Type{
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

func FakeVerifyPayloadRequest() *ethTypes.VerifyPayloadRequest {
	return &ethTypes.VerifyPayloadRequest{
		Data:      "my data to sign",
		Signature: "0x34334af7bacf5d82bb892c838beda65331232c29e122b3485f31e14eda731dbb0ebae9d1eed72c099ff4c3b462aebf449068f717f3638a6facd0b3dddf2529a500",
		Address:   "0x5Cc634233E4a454d47aACd9fC68801482Fb02610",
	}
}

func FakeZKSVerifyPayloadRequest() *zksTypes.VerifyPayloadRequest {
	return &zksTypes.VerifyPayloadRequest{
		Data:      "my data to sign",
		Signature: "0x34334af7bacf5d82bb892c838beda65331232c29e122b3485f31e14eda731dbb0ebae9d1eed72c099ff4c3b462aebf449068f717f3638a6facd0b3dddf2529a500",
		PublicKey: "16551006344732991963827342392501535507890487822471009342749102663105305595515",
	}
}

func FakeVerifyTypedDataPayloadRequest() *ethTypes.VerifyTypedDataRequest {
	return &ethTypes.VerifyTypedDataRequest{
		TypedData: *FakeSignTypedDataRequest(),
		Signature: "0x3399aeb23d6564b3a0b220447e9f1bb2057ffb82cfb766147620aa6bc84938e26941e7583d6460fea405d99da897e88cab07a7fd0991c6c2163645c45d25e4b201",
		Address:   "0x5Cc634233E4a454d47aACd9fC68801482Fb02610",
	}
}
