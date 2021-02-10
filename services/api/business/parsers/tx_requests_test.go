// +build unit

package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
)

func TestParsersTxRequest_NewTxRequestModelFromEntities(t *testing.T) {
	reqHash := "reqHash"
	expectedScheduleID := 1
	txReqEntity := testutils.FakeTxRequest()
	txReqModel := NewTxRequestModelFromEntities(txReqEntity, reqHash, expectedScheduleID)

	assert.Equal(t, txReqEntity.IdempotencyKey, txReqModel.IdempotencyKey)
	assert.Equal(t, txReqEntity.CreatedAt, txReqModel.CreatedAt)
	assert.Equal(t, txReqEntity.Params, txReqModel.Params)
	assert.Equal(t, &expectedScheduleID, txReqModel.ScheduleID)
}

func TestParsersTxRequest_NewJobEntityFromSendTx(t *testing.T) {
	txReqEntity := testutils.FakeTxRequest()
	chainUUID := "chainUUID"
	jobs := NewJobEntitiesFromTxRequest(txReqEntity, chainUUID ,"0xDATA")
	assert.Len(t, jobs, 1)
	
	job := jobs[0]
	assert.Equal(t, job.ScheduleUUID, txReqEntity.Schedule.UUID)
	assert.Equal(t, job.ChainUUID, chainUUID)
	assert.Equal(t, job.Type, entities.EthereumTransaction)
	assert.Equal(t, job.Labels, txReqEntity.Labels)

	assert.Equal(t, job.Transaction.From, txReqEntity.Params.From)
	assert.Equal(t, job.Transaction.To, txReqEntity.Params.To)
	assert.Equal(t, job.Transaction.Value, txReqEntity.Params.Value)
	assert.Equal(t, job.Transaction.GasPrice, txReqEntity.Params.GasPrice)
	assert.Equal(t, job.Transaction.Gas, txReqEntity.Params.Gas)
	assert.Equal(t, job.Transaction.Raw, txReqEntity.Params.Raw)
	assert.Equal(t, job.Transaction.PrivateFrom, txReqEntity.Params.PrivateFrom)
	assert.Equal(t, job.Transaction.PrivateFor, txReqEntity.Params.PrivateFor)
	assert.Equal(t, job.Transaction.PrivacyGroupID, txReqEntity.Params.PrivacyGroupID)
}


func TestParsersTxRequest_NewOrionJobEntityFromSendTx(t *testing.T) {
	txReqEntity := testutils.FakeTxRequest()
	txReqEntity.Params.Protocol = entities.OrionChainType
	chainUUID := "chainUUID"
	jobs := NewJobEntitiesFromTxRequest(txReqEntity, chainUUID ,"0xDATA")
	assert.Len(t, jobs, 2)

	privJob := jobs[0]
	assert.Equal(t, privJob.Type, entities.OrionEEATransaction)
	assert.False(t, privJob.InternalData.OneTimeKey)
	
	markingJob := jobs[1]
	assert.Equal(t, markingJob.ScheduleUUID, txReqEntity.Schedule.UUID)
	assert.Equal(t, markingJob.ChainUUID, chainUUID)
	assert.Equal(t, markingJob.Type, entities.OrionMarkingTransaction)
	assert.Equal(t, markingJob.Labels, txReqEntity.Labels)
	assert.True(t, markingJob.InternalData.OneTimeKey)
}


func TestParsersTxRequest_NewTesseraJobEntityFromSendTx(t *testing.T) {
	txReqEntity := testutils.FakeTxRequest()
	txReqEntity.Params.Protocol = entities.TesseraChainType
	txReqEntity.Params.PrivateFor = []string{"0xPrivateFor"}
	chainUUID := "chainUUID"
	jobs := NewJobEntitiesFromTxRequest(txReqEntity, chainUUID ,"0xDATA")
	assert.Len(t, jobs, 2)

	privJob := jobs[0]
	assert.Equal(t, privJob.Type, entities.TesseraPrivateTransaction)
	assert.False(t, privJob.InternalData.OneTimeKey)

	markingJob := jobs[1]
	assert.Equal(t, markingJob.ScheduleUUID, txReqEntity.Schedule.UUID)
	assert.Equal(t, markingJob.ChainUUID, chainUUID)
	assert.Equal(t, markingJob.Type, entities.TesseraMarkingTransaction)
	assert.Equal(t, markingJob.Transaction.PrivateFor, txReqEntity.Params.PrivateFor)
	assert.Equal(t, markingJob.Labels, txReqEntity.Labels)
	assert.Equal(t, markingJob.InternalData.OneTimeKey, txReqEntity.InternalData.OneTimeKey)
	assert.Equal(t, markingJob.Transaction.From, txReqEntity.Params.From)
}
