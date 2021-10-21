package parsers

import (
	"math/big"
	"strconv"

	"github.com/consensys/orchestrate/pkg/utils"
	quorumtypes "github.com/consensys/quorum/core/types"

	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/quorum/common/hexutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func ETHTransactionToTransaction(tx *entities.ETHTransaction, chainIDStr string) *types.Transaction {
	var txData types.TxData

	// No need to validate the data as we know that internally the values are correct
	value, _ := new(big.Int).SetString(tx.Value, 10)
	gasPrice, _ := new(big.Int).SetString(tx.GasPrice, 10)
	data, _ := hexutil.Decode(tx.Data)
	nonce, _ := strconv.ParseUint(tx.Nonce, 10, 64)
	gasLimit, _ := strconv.ParseUint(tx.Gas, 10, 64)
	chainID, _ := new(big.Int).SetString(chainIDStr, 10)

	var toAddr *common.Address
	if tx.To != "" {
		toAddr = utils.ToPtr(common.HexToAddress(tx.To)).(*common.Address)
	}

	switch tx.TransactionType {
	case entities.LegacyTxType:
		txData = &types.LegacyTx{
			Nonce:    nonce,
			GasPrice: gasPrice,
			Gas:      gasLimit,
			To:       toAddr,
			Value:    value,
			Data:     data,
		}
	default:
		gasTipCap, _ := new(big.Int).SetString(tx.GasTipCap, 10)
		gasFeeCap, _ := new(big.Int).SetString(tx.GasFeeCap, 10)

		txData = &types.DynamicFeeTx{
			ChainID:    chainID,
			Nonce:      nonce,
			GasTipCap:  gasTipCap,
			GasFeeCap:  gasFeeCap,
			Gas:        gasLimit,
			To:         toAddr,
			Value:      value,
			Data:       data,
			AccessList: tx.AccessList,
		}
	}

	return types.NewTx(txData)
}

func ETHTransactionToQuorumTransaction(tx *entities.ETHTransaction) *quorumtypes.Transaction {
	// No need to validate the data as we know that internally the values are correct
	amount := new(big.Int)
	amount, _ = amount.SetString(tx.Value, 10)
	gasPrice := new(big.Int)
	gasPrice, _ = gasPrice.SetString(tx.GasPrice, 10)
	data, _ := hexutil.Decode(tx.Data)
	nonce, _ := strconv.ParseUint(tx.Nonce, 10, 64)
	gasLimit, _ := strconv.ParseUint(tx.Gas, 10, 64)

	if tx.To == "" {
		return quorumtypes.NewContractCreation(nonce, amount, gasLimit, gasPrice, data)
	}

	to, _ := common.NewMixedcaseAddressFromString(tx.To)
	return quorumtypes.NewTransaction(nonce, to.Address(), amount, gasLimit, gasPrice, data)
}
