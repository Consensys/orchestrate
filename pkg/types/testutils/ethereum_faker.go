package testutils

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
)

func FakeETHTransaction() *types.ETHTransaction {
	return &types.ETHTransaction{
		From:           "From",
		To:             "To",
		Nonce:          "Nonce",
		Value:          "Value",
		GasPrice:       "GasPrice",
		GasLimit:       "GasLimit",
		Data:           "Data",
		Raw:            "Raw",
		PrivateFrom:    "PrivateFrom",
		PrivateFor:     []string{"val1", "val2"},
		PrivacyGroupID: "PrivacyGroupID",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

func FakeETHTransactionParams() *types.ETHTransactionParams {
	return &types.ETHTransactionParams{
		From:            "From",
		To:              "To",
		Value:           "Value",
		GasPrice:        "GasPrice",
		GasLimit:        "GasLimit",
		MethodSignature: "constructor(string,string)",
		Args:            []string{"val1", "val2"},
		Raw:             "Raw",
		ContractName:    "ContractName",
		ContractTag:     "ContractTag",
		Nonce:           "1",
		PrivateFrom:     "PrivateFrom",
		PrivateFor:      []string{"val1", "val2"},
		PrivacyGroupID:  "PrivacyGroupID",
	}
}
