package testutils

import (
	"github.com/gofrs/uuid"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
)

const FromAddress = "0x5Cc634233E4a454d47aACd9fC68801482Fb02610"

func FakeSendTransactionRequest() *types.SendTransactionRequest {
	return &types.SendTransactionRequest{
		ChainName: "ganache",
		Params: types.TransactionParams{
			From:            FromAddress,
			MethodSignature: "transfer()",
			To:              "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
		},
	}
}

func FakeSendRawTransactionRequest() *types.RawTransactionRequest {
	return &types.RawTransactionRequest{
		ChainName: "ganache",
		Params: types.RawTransactionParams{
			Raw: "0xabeabe",
		},
	}
}

func FakeSendTransferTransactionRequest() *types.TransferRequest {
	return &types.TransferRequest{
		ChainName: "ganache",
		Params: types.TransferParams{
			From:  FromAddress,
			Value: "1000000000000000000",
			To:    "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
		},
	}
}

func FakeDeployContractRequest() *types.DeployContractRequest {
	return &types.DeployContractRequest{
		ChainName: "ganache",
		Params: types.DeployContractParams{
			From:         FromAddress,
			ContractName: "MyContract",
			ContractTag:  "v1.0.0",
		},
	}
}

func FakeSendTesseraRequest() *types.SendTransactionRequest {
	return &types.SendTransactionRequest{
		ChainName: "ganache",
		Params: types.TransactionParams{
			From:            FromAddress,
			MethodSignature: "transfer()",
			To:              "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
			Protocol:        entities.TesseraChainType,
			PrivateFrom:     "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=",
			PrivateFor:      []string{"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="},
		},
	}
}

func FakeSendOrionRequest() *types.SendTransactionRequest {
	return &types.SendTransactionRequest{
		ChainName: "ganache",
		Params: types.TransactionParams{
			From:            FromAddress,
			MethodSignature: "transfer()",
			To:              "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
			Protocol:        entities.OrionChainType,
			PrivateFrom:     "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=",
			PrivacyGroupID:  "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=",
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
		Type:         entities.EthereumTransaction,
		Transaction:  *FakeETHTransaction(),
	}
}

func FakeJobUpdateRequest() *types.UpdateJobRequest {
	return &types.UpdateJobRequest{
		Transaction: FakeETHTransaction(),
		Status:      entities.StatusPending,
	}
}

func FakeJobResponse() *types.JobResponse {
	return &types.JobResponse{
		UUID:        uuid.Must(uuid.NewV4()).String(),
		ChainUUID:   uuid.Must(uuid.NewV4()).String(),
		Transaction: *FakeETHTransaction(),
		Status:      entities.StatusCreated,
		Labels:      make(map[string]string),
		Annotations: types.Annotations{
			GasPricePolicy: types.GasPriceParams{
				RetryPolicy: types.RetryParams{
					Interval: "5s",
				},
			},
		},
	}
}
