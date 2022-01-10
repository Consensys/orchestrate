package formatters

import (
	"math/big"

	"github.com/consensys/orchestrate/pkg/types/entities"
	quorumtypes "github.com/consensys/quorum/core/types"
	"github.com/ethereum/go-ethereum/core/types"
)

func ETHTransactionToTransaction(tx *entities.ETHTransaction, chainID *big.Int) *types.Transaction {
	var txData types.TxData

	var value *big.Int
	if tx.Value != nil {
		value = tx.Value.ToInt()
	}

	var gasPrice *big.Int
	if tx.GasPrice != nil {
		gasPrice = tx.GasPrice.ToInt()
	}
	var nonce uint64
	if tx.Nonce != nil {
		nonce = *tx.Nonce
	}
	var gasLimit uint64
	if tx.Gas != nil {
		gasLimit = *tx.Gas
	}

	switch tx.TransactionType {
	case entities.LegacyTxType:
		txData = &types.LegacyTx{
			Nonce:    nonce,
			GasPrice: gasPrice,
			Gas:      gasLimit,
			To:       tx.To,
			Value:    value,
			Data:     tx.Data,
		}
	default:
		var gasTipCap *big.Int
		if tx.GasTipCap != nil {
			gasTipCap = tx.GasTipCap.ToInt()
		}

		var gasFeeCap *big.Int
		if tx.GasFeeCap != nil {
			gasFeeCap = tx.GasFeeCap.ToInt()
		}

		txData = &types.DynamicFeeTx{
			ChainID:    chainID,
			Nonce:      nonce,
			GasTipCap:  gasTipCap,
			GasFeeCap:  gasFeeCap,
			Gas:        gasLimit,
			To:         tx.To,
			Value:      value,
			Data:       tx.Data,
			AccessList: tx.AccessList,
		}
	}

	return types.NewTx(txData)
}

func ETHTransactionToQuorumTransaction(tx *entities.ETHTransaction) *quorumtypes.Transaction {
	var value *big.Int
	if tx.Value != nil {
		value = tx.Value.ToInt()
	}

	var gasPrice *big.Int
	if tx.GasPrice != nil {
		gasPrice = tx.GasPrice.ToInt()
	}
	var nonce uint64
	if tx.Nonce != nil {
		nonce = *tx.Nonce
	}
	var gasLimit uint64
	if tx.Gas != nil {
		gasLimit = *tx.Gas
	}

	if tx.To == nil {
		return quorumtypes.NewContractCreation(nonce, value, gasLimit, gasPrice, tx.Data)
	}

	return quorumtypes.NewTransaction(nonce, *tx.To, value, gasLimit, gasPrice, tx.Data)
}
