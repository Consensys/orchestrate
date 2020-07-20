package testutils

import (
	"github.com/gofrs/uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

func FakeSendTransactionRequest() *types.SendTransactionRequest {
	return &types.SendTransactionRequest{
		BaseTransactionRequest: types.BaseTransactionRequest{
			ChainName: "chainName",
		},
		Params: types.TransactionParams{
			From:            "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4",
			MethodSignature: "transfer()",
			To:              "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
		},
	}
}

func FakeSendRawTransactionRequest() *types.RawTransactionRequest {
	return &types.RawTransactionRequest{
		BaseTransactionRequest: types.BaseTransactionRequest{
			ChainName: "chainName",
		},
		Params: types.RawTransactionParams{
			Raw: "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4",
		},
	}
}

func FakeSendTransferTransactionRequest() *types.TransferRequest {
	return &types.TransferRequest{
		BaseTransactionRequest: types.BaseTransactionRequest{
			ChainName: "chainName",
		},
		Params: types.TransferParams{
			From:  "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4",
			Value: "1000000000000000000",
			To:    "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
		},
	}
}

func FakeDeployContractRequest() *types.DeployContractRequest {
	return &types.DeployContractRequest{
		BaseTransactionRequest: types.BaseTransactionRequest{
			ChainName: "chainName",
		},
		Params: types.DeployContractParams{
			From:         "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4",
			ContractName: "MyContract",
			ContractTag:  "v1.0.0",
		},
	}
}

func FakeSendTesseraRequest() *types.SendTransactionRequest {
	return &types.SendTransactionRequest{
		BaseTransactionRequest: types.BaseTransactionRequest{
			ChainName: "chainName",
		},
		Params: types.TransactionParams{
			From:            "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4",
			MethodSignature: "transfer()",
			To:              "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
			PrivateTransactionParams: types.PrivateTransactionParams{
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
			ChainName: "chainName",
		},
		Params: types.TransactionParams{
			From:            "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4",
			MethodSignature: "transfer()",
			To:              "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
			PrivateTransactionParams: types.PrivateTransactionParams{
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
		ScheduleUUID: uuid.Must(uuid.NewV4()).String(),
		ChainUUID:    uuid.Must(uuid.NewV4()).String(),
		Type:         utils.EthereumTransaction,
		Labels:       nil,
		Annotations:  nil,
		Transaction:  FakeETHTransaction(),
	}
}

func FakeJobUpdateRequest() *types.UpdateJobRequest {
	return &types.UpdateJobRequest{
		Labels:      nil,
		Transaction: FakeETHTransaction(),
		Status:      utils.StatusPending,
	}
}

func FakeJobResponse() *types.JobResponse {
	return &types.JobResponse{
		UUID:        uuid.Must(uuid.NewV4()).String(),
		Transaction: FakeETHTransaction(),
		Status:      utils.StatusCreated,
	}
}
