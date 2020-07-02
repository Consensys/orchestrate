// +build unit

package parsers

import (
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
)

func TestParsersJob_NewModelFromEntity(t *testing.T) {
	jobEntity := testutils2.FakeJob()
	finalJobEntity := NewJobEntityFromModels(NewJobModelFromEntities(jobEntity, nil))

	expectedJSON, _ := json.Marshal(jobEntity)
	actualJSON, _ := json.Marshal(finalJobEntity)
	assert.Equal(t, string(expectedJSON), string(actualJSON))
}

func TestParsersJob_NewEntityFromModel(t *testing.T) {
	jobModel := testutils.FakeJobModel(1)
	jobEntity := NewJobEntityFromModels(jobModel)
	finalJobModel := NewJobModelFromEntities(jobEntity, jobModel.ScheduleID)
	finalJobModel.Schedule = jobModel.Schedule

	assert.Equal(t, finalJobModel.ScheduleID, jobModel.ScheduleID)
	assert.Equal(t, finalJobModel.UUID, jobModel.UUID)
	assert.Equal(t, finalJobModel.Type, jobModel.Type)
	assert.Equal(t, finalJobModel.Labels, jobModel.Labels)
	assert.Equal(t, finalJobModel.CreatedAt, jobModel.CreatedAt)
}

func TestParsersJob_NewEnvelopeFromModel(t *testing.T) {
	jobModel := testutils.FakeJobModel(1)
	headers := map[string]string{
		"Authorization": "Bearer MyToken",
	}
	txEnvelope := NewEnvelopeFromJobModel(jobModel, headers)

	txRequest := txEnvelope.GetTxRequest()
	assert.Equal(t, jobModel.ChainUUID, txEnvelope.GetChainUUID())
	assert.Equal(t, jobModel.UUID, txEnvelope.GetID())
	assert.Equal(t, tx.JobTypeMap[jobModel.Type], txRequest.GetJobType())
	assert.Equal(t, jobModel.Transaction.Sender, txRequest.Params.GetFrom())
	assert.Equal(t, jobModel.Transaction.Recipient, txRequest.Params.GetTo())
	assert.Equal(t, jobModel.Transaction.Data, txRequest.Params.GetData())
	assert.Equal(t, jobModel.Transaction.Nonce, txRequest.Params.GetNonce())
	assert.Equal(t, jobModel.Transaction.Raw, txRequest.Params.GetRaw())
	assert.Equal(t, jobModel.Transaction.GasPrice, txRequest.Params.GetGasPrice())
	assert.Equal(t, jobModel.Transaction.Gas, txRequest.Params.GetGas())
	assert.Equal(t, jobModel.Transaction.PrivateFor, txRequest.Params.GetPrivateFor())
	assert.Equal(t, jobModel.Transaction.PrivateFrom, txRequest.Params.GetPrivateFrom())
	assert.Equal(t, jobModel.Transaction.PrivacyGroupID, txRequest.Params.GetPrivacyGroupId())
}

func TestParsersJob_NewEnvelopeFromModelOneTimeKey(t *testing.T) {
	jobModel := testutils.FakeJobModel(1)
	jobModel.Annotations = &types.Annotations{
		OneTimeKey: true,
	}

	txEnvelope := NewEnvelopeFromJobModel(jobModel, map[string]string{})

	envelope, err := txEnvelope.Envelope()
	assert.NoError(t, err)
	assert.True(t, envelope.IsOneTimeKeySignature())
}
