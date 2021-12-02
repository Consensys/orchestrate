// +build unit

package parsers

import (
	"testing"

	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/types/testutils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	jobs, _ := NewJobEntitiesFromTxRequest(txReqEntity, chainUUID, hexutil.MustDecode("0x0ABC"))
	assert.Len(t, jobs, 1)

	job := jobs[0]
	assert.Equal(t, job.ScheduleUUID, txReqEntity.Schedule.UUID)
	assert.Equal(t, job.ChainUUID, chainUUID)
	assert.Equal(t, job.Type, entities.EthereumTransaction)
	assert.Equal(t, job.Labels, txReqEntity.Labels)

	assert.Equal(t, job.Transaction.From.Hex(), txReqEntity.Params.From.Hex())
	assert.Equal(t, job.Transaction.To.Hex(), txReqEntity.Params.To.Hex())
	assert.Equal(t, job.Transaction.Value, txReqEntity.Params.Value)
	assert.Equal(t, job.Transaction.GasPrice, txReqEntity.Params.GasPrice)
	assert.Equal(t, job.Transaction.Gas, txReqEntity.Params.Gas)
	assert.Equal(t, job.Transaction.Raw, txReqEntity.Params.Raw)
	assert.Equal(t, job.Transaction.PrivateFrom, txReqEntity.Params.PrivateFrom)
	assert.Equal(t, job.Transaction.PrivateFor, txReqEntity.Params.PrivateFor)
	assert.Equal(t, job.Transaction.PrivacyGroupID, txReqEntity.Params.PrivacyGroupID)
}

func TestParsersTxRequest_NewJobEntityFromSendRawTx(t *testing.T) {
	txReqEntity := testutils.FakeTxRequest()
	txReqEntity.Params.Raw = hexutil.MustDecode("0xf85380839896808252088083989680808216b4a0d35c752d3498e6f5ca1630d264802a992a141ca4b6a3f439d673c75e944e5fb0a05278aaa5fabbeac362c321b54e298dedae2d31471e432c26ea36a8d49cf08f1e")
	jobs, err := NewJobEntitiesFromTxRequest(txReqEntity, "", nil)
	require.NoError(t, err)
	require.Len(t, jobs, 1)

	job := jobs[0]
	assert.Equal(t, job.ScheduleUUID, txReqEntity.Schedule.UUID)
	assert.Equal(t, job.Type, entities.EthereumRawTransaction)
	assert.Equal(t, job.Labels, txReqEntity.Labels)

	assert.Equal(t, "0x7357589f8e367c2C31F51242fB77B350A11830F3", txReqEntity.Params.From.Hex())
}

func TestParsersTxRequest_NewEEAJobEntityFromSendTx(t *testing.T) {
	txReqEntity := testutils.FakeTxRequest()
	txReqEntity.Params.Protocol = entities.EEAChainType
	chainUUID := "chainUUID"
	jobs, _ := NewJobEntitiesFromTxRequest(txReqEntity, chainUUID, hexutil.MustDecode("0x0ABC"))
	assert.Len(t, jobs, 2)

	privJob := jobs[0]
	assert.Equal(t, privJob.Type, entities.EEAPrivateTransaction)
	assert.False(t, privJob.InternalData.OneTimeKey)

	markingJob := jobs[1]
	assert.Equal(t, markingJob.ScheduleUUID, txReqEntity.Schedule.UUID)
	assert.Equal(t, markingJob.ChainUUID, chainUUID)
	assert.Equal(t, markingJob.Type, entities.EEAMarkingTransaction)
	assert.Equal(t, markingJob.Labels, txReqEntity.Labels)
	assert.True(t, markingJob.InternalData.OneTimeKey)
}

func TestParsersTxRequest_NewTesseraJobEntityFromSendTx(t *testing.T) {
	txReqEntity := testutils.FakeTxRequest()
	txReqEntity.Params.Protocol = entities.TesseraChainType
	txReqEntity.Params.PrivateFor = []string{"0xPrivateFor"}
	chainUUID := "chainUUID"
	jobs, _ := NewJobEntitiesFromTxRequest(txReqEntity, chainUUID, hexutil.MustDecode("0x0ABC"))
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
	assert.Equal(t, markingJob.Transaction.From.Hex(), txReqEntity.Params.From.Hex())
}
