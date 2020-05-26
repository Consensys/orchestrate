// +build unit

package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/testutils"
)

func TestParsersTxRequest_NewModelFromEntity(t *testing.T) {
	reqHash := "reqHash"
	tenantID := "_"
	txReqEntity := testutils2.FakeTxRequestEntity()
	txReqModel, _ := NewTxRequestModelFromEntities(txReqEntity, reqHash, tenantID)

	paramsBytes, _ := json.Marshal(txReqEntity.Params)
	assert.Equal(t, txReqEntity.IdempotencyKey, txReqModel.IdempotencyKey)
	assert.Equal(t, txReqEntity.CreatedAt, txReqModel.CreatedAt)
	assert.Equal(t, string(paramsBytes), txReqModel.Params)
	assert.Equal(t, txReqEntity.Schedule.UUID, txReqModel.Schedules[0].UUID)
	assert.Equal(t, txReqEntity.Schedule.ChainUUID, txReqModel.Schedules[0].ChainUUID)
}


func TestParsersTxRequest_NewJobEntityFromSendTx(t *testing.T) {
	txReqEntity := testutils2.FakeTxRequestEntity()
	job, _:= NewJobEntityFromSendTxRequest(txReqEntity)
	
	assert.Equal(t, job.ScheduleUUID, txReqEntity.Schedule.UUID)
	assert.Equal(t, job.Type, entities.JobConstantinopleTransaction)
	assert.Equal(t, job.Labels, txReqEntity.Labels)
	
	assert.Equal(t, job.Transaction.From, *txReqEntity.Params.From)
	assert.Equal(t, job.Transaction.To, *txReqEntity.Params.To)
	assert.Equal(t, job.Transaction.Value, *txReqEntity.Params.Value)
	assert.Equal(t, job.Transaction.GasPrice, *txReqEntity.Params.GasPrice)
	assert.Equal(t, job.Transaction.GasLimit, *txReqEntity.Params.GasLimit)
	assert.NotEmpty(t, job.Transaction.Data)
}
