package quorumkeymanager

import (
	"strings"

	"github.com/ConsenSys/orchestrate/pkg/quorum-key-manager/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func FakeEth1AccountResponse(addr string, tenants []string) *types.Eth1AccountResponse {
	return &types.Eth1AccountResponse{
		Address: ethcommon.HexToAddress(addr),
		Tags: map[string]string{
			TagIDAllowedTenants: strings.Join(tenants, TagSeparatorAllowedTenants),
		},
	}
}

func FakeSignTypedDataRequest() *types.SignTypedDataRequest {
	return &types.SignTypedDataRequest{
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

func FakeVerifyPayloadRequest() *types.VerifyEth1SignatureRequest {
	return &types.VerifyEth1SignatureRequest{
		Data:      hexutil.MustDecode(hexutil.Encode([]byte("my data to sign"))),
		Signature: hexutil.MustDecode("0x34334af7bacf5d82bb892c838beda65331232c29e122b3485f31e14eda731dbb0ebae9d1eed72c099ff4c3b462aebf449068f717f3638a6facd0b3dddf2529a500"),
		Address:   ethcommon.HexToAddress("0x5Cc634233E4a454d47aACd9fC68801482Fb02610"),
	}
}

func FakeVerifyTypedDataPayloadRequest() *types.VerifyTypedDataRequest {
	return &types.VerifyTypedDataRequest{
		TypedData: *FakeSignTypedDataRequest(),
		Signature: hexutil.MustDecode("0x3399aeb23d6564b3a0b220447e9f1bb2057ffb82cfb766147620aa6bc84938e26941e7583d6460fea405d99da897e88cab07a7fd0991c6c2163645c45d25e4b201"),
		Address:   ethcommon.HexToAddress("0x5Cc634233E4a454d47aACd9fC68801482Fb02610"),
	}
}
