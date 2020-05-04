package testutils

import (
	uuid "github.com/satori/go.uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
)

func FakeTransactionRequest() *types.TransactionRequest {
	return &types.TransactionRequest{
		BaseTransactionRequest: types.BaseTransactionRequest{
			IdempotencyKey: uuid.NewV4().String(),
			ChainID:        uuid.NewV4().String(),
		},
		Params: types.TransactionParams{
			From:            "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4",
			To:              "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
			MethodSignature: "constructor()",
		},
	}
}

func FakeScheduleRequest() *types.ScheduleRequest {
	return &types.ScheduleRequest{ChainID: uuid.NewV4().String()}
}

func FakeScheduleResponse() *types.ScheduleResponse {
	return &types.ScheduleResponse{
		UUID:    uuid.NewV4().String(),
		ChainID: uuid.NewV4().String(),
	}
}
