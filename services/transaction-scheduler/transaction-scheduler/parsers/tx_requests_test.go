// +build unit

package parsers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsersTxRequest_NewTxRequestModelFromEntities(t *testing.T) {
	reqHash := "reqHash"
	expectedScheduleID := 1
	txReqEntity := testutils.FakeTxRequest()
	txReqModel := NewTxRequestModelFromEntities(txReqEntity, reqHash, expectedScheduleID)

	assert.Equal(t, txReqEntity.IdempotencyKey, txReqModel.IdempotencyKey)
	assert.Equal(t, txReqEntity.CreatedAt, txReqModel.CreatedAt)
	assert.Equal(t, txReqEntity.UUID, txReqModel.UUID)
	assert.Equal(t, txReqEntity.Params, txReqModel.Params)
	assert.Equal(t, &expectedScheduleID, txReqModel.ScheduleID)
}

func TestParsersTxRequest_NewJobEntityFromSendTx(t *testing.T) {
	txReqEntity := testutils.FakeTxRequest()
	chainUUID := "chainUUID"
	job := NewJobEntityFromTxRequest(txReqEntity, utils.EthereumTransaction, chainUUID)

	assert.Equal(t, job.ScheduleUUID, txReqEntity.Schedule.UUID)
	assert.Equal(t, job.ChainUUID, chainUUID)
	assert.Equal(t, job.Type, utils.EthereumTransaction)
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
