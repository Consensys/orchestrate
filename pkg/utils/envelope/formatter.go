package envelope

import (
	"context"
	"math/big"

	authutils "github.com/consensys/orchestrate/pkg/toolkit/app/auth/utils"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/types/ethereum"
	"github.com/consensys/orchestrate/pkg/types/tx"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func NewEnvelopeFromJob(job *entities.Job, headers map[string]string) *tx.TxEnvelope {
	contextLabels := job.Labels
	if contextLabels == nil {
		contextLabels = map[string]string{}
	}

	contextLabels[tx.NextJobUUIDLabel] = job.NextJobUUID
	contextLabels[tx.PriorityLabel] = job.InternalData.Priority
	contextLabels[tx.ParentJobUUIDLabel] = job.InternalData.ParentJobUUID

	txEnvelope := &tx.TxEnvelope{
		Msg: &tx.TxEnvelope_TxRequest{TxRequest: &tx.TxRequest{
			Id:      job.ScheduleUUID,
			Headers: headers,
			Params: &tx.Params{
				From:            job.Transaction.From,
				To:              job.Transaction.To,
				Gas:             job.Transaction.Gas,
				GasPrice:        job.Transaction.GasPrice,
				GasFeeCap:       job.Transaction.GasFeeCap,
				GasTipCap:       job.Transaction.GasTipCap,
				Value:           job.Transaction.Value,
				Nonce:           job.Transaction.Nonce,
				Data:            job.Transaction.Data,
				Raw:             job.Transaction.Raw,
				PrivateFrom:     job.Transaction.PrivateFrom,
				PrivateFor:      job.Transaction.PrivateFor,
				MandatoryFor:    job.Transaction.MandatoryFor,
				PrivacyGroupId:  job.Transaction.PrivacyGroupID,
				PrivacyFlag:     int32(job.Transaction.PrivacyFlag),
				TransactionType: string(job.Transaction.TransactionType),
				AccessList:      ConvertFromAccessList(job.Transaction.AccessList),
			},
			ContextLabels: contextLabels,
			JobType:       tx.JobTypeMap[job.Type],
		}},
		InternalLabels: make(map[string]string),
	}

	txEnvelope.SetChainUUID(job.ChainUUID)

	chainID := new(big.Int)
	chainID.SetString(job.InternalData.ChainID, 10)
	txEnvelope.SetChainID(chainID)
	txEnvelope.SetScheduleUUID(job.ScheduleUUID)
	txEnvelope.SetJobUUID(job.UUID)

	if job.InternalData.OneTimeKey {
		txEnvelope.EnableTxFromOneTimeKey()
	}

	if job.InternalData.ParentJobUUID != "" {
		txEnvelope.SetParentJobUUID(job.InternalData.ParentJobUUID)
	}

	if job.InternalData.Priority != "" {
		txEnvelope.SetPriority(job.InternalData.Priority)
	}

	if job.Transaction.Hash != "" {
		txEnvelope.SetTxHash(job.Transaction.Hash)
	}

	return txEnvelope
}

func NewContextFromEnvelope(ctx context.Context, envelope *tx.Envelope) context.Context {
	return multitenancy.WithUserInfo(ctx, multitenancy.NewUserInfo(
		envelope.GetHeadersValue(authutils.TenantIDHeader),
		envelope.GetHeadersValue(authutils.UsernameHeader),
	))
}

func NewJobFromEnvelope(envelope *tx.Envelope) *entities.Job {
	return &entities.Job{
		UUID:         envelope.GetJobUUID(),
		NextJobUUID:  envelope.GetNextJobUUID(),
		ChainUUID:    envelope.GetChainUUID(),
		ScheduleUUID: envelope.GetScheduleUUID(),
		Type:         entities.JobType(envelope.GetJobTypeString()),
		InternalData: &entities.InternalData{
			OneTimeKey:    envelope.IsOneTimeKeySignature(),
			ChainID:       envelope.GetChainIDString(),
			ParentJobUUID: envelope.GetParentJobUUID(),
			Priority:      envelope.GetPriority(),
		},
		TenantID: envelope.GetHeadersValue(authutils.TenantIDHeader),
		OwnerID:  envelope.GetHeadersValue(authutils.UsernameHeader),
		Transaction: &entities.ETHTransaction{
			Hash:            envelope.GetTxHashString(),
			From:            envelope.GetFromString(),
			To:              envelope.GetToString(),
			Nonce:           envelope.GetNonceString(),
			Value:           envelope.GetValueString(),
			GasPrice:        envelope.GetGasPriceString(),
			Gas:             envelope.GetGasString(),
			GasFeeCap:       envelope.GetGasFeeCapString(),
			GasTipCap:       envelope.GetGasTipCapString(),
			AccessList:      ConvertToAccessList(envelope.GetAccessList()),
			TransactionType: entities.TransactionType(envelope.GetTransactionType()),
			Data:            envelope.GetData(),
			Raw:             envelope.GetRaw(),
			PrivateFrom:     envelope.GetPrivateFrom(),
			PrivateFor:      envelope.GetPrivateFor(),
			MandatoryFor:    envelope.GetMandatoryFor(),
			PrivacyGroupID:  envelope.GetPrivacyGroupID(),
			PrivacyFlag:     envelope.GetPrivacyFlag(),
			EnclaveKey:      envelope.GetEnclaveKey(),
		},
	}
}

func ConvertFromAccessList(accessList types.AccessList) []*ethereum.AccessTuple {
	result := []*ethereum.AccessTuple{}
	for _, t := range accessList {
		tupl := &ethereum.AccessTuple{
			Address:     t.Address.Hex(),
			StorageKeys: []string{},
		}

		for _, k := range t.StorageKeys {
			tupl.StorageKeys = append(tupl.StorageKeys, k.Hex())
		}

		result = append(result, tupl)
	}

	return result
}

func ConvertToAccessList(accessList []*ethereum.AccessTuple) types.AccessList {
	result := types.AccessList{}
	for _, item := range accessList {
		storageKeys := []ethcommon.Hash{}
		for _, sk := range item.StorageKeys {
			storageKeys = append(storageKeys, ethcommon.HexToHash(sk))
		}

		result = append(result, types.AccessTuple{
			Address:     ethcommon.HexToAddress(item.Address),
			StorageKeys: storageKeys,
		})
	}

	return result
}
