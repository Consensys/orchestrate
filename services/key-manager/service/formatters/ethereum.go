package formatters

import (
	"math/big"

	quorumtypes "github.com/consensys/quorum/core/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/keymanager/ethereum"
)

func FormatETHAccountResponse(account *entities.ETHAccount) *types.ETHAccountResponse {
	return &types.ETHAccountResponse{
		Address:             account.Address,
		PublicKey:           account.PublicKey,
		CompressedPublicKey: account.CompressedPublicKey,
		Namespace:           account.Namespace,
	}
}

func FormatSignETHTransactionRequest(request *types.SignETHTransactionRequest) *ethtypes.Transaction {
	// No need to check the "ok" values because we know that at the fields are valid big ints and hex string,
	// this also avoids this function returning an error
	amount, _ := new(big.Int).SetString(request.Amount, 10)
	gasPrice, _ := new(big.Int).SetString(request.GasPrice, 10)
	data, _ := hexutil.Decode(request.Data)

	if request.To == "" {
		return ethtypes.NewContractCreation(request.Nonce, amount, request.GasLimit, gasPrice, data)
	}

	return ethtypes.NewTransaction(request.Nonce, common.HexToAddress(request.To), amount, request.GasLimit, gasPrice, data)
}

func FormatSignTesseraTransactionRequest(request *types.SignTesseraTransactionRequest) *quorumtypes.Transaction {
	// No need to check the "ok" values because we know that at the fields are valid big ints and hex string,
	// this also avoids this function returning an error
	amount, _ := new(big.Int).SetString(request.Amount, 10)
	gasPrice, _ := new(big.Int).SetString(request.GasPrice, 10)
	data, _ := hexutil.Decode(request.Data)

	if request.To == "" {
		return quorumtypes.NewContractCreation(request.Nonce, amount, request.GasLimit, gasPrice, data)
	}

	return quorumtypes.NewTransaction(request.Nonce, common.HexToAddress(request.To), amount, request.GasLimit, gasPrice, data)
}
