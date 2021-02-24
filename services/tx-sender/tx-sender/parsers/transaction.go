package parsers

import (
	"math/big"
	"strconv"

	quorumtypes "github.com/consensys/quorum/core/types"

	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/consensys/quorum/common/hexutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func ETHTransactionToTransaction(tx *entities.ETHTransaction) *types.Transaction {
	// No need to validate the data as we know that internally the values are correct
	amount := new(big.Int)
	amount, _ = amount.SetString(tx.Value, 10)
	gasPrice := new(big.Int)
	gasPrice, _ = gasPrice.SetString(tx.GasPrice, 10)
	data, _ := hexutil.Decode(tx.Data)
	nonce, _ := strconv.ParseUint(tx.Nonce, 10, 64)
	gasLimit, _ := strconv.ParseUint(tx.Gas, 10, 64)

	if tx.To == "" {
		return types.NewContractCreation(nonce, amount, gasLimit, gasPrice, data)
	}

	to, _ := common.NewMixedcaseAddressFromString(tx.To)
	return types.NewTransaction(nonce, to.Address(), amount, gasLimit, gasPrice, data)
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
