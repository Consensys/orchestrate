package testutils

import (
	types "github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/utils"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gofrs/uuid"
)

var FromAddress = ethcommon.HexToAddress("0x5Cc634233E4a454d47aACd9fC68801482Fb02610")

func FakeSendTransactionRequest() *types.SendTransactionRequest {
	return &types.SendTransactionRequest{
		ChainName: "ganache",
		Params: types.TransactionParams{
			From:            &FromAddress,
			MethodSignature: "transfer(address,uint256)",
			To:              FakeAddress(),
			ContractTag:     "contractTag",
			ContractName:    "contractName",
			Args:            []interface{}{"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18", 5},
		},
	}
}

func FakeSendRawTransactionRequest() *types.RawTransactionRequest {
	return &types.RawTransactionRequest{
		ChainName: "ganache",
		Params: types.RawTransactionParams{
			Raw: hexutil.MustDecode("0xf85380839896808252088083989680808216b4a0d35c752d3498e6f5ca1630d264802a992a141ca4b6a3f439d673c75e944e5fb0a05278aaa5fabbeac362c321b54e298dedae2d31471e432c26ea36a8d49cf08f1e"),
		},
	}
}

func FakeSendTransferTransactionRequest() *types.TransferRequest {
	return &types.TransferRequest{
		ChainName: "ganache",
		Params: types.TransferParams{
			From:  FromAddress,
			Value: (*hexutil.Big)(hexutil.MustDecodeBig("0xDE0B6B3A7640000")),
			To:    *FakeAddress(),
		},
	}
}

func FakeDeployContractRequest() *types.DeployContractRequest {
	return &types.DeployContractRequest{
		ChainName: "ganache",
		Params: types.DeployContractParams{
			From:         &FromAddress,
			ContractName: "MyContract",
			ContractTag:  "v1.0.0",
		},
	}
}

func FakeSendTesseraRequest() *types.SendTransactionRequest {
	return &types.SendTransactionRequest{
		ChainName: "ganache",
		Params: types.TransactionParams{
			From:            &FromAddress,
			MethodSignature: "transfer(address,uint256)",
			To:              utils.ToPtr(ethcommon.HexToAddress("0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18")).(*ethcommon.Address),
			Protocol:        entities.TesseraChainType,
			PrivateFrom:     "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=",
			PrivateFor:      []string{"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="},
			Args:            []interface{}{"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18", 5},
		},
	}
}

func FakeSendEEARequest() *types.SendTransactionRequest {
	return &types.SendTransactionRequest{
		ChainName: "ganache",
		Params: types.TransactionParams{
			From:            &FromAddress,
			MethodSignature: "transfer(address,uint256)",
			To:              utils.ToPtr(ethcommon.HexToAddress("0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18")).(*ethcommon.Address),
			Protocol:        entities.EEAChainType,
			PrivateFrom:     "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=",
			PrivacyGroupID:  "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=",
			Args:            []interface{}{"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18", 5},
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
