package testutils

import (
	uuid "github.com/satori/go.uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
)

func FakeTxRequest() *models.TransactionRequest {
	return &models.TransactionRequest{
		IdempotencyKey: uuid.NewV4().String(),
		Chain:          uuid.NewV4().String(),
		Method:         types.MethodSendRawTransaction,
		Params:         "{\"field0\":\"field0Value\"}",
		Labels:         nil,
	}
}
