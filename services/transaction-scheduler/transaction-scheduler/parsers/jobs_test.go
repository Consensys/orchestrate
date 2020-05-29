// +build unit

package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/testutils"
)

func TestParsersJob_NewModelFromEntity(t *testing.T) {
	jobEntity := testutils2.FakeJobEntity()
	jobModel := NewJobModelFromEntities(jobEntity, nil)
	finalJobEntity := NewJobEntityFromModels(jobModel)

	expectedJSON, _ := json.Marshal(jobEntity)
	actualJOSN, _ := json.Marshal(finalJobEntity)
	assert.Equal(t, string(expectedJSON), string(actualJOSN))
}

func TestParsersJob_NewEntityFromModel(t *testing.T) {
	jobModel := testutils.FakeJob(1)
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
	jobModel := testutils.FakeJob(1)
	headers := map[string]string{
		"Authorization": "Bearer MyToken",
	}
	txEnvelope := NewEnvelopeFromJobModel(jobModel, headers)

	txRequest := txEnvelope.GetTxRequest()
	assert.Equal(t, tx.JobType_ETH_TX, txRequest.JobType)
	assert.Equal(t, jobModel.Schedule.ChainUUID, txEnvelope.GetChainUUID())
	assert.Equal(t, jobModel.Transaction.Sender, txRequest.Params.From)
	assert.Equal(t, jobModel.Transaction.Recipient, txRequest.Params.To)
	assert.Equal(t, jobModel.Transaction.Data, txRequest.Params.Data)
	assert.Equal(t, jobModel.Transaction.Nonce, txRequest.Params.Nonce)
	assert.Equal(t, jobModel.Transaction.Raw, txRequest.Params.Raw)
	assert.Equal(t, jobModel.Transaction.GasPrice, txRequest.Params.GasPrice)
	assert.Equal(t, jobModel.Transaction.PrivateFor, txRequest.Params.PrivateFor)
	assert.Equal(t, jobModel.Transaction.PrivateFrom, txRequest.Params.PrivateFrom)
	assert.Equal(t, jobModel.Transaction.PrivacyGroupID, txRequest.Params.PrivacyGroupId)
	assert.Equal(t, headers, txRequest.Headers)
}
