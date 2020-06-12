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
		From:            "0x7357589f8e367c2C31F51242fB77B350A11830F3",
		To:              "0x7357589f8e367c2C31F51242fB77B350A11830F2",
		Value:           "1",
		GasPrice:        "0",
		GasLimit:        "0",
		MethodSignature: "method(string,string)",
		Args:            []string{"val1", "val2"},
		ContractName:    "ContractName",
		ContractTag:     "ContractTag",
		Nonce:           "1",
	}
}

func FakeTesseraTransactionParams() *types.ETHTransactionParams {
	tx := FakeETHTransactionParams()
	tx.PrivateFrom = "ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc="
	tx.PrivateFor = []string{"ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bd="}
	tx.Protocol = types.TesseraPrivateTransaction

	return tx
}

func FakeOrionTransactionParams() *types.ETHTransactionParams {
	tx := FakeETHTransactionParams()
	tx.PrivateFrom = "ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Be="
	tx.PrivacyGroupID = "ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bf="
	tx.Protocol = types.OrionEEATransaction

	return tx
}

func FakeRawTransactionParams() *types.ETHTransactionParams {
	tx := FakeETHTransactionParams()
	tx.PrivateFrom = "ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Ba="
	tx.Raw = "0xABCDE012312312"
	tx.Protocol = types.OrionEEATransaction

	return tx
}
