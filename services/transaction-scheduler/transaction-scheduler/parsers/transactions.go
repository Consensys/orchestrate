package parsers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
)

func NewTransactionModelFromEntities(tx *types.ETHTransaction) *models.Transaction {
	return &models.Transaction{
		Hash:           tx.Hash,
		Sender:         tx.From,
		Recipient:      tx.To,
		Nonce:          tx.Nonce,
		Value:          tx.Value,
		GasPrice:       tx.GasPrice,
		GasLimit:       tx.GasLimit,
		Data:           tx.Data,
		Raw:            tx.Raw,
		PrivateFrom:    tx.PrivateFrom,
		PrivateFor:     tx.PrivateFor,
		PrivacyGroupID: tx.PrivacyGroupID,
		CreatedAt:      tx.CreatedAt,
		UpdatedAt:      tx.UpdatedAt,
	}
}

func NewTransactionEntityFromModels(tx *models.Transaction) *types.ETHTransaction {
	return &types.ETHTransaction{
		Hash:           tx.Hash,
		From:           tx.Sender,
		To:             tx.Recipient,
		Nonce:          tx.Nonce,
		Value:          tx.Value,
		GasPrice:       tx.GasPrice,
		GasLimit:       tx.GasLimit,
		Data:           tx.Data,
		PrivateFrom:    tx.PrivateFrom,
		PrivateFor:     tx.PrivateFor,
		PrivacyGroupID: tx.PrivacyGroupID,
		Raw:            tx.Raw,
		CreatedAt:      tx.CreatedAt,
		UpdatedAt:      tx.UpdatedAt,
	}
}

func UpdateTransactionModelFromEntities(txModel *models.Transaction, tx *types.ETHTransaction) {
	if tx.Hash != "" {
		txModel.Hash = tx.Hash
	}
	if tx.From != "" {
		txModel.Sender = tx.From
	}
	if tx.To != "" {
		txModel.Recipient = tx.To
	}
	if tx.Nonce != "" {
		txModel.Nonce = tx.Nonce
	}
	if tx.Value != "" {
		txModel.Value = tx.Value
	}
	if tx.GasPrice != "" {
		txModel.GasPrice = tx.GasPrice
	}
	if tx.GasLimit != "" {
		txModel.GasLimit = tx.GasLimit
	}
	if tx.Data != "" {
		txModel.Data = tx.Data
	}
	if tx.PrivateFrom != "" {
		txModel.PrivateFrom = tx.PrivateFrom
	}
	if len(tx.PrivateFor) > 0 {
		txModel.PrivateFor = tx.PrivateFor
	}
	if tx.PrivateFrom != "" {
		txModel.PrivateFrom = tx.PrivateFrom
	}
	if tx.PrivacyGroupID != "" {
		txModel.PrivacyGroupID = tx.PrivacyGroupID
	}
	if tx.Raw != "" {
		txModel.Raw = tx.Raw
	}
}
