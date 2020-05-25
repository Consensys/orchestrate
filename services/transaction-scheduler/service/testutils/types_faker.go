package testutils

import (
	uuid "github.com/satori/go.uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
)

func FakeSendTransactionRequest() *types.SendTransactionRequest {
	return &types.SendTransactionRequest{
		BaseTransactionRequest: types.BaseTransactionRequest{
			IdempotencyKey: uuid.NewV4().String(),
		},
		Params: types.TransactionParams{
			From:            "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4",
			MethodSignature: "transfer()",
			To:              "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
		},
	}
}

func FakeTransactionResponse() *types.TransactionResponse {
	return &types.TransactionResponse{
		IdempotencyKey: uuid.NewV4().String(),
		Schedule:       FakeScheduleResponse(),
	}
}

func FakeCreateScheduleRequest() *types.CreateScheduleRequest {
	return &types.CreateScheduleRequest{
		ChainUUID: uuid.NewV4().String(),
	}
}

func FakeScheduleResponse() *types.ScheduleResponse {
	return &types.ScheduleResponse{
		UUID:      uuid.NewV4().String(),
		ChainUUID: uuid.NewV4().String(),
		Jobs:      []*types.JobResponse{FakeJobResponse()},
	}
}

func FakeCreateJobRequest() *types.CreateJobRequest {
	return &types.CreateJobRequest{
		ScheduleUUID: uuid.NewV4().String(),
		Type:         entities.JobConstantinopleTransaction,
		Labels:       nil,
		Transaction:  *FakeETHTransaction(),
	}
}

func FakeJobUpdateRequest() *types.UpdateJobRequest {
	return &types.UpdateJobRequest{
		Labels:      nil,
		Transaction: *FakeETHTransaction(),
	}
}

func FakeJobResponse() *types.JobResponse {
	return &types.JobResponse{
		UUID:        uuid.NewV4().String(),
		Transaction: *FakeETHTransaction(),
		Status:      entities.JobStatusCreated,
	}
}

func FakeETHTransaction() *types.ETHTransaction {
	return &types.ETHTransaction{
		Hash: "0xhash",
		From: "0xfrom",
		To:   "0xto",
		Data: "0xdede",
	}
}
