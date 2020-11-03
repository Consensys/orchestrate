package service

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
)

func EnvelopeToJob(envelope *tx.Envelope, tenantID string) *entities.Job {
	return &entities.Job{
		UUID:         envelope.GetJobUUID(),
		NextJobUUID:  envelope.GetNextJobUUID(),
		ChainUUID:    envelope.GetChainUUID(),
		ScheduleUUID: envelope.GetScheduleUUID(),
		Type:         envelope.GetJobTypeString(),
		InternalData: &entities.InternalData{
			OneTimeKey: envelope.IsOneTimeKeySignature(),
			ChainID:    envelope.GetChainIDString(),
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
