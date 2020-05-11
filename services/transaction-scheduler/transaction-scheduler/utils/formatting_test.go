// +build unit

package utils

import (
	"github.com/stretchr/testify/assert"
	storetestutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
	"testing"
)

func TestUtils_ObjectToJSON(t *testing.T) {
	jsonStr, _ := ObjectToJSON(&types.TransactionRequest{
		BaseTransactionRequest: types.BaseTransactionRequest{
			IdempotencyKey: "myKey",
			ChainUUID:      "myChain",
		},
		Params: types.TransactionParams{
			From:            "from",
			To:              "to",
			MethodSignature: "constructor()",
		},
	})

	expectedJSONStr := "{\"idempotencyKey\":\"myKey\",\"chainUUID\":\"myChain\",\"params\":{\"from\":\"from\",\"to\":\"to\",\"methodSignature\":\"constructor()\"}}"
	assert.Equal(t, expectedJSONStr, jsonStr)
}

func TestUtils_FormatTxResponse(t *testing.T) {
	txRequest := storetestutils.FakeTxRequest(1)

	txResponse, _ := FormatTxResponse(txRequest)

	assert.Equal(t, txRequest.IdempotencyKey, txResponse.IdempotencyKey)

	jsonMap := make(map[string]interface{})
	jsonMap["field0"] = "field0Value"
	assert.Equal(t, jsonMap, txResponse.Params)
}

func TestUtils_FormatScheduleResponse(t *testing.T) {
	schedule := storetestutils.FakeSchedule()

	scheduleResponse := FormatScheduleResponse(schedule)

	assert.Equal(t, schedule.UUID, scheduleResponse.UUID)
	assert.Equal(t, schedule.ChainUUID, scheduleResponse.ChainUUID)
	assert.Equal(t, schedule.CreatedAt, scheduleResponse.CreatedAt)
	assert.Equal(t, schedule.Jobs[0].UUID, scheduleResponse.Jobs[0].UUID)
}

func TestUtils_FormatJobResponse(t *testing.T) {
	job := storetestutils.FakeJob(1)

	jobResponse := FormatJobResponse(job)

	assert.Equal(t, job.UUID, jobResponse.UUID)
	assert.Equal(t, job.GetStatus(), jobResponse.Status)
	assert.Equal(t, job.CreatedAt, jobResponse.CreatedAt)
}
