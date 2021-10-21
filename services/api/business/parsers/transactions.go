package parsers

import (
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/consensys/orchestrate/services/api/store/models"
	"github.com/ethereum/go-ethereum/core/types"
)

func NewTransactionModelFromEntities(tx *entities.ETHTransaction) *models.Transaction {
	return &models.Transaction{
		Hash:           tx.Hash,
		Sender:         tx.From,
		Recipient:      tx.To,
		Nonce:          tx.Nonce,
		Value:          tx.Value,
		GasPrice:       tx.GasPrice,
		GasFeeCap:      tx.GasFeeCap,
		GasTipCap:      tx.GasTipCap,
		Gas:            tx.Gas,
		Data:           tx.Data,
		Raw:            tx.Raw,
		TxType:         string(tx.TransactionType),
		AccessList:     tx.AccessList,
		PrivateFrom:    tx.PrivateFrom,
		PrivateFor:     tx.PrivateFor,
		PrivacyGroupID: tx.PrivacyGroupID,
		EnclaveKey:     tx.EnclaveKey,
		CreatedAt:      tx.CreatedAt,
		UpdatedAt:      tx.UpdatedAt,
	}
}

func NewTransactionEntityFromModels(tx *models.Transaction) *entities.ETHTransaction {
	accessList := types.AccessList{}
	_ = utils.CastInterfaceToObject(tx.AccessList, &accessList)

	return &entities.ETHTransaction{
		Hash:            tx.Hash,
		From:            tx.Sender,
		To:              tx.Recipient,
		Nonce:           tx.Nonce,
		Value:           tx.Value,
		GasPrice:        tx.GasPrice,
		Gas:             tx.Gas,
		GasTipCap:       tx.GasTipCap,
		GasFeeCap:       tx.GasFeeCap,
		Data:            tx.Data,
		TransactionType: entities.TransactionType(tx.TxType),
		AccessList:      accessList,
		PrivateFrom:     tx.PrivateFrom,
		PrivateFor:      tx.PrivateFor,
		PrivacyGroupID:  tx.PrivacyGroupID,
		EnclaveKey:      tx.EnclaveKey,
		Raw:             tx.Raw,
		CreatedAt:       tx.CreatedAt,
		UpdatedAt:       tx.UpdatedAt,
	}
}

func UpdateTransactionModelFromEntities(txModel *models.Transaction, tx *entities.ETHTransaction) {
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
	if tx.GasFeeCap != "" {
		txModel.GasFeeCap = tx.GasFeeCap
	}
	if tx.GasTipCap != "" {
		txModel.GasTipCap = tx.GasTipCap
	}
	if tx.Gas != "" {
		txModel.Gas = tx.Gas
	}
	if tx.Data != "" {
		txModel.Data = tx.Data
	}
	if tx.TransactionType != "" {
		txModel.TxType = string(tx.TransactionType)
	}
	if len(tx.AccessList) > 0 {
		txModel.AccessList = tx.AccessList
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
	if tx.EnclaveKey != "" {
		txModel.EnclaveKey = tx.EnclaveKey
	}
	if tx.Raw != "" {
		txModel.Raw = tx.Raw
	}
}
