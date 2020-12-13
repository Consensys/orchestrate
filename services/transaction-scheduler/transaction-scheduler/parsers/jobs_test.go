// +build unit

package parsers

import (
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils/envelope"

	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/tx"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/store/models/testutils"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/json"
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
	assert.Equal(t, finalJobModel.NextJobUUID, jobModel.NextJobUUID)
	assert.Equal(t, finalJobModel.Type, jobModel.Type)
	assert.Equal(t, finalJobModel.Labels, jobModel.Labels)
	assert.Equal(t, finalJobModel.CreatedAt, jobModel.CreatedAt)
}

func TestParsersJob_NewEnvelopeFromModel(t *testing.T) {
	job := testutils2.FakeJob()
	headers := map[string]string{
		"Authorization": "Bearer MyToken",
	}
	txEnvelope := envelope.NewEnvelopeFromJob(job, headers)

	txRequest := txEnvelope.GetTxRequest()
	assert.Equal(t, job.ChainUUID, txEnvelope.GetChainUUID())
	assert.Equal(t, tx.JobTypeMap[job.Type], txRequest.GetJobType())
	assert.Equal(t, job.Transaction.From, txRequest.Params.GetFrom())
	assert.Equal(t, job.Transaction.To, txRequest.Params.GetTo())
	assert.Equal(t, job.Transaction.Data, txRequest.Params.GetData())
	assert.Equal(t, job.Transaction.Nonce, txRequest.Params.GetNonce())
	assert.Equal(t, job.Transaction.Raw, txRequest.Params.GetRaw())
	assert.Equal(t, job.Transaction.GasPrice, txRequest.Params.GetGasPrice())
	assert.Equal(t, job.Transaction.Gas, txRequest.Params.GetGas())
	assert.Equal(t, job.Transaction.PrivateFor, txRequest.Params.GetPrivateFor())
	assert.Equal(t, job.Transaction.PrivateFrom, txRequest.Params.GetPrivateFrom())
	assert.Equal(t, job.Transaction.PrivacyGroupID, txRequest.Params.GetPrivacyGroupId())
	assert.Equal(t, job.InternalData.ChainID, txEnvelope.GetChainID())
}

func TestParsersJob_NewEnvelopeFromModelOneTimeKey(t *testing.T) {
	job := testutils2.FakeJob()
	job.InternalData = &entities.InternalData{
		OneTimeKey: true,
	}

	txEnvelope := envelope.NewEnvelopeFromJob(job, map[string]string{})

	evlp, err := txEnvelope.Envelope()
	assert.NoError(t, err)
	assert.True(t, evlp.IsOneTimeKeySignature())
}
