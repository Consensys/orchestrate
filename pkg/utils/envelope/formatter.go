package envelope

import (
	"math/big"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/tx"
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
				From:           job.Transaction.From,
				To:             job.Transaction.To,
				Gas:            job.Transaction.Gas,
				GasPrice:       job.Transaction.GasPrice,
				Value:          job.Transaction.Value,
				Nonce:          job.Transaction.Nonce,
				Data:           job.Transaction.Data,
				Raw:            job.Transaction.Raw,
				PrivateFor:     job.Transaction.PrivateFor,
				PrivateFrom:    job.Transaction.PrivateFrom,
				PrivacyGroupId: job.Transaction.PrivacyGroupID,
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

	if job.InternalData.ExpectedNonce != "" {
		txEnvelope.SetExpectedNonce(job.InternalData.ExpectedNonce)
	}

	if job.InternalData.ParentJobUUID != "" {
		txEnvelope.SetParentJobUUID(job.InternalData.ParentJobUUID)
	}

	return txEnvelope
}

func NewJobFromEnvelope(envelope *tx.Envelope, tenantID string) *entities.Job {
	return &entities.Job{
		UUID:         envelope.GetJobUUID(),
		NextJobUUID:  envelope.GetNextJobUUID(),
		ChainUUID:    envelope.GetChainUUID(),
		ScheduleUUID: envelope.GetScheduleUUID(),
		Type:         envelope.GetJobTypeString(),
		InternalData: &entities.InternalData{
			OneTimeKey:    envelope.IsOneTimeKeySignature(),
			ChainID:       envelope.GetChainIDString(),
			ParentJobUUID: envelope.GetParentJobUUID(),
			ExpectedNonce: envelope.GetExpectedNonce(),
		},
		TenantID: tenantID,
		Transaction: &entities.ETHTransaction{
			Hash:           envelope.GetTxHashString(),
			From:           envelope.GetFromString(),
			To:             envelope.GetToString(),
			Nonce:          envelope.GetNonceString(),
			Value:          envelope.GetValueString(),
			GasPrice:       envelope.GetGasPriceString(),
			Gas:            envelope.GetGasString(),
			Data:           envelope.GetData(),
			Raw:            envelope.GetRaw(),
			PrivateFrom:    envelope.GetPrivateFrom(),
			PrivateFor:     envelope.GetPrivateFor(),
			PrivacyGroupID: envelope.GetPrivacyGroupID(),
			EnclaveKey:     envelope.GetEnclaveKey(),
		},
	}
}
