// +build unit

package utils

import (
	"context"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
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
	txRequest := testutils.FakeTxRequest()

	txResponse, _ := FormatTxResponse(context.Background(), txRequest)

	assert.Equal(t, txRequest.IdempotencyKey, txResponse.IdempotencyKey)
	assert.Equal(t, txRequest.Chain, txResponse.ChainID)
	assert.Equal(t, txRequest.Method, types.MethodSendRawTransaction)

	jsonMap := make(map[string]interface{})
	jsonMap["field0"] = "field0Value"
	assert.Equal(t, jsonMap, txResponse.Params)
}
