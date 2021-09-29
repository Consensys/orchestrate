// +build unit

package parsers

import (
	testutils2 "github.com/consensys/orchestrate/pkg/types/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/consensys/orchestrate/services/api/store/models/testutils"

	"github.com/consensys/orchestrate/pkg/encoding/json"
)

func TestParsersTransaction_NewModelFromEntity(t *testing.T) {
	txEntity := testutils2.FakeETHTransaction()
	txModel := NewTransactionModelFromEntities(txEntity)
	finalTxEntity := NewTransactionEntityFromModels(txModel)

	expectedJSON, _ := json.Marshal(txEntity)
	actualJOSN, _ := json.Marshal(finalTxEntity)
	assert.Equal(t, string(expectedJSON), string(actualJOSN))
}

func TestParsersTransaction_NewEntityFromModel(t *testing.T) {
	txModel := testutils.FakeTransaction()
	txEntity := NewTransactionEntityFromModels(txModel)
	finalTxModel := NewTransactionModelFromEntities(txEntity)
	finalTxModel.UUID = txModel.UUID

	expectedJSON, _ := json.Marshal(txModel)
	actualJOSN, _ := json.Marshal(finalTxModel)
	assert.Equal(t, string(expectedJSON), string(actualJOSN))
}

func TestParsersTransaction_UpdateTransactionModel(t *testing.T) {
	txModel := testutils.FakeTransaction()
	txEntity := testutils2.FakeETHTransaction()
	UpdateTransactionModelFromEntities(txModel, txEntity)

	expectedTxModel := NewTransactionModelFromEntities(txEntity)
	expectedTxModel.UUID = txModel.UUID
	expectedTxModel.CreatedAt = txModel.CreatedAt
	expectedTxModel.UpdatedAt = txModel.UpdatedAt

	expectedJSON, _ := json.Marshal(txModel)
	actualJOSN, _ := json.Marshal(expectedTxModel)
	assert.Equal(t, string(expectedJSON), string(actualJOSN))
}
