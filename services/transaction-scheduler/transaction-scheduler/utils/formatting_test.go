// +build unit

package utils

import (
	"github.com/stretchr/testify/assert"
	storetestutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
	typestestutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types/testutils"
	"testing"
)

func TestUtils_ObjectToJSON(t *testing.T) {
	jsonStr, _ := ObjectToJSON(&types.TransactionRequest{
		BaseTransactionRequest: types.BaseTransactionRequest{
			IdempotencyKey: "myKey",
			ChainID:        "myChain",
		},
		Params: types.TransactionParams{
			From:            "from",
			To:              "to",
			MethodSignature: "constructor()",
		},
	})

	expectedJSONStr := "{\"idempotencyKey\":\"myKey\",\"chainID\":\"myChain\",\"params\":{\"from\":\"from\",\"to\":\"to\",\"methodSignature\":\"constructor()\"}}"
	assert.Equal(t, expectedJSONStr, jsonStr)
}

func TestUtils_FormatTxResponse(t *testing.T) {
	txRequest := storetestutils.FakeTxRequest(1)
	scheduleResponse := typestestutils.FakeScheduleResponse()

	txResponse, _ := FormatTxResponse(txRequest, scheduleResponse)

	assert.Equal(t, txRequest.IdempotencyKey, txResponse.IdempotencyKey)

	jsonMap := make(map[string]interface{})
	jsonMap["field0"] = "field0Value"
	assert.Equal(t, jsonMap, txResponse.Params)
}
