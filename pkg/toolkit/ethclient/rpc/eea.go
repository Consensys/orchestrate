package rpc

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/ethclient/utils"
	proto "github.com/consensys/orchestrate/pkg/types/ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Transaction Receipt
type privateReceipt struct {
	ContractAddress string       `json:"contractAddress,omitempty"`
	From            string       `json:"from,omitempty"`
	Output          string       `json:"output,omitempty"`
	CommitmentHash  string       `json:"commitmentHash,omitempty"`
	TransactionHash string       `json:"transactionHash,omitempty"`
	PrivateFrom     string       `json:"privateFrom,omitempty"`
	PrivateFor      []string     `json:"privateFor,omitempty"`
	PrivacyGroupID  string       `json:"privacyGroupId,omitempty"`
	Status          string       `json:"status,omitempty"`
	Logs            []*proto.Log `json:"logs,omitempty"`
}

// Distributes a signed, RLP encoded private transaction.
// https://besu.hyperledger.org/en/stable/Reference/API-Methods/#priv_distributerawtransaction
func (ec *Client) PrivDistributeRawTransaction(ctx context.Context, endpoint, raw string) (txHash ethcommon.Hash, err error) {
	err = ec.Call(ctx, endpoint, utils.ProcessResult(&txHash), "priv_distributeRawTransaction", raw)
	if err != nil {
		return ethcommon.Hash{}, errors.FromError(err).ExtendComponent(component)
	}
	return txHash, nil
}

func (ec *Client) PrivCreatePrivacyGroup(ctx context.Context, endpoint string, addresses []string) (string, error) {
	var privGroupID string
	err := ec.Call(ctx, endpoint, utils.ProcessResult(&privGroupID), "priv_createPrivacyGroup",
		map[string][]string{"addresses": addresses})
	if err != nil {
		return "", errors.FromError(err).ExtendComponent(component)
	}
	return privGroupID, nil
}

// PrivEEANonce Returns the private transaction count for specified account and privacy group
func (ec *Client) PrivEEANonce(ctx context.Context, endpoint string, account ethcommon.Address, privateFrom string, privateFor []string) (uint64, error) {
	var nonce hexutil.Uint64
	err := ec.Call(ctx, endpoint, utils.ProcessResult(&nonce), "priv_getEeaTransactionCount", account.Hex(), privateFrom, privateFor)
	if err != nil {
		return 0, errors.FromError(err).ExtendComponent(component)
	}
	return uint64(nonce), nil
}

// PrivNonce returns the private transaction count for the specified account and group of sender and recipients
func (ec *Client) PrivNonce(ctx context.Context, endpoint string, account ethcommon.Address, privacyGroupID string) (uint64, error) {
	var nonce hexutil.Uint64
	err := ec.Call(ctx, endpoint, utils.ProcessResult(&nonce), "priv_getTransactionCount", account.Hex(), privacyGroupID)
	if err != nil {
		return 0, errors.FromError(err).ExtendComponent(component)
	}
	return uint64(nonce), nil
}

// Returns a list of privacy groups containing only the listed members. For example, if the listed members are A and B, a privacy group containing A, B, and C is not returned.
func (ec *Client) PrivFindPrivacyGroup(ctx context.Context, endpoint string, members []string) ([]string, error) {
	var groupIDs []string
	err := ec.Call(ctx, endpoint, utils.ProcessResult(&members), "priv_findPrivacyGroup", members)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return groupIDs, nil
}

// PrivCodeAt returns the contract code of the given account.
// The block number can be nil, in which case the code is taken from the latest known block.
// https://besu.hyperledger.org/en/stable/Reference/API-Methods/#priv_getcode
func (ec *Client) PrivCodeAt(ctx context.Context, endpoint string, account ethcommon.Address, privateGroupID string, blockNumber *big.Int) ([]byte, error) {
	var code hexutil.Bytes
	err := ec.Call(ctx, endpoint, utils.ProcessResult(&code), "priv_getCode", privateGroupID, account, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return code, nil
}

func (ec *Client) EEAPrivPrecompiledContractAddr(ctx context.Context, endpoint string) (ethcommon.Address, error) {
	var hash string
	err := ec.Call(ctx, endpoint, utils.ProcessResult(&hash), "priv_getPrivacyPrecompileAddress")
	if err != nil {
		return ethcommon.Address{}, errors.FromError(err).ExtendComponent(component)
	}
	return ethcommon.HexToAddress(hash), nil
}

// TransactionReceipt returns the receipt of a transaction by transaction hash.
// Note that the receipt is not available for pending transactions.
func (ec *Client) PrivateTransactionReceipt(ctx context.Context, endpoint string, txHash ethcommon.Hash) (*proto.Receipt, error) {
	var r *proto.Receipt
	err := ec.Call(ctx, endpoint, utils.ProcessReceiptResult(&r), "eth_getTransactionReceipt", txHash)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	// https://besu.hyperledger.org/en/stable/Reference/API-Objects/#private-transaction-receipt-object
	// We do not need to retry for private Receipt as public receipt is available it means the private one
	// is too as private chain are implemented with instant finality
	var pr *privateReceipt
	err = ec.Call(ctx, endpoint, processPrivateReceiptResult(&pr), "priv_getTransactionReceipt", txHash)

	// In case private receipt is not available, we return the public receipt
	if err != nil && errors.IsInvalidParameterError(err) {
		return r, err
	} else if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	// Once we have both receipts, we create a hybrid version as follow
	r.Status, _ = hexutil.DecodeUint64(pr.Status)
	r.ContractAddress = pr.ContractAddress
	r.Logs = pr.Logs
	r.Output = pr.Output
	r.TxHash = pr.TransactionHash
	r.PrivateFrom = pr.PrivateFrom
	r.PrivateFor = pr.PrivateFor
	r.PrivacyGroupId = pr.PrivacyGroupID

	return r, nil
}

func processPrivateReceiptResult(receipt **privateReceipt) ProcessResultFunc {
	return func(result json.RawMessage) error {
		err := utils.ProcessResult(&receipt)(result)
		if err != nil {
			return err
		}

		if receipt == nil || *receipt == nil {
			// Receipt was not found
			return errors.NotFoundError("private receipt not found")
		}

		return nil
	}
}
