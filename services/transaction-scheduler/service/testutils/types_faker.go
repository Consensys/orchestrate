package testutils

import (
	uuid "github.com/satori/go.uuid"
	types2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/types"
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

func FakeSendTesseraRequest() *types.SendTransactionRequest {
	return &types.SendTransactionRequest{
		BaseTransactionRequest: types.BaseTransactionRequest{
			IdempotencyKey: uuid.NewV4().String(),
		},
		Params: types.TransactionParams{
			From:            "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4",
			MethodSignature: "transfer()",
			To:              "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
			PrivateTransactionParams: types2.PrivateTransactionParams{
				Protocol:    utils.TesseraChainType,
				PrivateFrom: "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=",
				PrivateFor:  []string{"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="},
			},
		},
	}
}

func FakeSendOrionRequest() *types.SendTransactionRequest {
	return &types.SendTransactionRequest{
		BaseTransactionRequest: types.BaseTransactionRequest{
			IdempotencyKey: uuid.NewV4().String(),
		},
		Params: types.TransactionParams{
			From:            "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4",
			MethodSignature: "transfer()",
			To:              "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
			PrivateTransactionParams: types2.PrivateTransactionParams{
				Protocol:       utils.OrionChainType,
				PrivateFrom:    "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=",
				PrivacyGroupID: "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=",
			},
		},
	}
}

func FakeCreateScheduleRequest() *types.CreateScheduleRequest {
	return &types.CreateScheduleRequest{}
}

func FakeCreateJobRequest() *types.CreateJobRequest {
	return &types.CreateJobRequest{
		ScheduleUUID: uuid.NewV4().String(),
		ChainUUID:    uuid.NewV4().String(),
		Type:         types2.EthereumTransaction,
		Labels:       nil,
		Transaction:  testutils2.FakeETHTransaction(),
	}
}

func FakeJobUpdateRequest() *types.UpdateJobRequest {
	return &types.UpdateJobRequest{
		Labels:      nil,
		Transaction: testutils2.FakeETHTransaction(),
		Status:      types2.StatusPending,
	}
}

func FakeJobResponse() *types.JobResponse {
	return &types.JobResponse{
		UUID:        uuid.NewV4().String(),
		Transaction: testutils2.FakeETHTransaction(),
		Status:      types2.StatusCreated,
	}
}
