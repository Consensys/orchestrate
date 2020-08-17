package testutils

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

func ParseIArray(args ...interface{}) (ret []interface{}) {
	ret = make([]interface{}, len(args))
	copy(ret, args)
	return
}

func FakeETHTransaction() *entities.ETHTransaction {
	return &entities.ETHTransaction{
		From:           "0x1abae27a0cbfb02945720425d3b80c7e09728534",
		To:             "0x4fed1fc4144c223ae3c1553be203cdfcbd38c581",
		Nonce:          "1",
		Value:          "50000",
		GasPrice:       "10000",
		Gas:            "21000",
		Data:           "0x",
		Raw:            "0x",
		PrivateFrom:    "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=",
		PrivateFor:     []string{"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=", "B1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="},
		PrivacyGroupID: "PrivacyGroupID",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

func FakeETHTransactionParams() *entities.ETHTransactionParams {
	return &entities.ETHTransactionParams{
		From:            "0x7357589f8e367c2C31F51242fB77B350A11830F3",
		To:              "0x7357589f8e367c2C31F51242fB77B350A11830F2",
		Value:           "1",
		GasPrice:        "0",
		Gas:             "0",
		MethodSignature: "method(string,string)",
		Args:            ParseIArray("val1", "val2"),
		ContractName:    "ContractName",
		ContractTag:     "ContractTag",
		Nonce:           "1",
	}
}

func FakeTesseraTransactionParams() *entities.ETHTransactionParams {
	tx := FakeETHTransactionParams()
	tx.PrivateFrom = "ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc="
	tx.PrivateFor = []string{"ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bd="}
	tx.Protocol = utils.TesseraPrivateTransaction

	return tx
}

func FakeOrionTransactionParams() *entities.ETHTransactionParams {
	tx := FakeETHTransactionParams()
	tx.PrivateFrom = "ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Be="
	tx.PrivacyGroupID = "ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bf="
	tx.Protocol = utils.OrionEEATransaction

	return tx
}

func FakeRawTransactionParams() *entities.ETHTransactionParams {
	return &entities.ETHTransactionParams{
		Raw: "0xABCDE012312312",
	}
}

func FakeTransferTransactionParams() *entities.ETHTransactionParams {
	return &entities.ETHTransactionParams{
		From:  "0x7357589f8e367c2C31F51242fB77B350A11830FA",
		To:    "0x7357589f8e367c2C31F51242fB77B350A11830FB",
		Value: "10000000000",
	}
}
