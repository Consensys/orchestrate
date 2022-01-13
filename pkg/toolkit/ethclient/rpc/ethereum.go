package rpc

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/ethereum/types"
	"github.com/consensys/orchestrate/pkg/toolkit/ethclient/utils"
	proto "github.com/consensys/orchestrate/pkg/types/ethereum"
	eth "github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
)

// BlockByHash returns the given full block.
func (ec *Client) BlockByHash(ctx context.Context, endpoint string, hash ethcommon.Hash) (*ethtypes.Block, error) {
	// Perform RPC call
	var header *ethtypes.Header
	var body *Body
	err := ec.Call(ctx, endpoint, processBlockResult(&header, &body), "eth_getBlockByHash", hash, true)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return ethtypes.NewBlock(header, body.Transactions, []*ethtypes.Header{}, []*ethtypes.Receipt{}, new(trie.Trie)), nil
}

// BlockByNumber returns a block from the current canonical chain. If number is nil, the
// latest known block is returned.
//
// Note that loading full blocks requires two requests. Use HeaderByNumber
// if you don't need all transactions or uncle headers.
func (ec *Client) BlockByNumber(ctx context.Context, endpoint string, number *big.Int) (*ethtypes.Block, error) {
	// Perform RPC call
	var header *ethtypes.Header
	var body *Body
	err := ec.Call(ctx, endpoint, processBlockResult(&header, &body), "eth_getBlockByNumber", toBlockNumArg(number), true)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return ethtypes.NewBlock(header, body.Transactions, []*ethtypes.Header{}, []*ethtypes.Receipt{}, new(trie.Trie)), nil
}

func processHeaderResult(head **ethtypes.Header) ParseResultFunc {
	return func(result json.RawMessage) error {
		err := utils.ProcessResult(head)(result)
		if err != nil {
			return err
		}

		if *head == nil {
			// Block was not found
			return errors.NotFoundError("block not found")
		}

		return nil
	}
}

// HeaderByHash returns the block header with the given hash.
func (ec *Client) HeaderByHash(ctx context.Context, endpoint string, hash ethcommon.Hash) (*ethtypes.Header, error) {
	var head *ethtypes.Header
	err := ec.Call(ctx, endpoint, processHeaderResult(&head), "eth_getBlockByHash", hash, false)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return head, nil
}

// HeaderByNumber returns a block header from the current canonical chain. If number is
// nil, the latest known header is returned.
func (ec *Client) HeaderByNumber(ctx context.Context, endpoint string, number *big.Int) (*ethtypes.Header, error) {
	var head *ethtypes.Header
	err := ec.Call(ctx, endpoint, processHeaderResult(&head), "eth_getBlockByNumber", toBlockNumArg(number), false)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return head, nil
}

func processTxResult(tx **ethtypes.Transaction, extra **txExtraInfo) ParseResultFunc {
	return func(result json.RawMessage) error {
		var raw json.RawMessage
		err := utils.ProcessResult(&raw)(result)
		if err != nil {
			return err
		}

		if len(raw) == 0 {
			// Block was not found
			return errors.NotFoundError("transaction not found in the body of the response")
		}

		if err := json.Unmarshal(raw, tx); err != nil {
			return errors.FromError(err)
		}

		if _, r, _ := (*tx).RawSignatureValues(); r == nil {
			return errors.DataCorruptedError("transaction without signature")
		}

		// Unmarshal block body information
		if err := json.Unmarshal(raw, extra); err != nil {
			return errors.FromError(err)
		}

		return nil
	}
}

// TransactionByHash returns the transaction with the given hash.
func (ec *Client) TransactionByHash(ctx context.Context, endpoint string, hash ethcommon.Hash) (*ethtypes.Transaction, bool, error) {
	var tx *ethtypes.Transaction
	var extra *txExtraInfo
	err := ec.Call(ctx, endpoint, processTxResult(&tx, &extra), "eth_getTransactionByHash", hash)
	if err != nil {
		return nil, false, errors.FromError(err).ExtendComponent(component)
	}

	return tx, extra.BlockNumber == nil, nil
}

// TransactionReceipt returns the receipt of a transaction by transaction hash.
// Note that the receipt is not available for pending transactions.
func (ec *Client) TransactionReceipt(ctx context.Context, endpoint string, txHash ethcommon.Hash) (*proto.Receipt, error) {
	var r *proto.Receipt
	err := ec.Call(ctx, endpoint, utils.ProcessReceiptResult(&r), "eth_getTransactionReceipt", txHash)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return r, nil
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	return hexutil.EncodeBig(number)
}

type Progress struct {
	StartingBlock hexutil.Uint64
	CurrentBlock  hexutil.Uint64
	HighestBlock  hexutil.Uint64
	PulledStates  hexutil.Uint64
	KnownStates   hexutil.Uint64
}

func processProgressResult(progress **Progress) ParseResultFunc {
	return func(result json.RawMessage) error {
		var raw json.RawMessage
		err := utils.ProcessResult(&raw)(result)
		if err != nil {
			return err
		}

		var syncing bool
		if err = json.Unmarshal(raw, &syncing); err == nil {
			return nil
		}

		err = json.Unmarshal(raw, progress)
		if err != nil {
			return errors.EncodingError(err.Error())
		}

		return nil
	}
}

// SyncProgress retrieves the current progress of the sync algorithm. If there's
// no sync currently running, it returns nil.
func (ec *Client) SyncProgress(ctx context.Context, endpoint string) (*eth.SyncProgress, error) {
	var progress *Progress
	if err := ec.Call(ctx, endpoint, processProgressResult(&progress), "eth_syncing"); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	if progress == nil {
		return nil, nil
	}

	return &eth.SyncProgress{
		StartingBlock: uint64(progress.StartingBlock),
		CurrentBlock:  uint64(progress.CurrentBlock),
		HighestBlock:  uint64(progress.HighestBlock),
		PulledStates:  uint64(progress.PulledStates),
		KnownStates:   uint64(progress.KnownStates),
	}, nil
}

// SendRawPrivateTransaction send a raw transaction to an Ethereum node supporting EEA extension
func (ec *Client) Network(ctx context.Context, endpoint string) (*big.Int, error) {
	var version string
	if err := ec.Call(ctx, endpoint, utils.ProcessResult(&version), "net_version"); err != nil {
		return nil, err
	}

	chain, ok := big.NewInt(0).SetString(version, 10)
	if !ok {
		return nil, errors.EncodingError("invalid network id %q", version)
	}

	return chain, nil
}

// State Access

// BalanceAt returns the wei balance of the given account.
// The block number can be nil, in which case the balance is taken from the latest known block.
func (ec *Client) BalanceAt(ctx context.Context, endpoint string, account ethcommon.Address, blockNumber *big.Int) (*big.Int, error) {
	var balance hexutil.Big
	err := ec.Call(ctx, endpoint, utils.ProcessResult(&balance), "eth_getBalance", account, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return (*big.Int)(&balance), nil
}

// StorageAt returns the value of key in the contract storage of the given account.
// The block number can be nil, in which case the value is taken from the latest known block.
func (ec *Client) StorageAt(ctx context.Context, endpoint string, account ethcommon.Address, key ethcommon.Hash, blockNumber *big.Int) ([]byte, error) {
	var storage hexutil.Bytes
	err := ec.Call(ctx, endpoint, utils.ProcessResult(&storage), "eth_getStorageAt", account, key, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return storage, nil
}

// CodeAt returns the contract code of the given account.
// The block number can be nil, in which case the code is taken from the latest known block.
func (ec *Client) CodeAt(ctx context.Context, endpoint string, account ethcommon.Address, blockNumber *big.Int) ([]byte, error) {
	var code hexutil.Bytes
	err := ec.Call(ctx, endpoint, utils.ProcessResult(&code), "eth_getCode", account, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return code, nil
}

// NonceAt returns the account nonce of the given account.
// The block number can be nil, in which case the nonce is taken from the latest known block.
func (ec *Client) NonceAt(ctx context.Context, endpoint string, account ethcommon.Address, blockNumber *big.Int) (uint64, error) {
	var nonce hexutil.Uint64
	err := ec.Call(ctx, endpoint, utils.ProcessResult(&nonce), "eth_getTransactionCount", account, toBlockNumArg(blockNumber))
	if err != nil {
		return 0, errors.FromError(err).ExtendComponent(component)
	}
	return uint64(nonce), nil
}

// Pending State

// PendingBalanceAt returns the wei balance of the given account in the pending state.
func (ec *Client) PendingBalanceAt(ctx context.Context, endpoint string, account ethcommon.Address) (*big.Int, error) {
	var balance hexutil.Big
	err := ec.Call(ctx, endpoint, utils.ProcessResult(&balance), "eth_getBalance", account, "pending")
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return (*big.Int)(&balance), nil
}

// PendingStorageAt returns the value of key in the contract storage of the given account in the pending state.
func (ec *Client) PendingStorageAt(ctx context.Context, endpoint string, account ethcommon.Address, key ethcommon.Hash) ([]byte, error) {
	var storage hexutil.Bytes
	err := ec.Call(ctx, endpoint, utils.ProcessResult(&storage), "eth_getStorageAt", account, key, "pending")
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return storage, nil
}

// PendingCodeAt returns the contract code of the given account in the pending state.
func (ec *Client) PendingCodeAt(ctx context.Context, endpoint string, account ethcommon.Address) ([]byte, error) {
	var code hexutil.Bytes
	err := ec.Call(ctx, endpoint, utils.ProcessResult(&code), "eth_getCode", account, "pending")
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return code, nil
}

// PendingNonceAt returns the account nonce of the given account in the pending state.
// This is the nonce that should be used for the next transaction.
func (ec *Client) PendingNonceAt(ctx context.Context, endpoint string, account ethcommon.Address) (uint64, error) {
	var nonce hexutil.Uint64
	err := ec.Call(ctx, endpoint, utils.ProcessResult(&nonce), "eth_getTransactionCount", account, "pending")
	if err != nil {
		return 0, errors.FromError(err).ExtendComponent(component)
	}
	return uint64(nonce), nil
}

// Contract Calling

// CallContract executes a message call transaction, which is directly executed in the VM
// of the node, but never mined into the blockchain.
//
// blockNumber selects the block height at which the call runs. It can be nil, in which
// case the code is taken from the latest known block. Note that state from very old
// blocks might not be available.
func (ec *Client) CallContract(ctx context.Context, endpoint string, msg *eth.CallMsg, blockNumber *big.Int) ([]byte, error) {
	var hex hexutil.Bytes
	err := ec.Call(ctx, endpoint, utils.ProcessResult(&hex), "eth_call", toCallArg(msg), toBlockNumArg(blockNumber))
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return hex, nil
}

// PendingCallContract executes a message call transaction using the EVM.
// The state seen by the contract call is the pending state.
func (ec *Client) PendingCallContract(ctx context.Context, endpoint string, msg *eth.CallMsg) ([]byte, error) {
	var hex hexutil.Bytes
	err := ec.Call(ctx, endpoint, utils.ProcessResult(&hex), "eth_call", toCallArg(msg), "pending")
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return hex, nil
}

// SuggestGasPrice retrieves the currently suggested gas price to allow a timely
// execution of a transaction.
func (ec *Client) SuggestGasPrice(ctx context.Context, endpoint string) (*big.Int, error) {
	var hex hexutil.Big
	err := ec.Call(ctx, endpoint, utils.ProcessResult(&hex), "eth_gasPrice")
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return hex.ToInt(), nil
}

type FeeHistory struct {
	OldestBlock   hexutil.Big
	Reward        [][]hexutil.Big
	BaseFeePerGas []hexutil.Big
	GasUsedRatio  []float64
}

func parseFeeHistoryResult(feeHistory **FeeHistory) ParseResultFunc {
	return func(result json.RawMessage) error {
		var raw json.RawMessage
		err := utils.ProcessResult(&raw)(result)
		if err != nil {
			return err
		}

		err = json.Unmarshal(raw, feeHistory)
		if err != nil {
			return errors.EncodingError(err.Error())
		}

		return nil
	}
}

func (ec *Client) FeeHistory(ctx context.Context, endpoint string, blockCount int, newestBlock string) (*FeeHistory, error) {
	var feeHistory *FeeHistory
	if err := ec.Call(ctx, endpoint, parseFeeHistoryResult(&feeHistory), "eth_feeHistory", blockCount, newestBlock, []interface{}{}); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return feeHistory, nil
}

// EstimateGas tries to estimate the gas needed to execute a specific transaction based on
// the current pending state of the backend blockchain. There is no guarantee that this is
// the true gas limit requirement as other transactions may be added or removed by miners,
// but it should provide a basis for setting a reasonable default.
func (ec *Client) EstimateGas(ctx context.Context, endpoint string, msg *eth.CallMsg) (uint64, error) {
	var hex hexutil.Uint64
	err := ec.Call(ctx, endpoint, utils.ProcessResult(&hex), "eth_estimateGas", toCallArg(msg))
	if err != nil {
		return 0, errors.FromError(err).ExtendComponent(component)
	}
	return uint64(hex), nil
}

func toCallArg(msg *eth.CallMsg) interface{} {
	arg := map[string]interface{}{
		"from": msg.From,
		"to":   msg.To,
	}
	if len(msg.Data) > 0 {
		arg["data"] = hexutil.Bytes(msg.Data)
	}
	if msg.Value != nil {
		arg["value"] = (*hexutil.Big)(msg.Value)
	}
	if msg.Gas != 0 {
		arg["gas"] = hexutil.Uint64(msg.Gas)
	}
	if msg.GasPrice != nil {
		arg["gasPrice"] = (*hexutil.Big)(msg.GasPrice)
	}
	return arg
}

// SendRawTransaction allows to send a raw transaction
func (ec *Client) SendRawTransaction(ctx context.Context, endpoint string, raw hexutil.Bytes) (txHash ethcommon.Hash, err error) {
	err = ec.Call(ctx, endpoint, utils.ProcessResult(&txHash), "eth_sendRawTransaction", raw)
	if err != nil {
		return ethcommon.Hash{}, errors.FromError(err).ExtendComponent(component)
	}
	return txHash, nil
}

// SendTransaction send transaction to an Ethereum node
func (ec *Client) SendTransaction(ctx context.Context, endpoint string, args *types.SendTxArgs) (txHash ethcommon.Hash, err error) {
	err = ec.Call(ctx, endpoint, utils.ProcessResult(&txHash), "eth_sendTransaction", &args)
	if err != nil {
		return ethcommon.Hash{}, errors.FromError(err).ExtendComponent(component)
	}
	return txHash, nil
}

// SendRawPrivateTransaction send a raw transaction to an Ethereum node supporting EEA extension
func (ec *Client) SendRawPrivateTransaction(ctx context.Context, endpoint string, raw hexutil.Bytes) (ethcommon.Hash, error) {
	// Send a raw signed transactions using EEA extension method
	// MethodSignature documentation here: https://besu.hyperledger.org/en/stable/Reference/API-Methods/#eea_sendrawtransaction
	var hash string
	err := ec.Call(ctx, endpoint, utils.ProcessResult(&hash), "eea_sendRawTransaction", raw)
	if err != nil {
		return ethcommon.Hash{}, errors.FromError(err).ExtendComponent(component)
	}
	return ethcommon.HexToHash(hash), nil
}
