// +build unit

package parsers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/testutils"
)

func TestParsersTxRequest_NewModelFromEntity(t *testing.T) {
	reqHash := "reqHash"
	txReqEntity := testutils2.FakeTxRequestEntity()
	txReqModel, _ := NewTxRequestModelFromEntities(txReqEntity, reqHash)

	paramsBytes, _ := json.Marshal(txReqEntity.Params)
	assert.Equal(t, txReqEntity.IdempotencyKey, txReqModel.IdempotencyKey)
	assert.Equal(t, txReqEntity.CreatedAt, txReqModel.CreatedAt)
	assert.Equal(t, txReqEntity.UUID, txReqModel.UUID)
	assert.Equal(t, string(paramsBytes), txReqModel.Params)
}

func TestParsersTxRequest_NewJobEntityFromSendTx(t *testing.T) {
	txReqEntity := testutils2.FakeTxRequestEntity()
	chainUUID := "chainUUID"
	job := NewJobEntityFromTxRequest(txReqEntity, types.EthereumTransaction, chainUUID)

	assert.Equal(t, job.ScheduleUUID, txReqEntity.Schedule.UUID)
	assert.Equal(t, job.ChainUUID, chainUUID)
	assert.Equal(t, job.Type, types.EthereumTransaction)
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
