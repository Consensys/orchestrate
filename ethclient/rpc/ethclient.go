package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"

	eth "github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/rpc/geth"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/types"
)

// Client is a connector to Ethereum blockchains that uses Geth rpc client
type Client struct {
	mux  *sync.Mutex
	rpcs map[string]rpc.Client

	conf *geth.Config
}

// NewClient creates a new MultiClient
func NewClient(conf *geth.Config) *Client {
	return &Client{
		mux:  &sync.Mutex{},
		rpcs: make(map[string]rpc.Client),
		conf: conf,
	}
}

// Dial an Ethereum client
func (ec *Client) Dial(ctx context.Context, rawurl string) (*big.Int, error) {
	// Dial using an rpc client
	c, err := geth.DialContext(ctx, rawurl, ec.conf)
	if err != nil {
		return nil, err
	}

	// Retrieve network version
	var version string
	if err := c.CallContext(ctx, &version, "net_version"); err != nil {
		return nil, err
	}

	chainID, ok := big.NewInt(0).SetString(version, 10)
	if !ok {
		return nil, fmt.Errorf("invalid net_version result %q", version)
	}

	// Register client
	ec.mux.Lock()
	ec.rpcs[chainID.Text(10)] = c
	ec.mux.Unlock()

	return chainID, nil
}

// Close client and all underlying Geth RPC client
func (ec *Client) Close() {
	for _, c := range ec.rpcs {
		go c.Close()
	}
}

func (ec *Client) getRPC(chainID *big.Int) (rpc.Client, error) {
	c, ok := ec.rpcs[chainID.Text(10)]
	if !ok {
		return nil, fmt.Errorf("no RPC connection registered for chain %q", chainID.Text(10))
	}
	return c, nil
}

type txExtraInfo struct {
	BlockNumber *string            `json:"blockNumber,omitempty"`
	BlockHash   *ethcommon.Hash    `json:"blockHash,omitempty"`
	From        *ethcommon.Address `json:"from,omitempty"`
}

type Body struct {
	Transactions []*ethtypes.Transaction `json:"transactions"`
}

func blockFromRaw(raw json.RawMessage) (*ethtypes.Block, error) {
	// Unmarshal block header information
	var header *ethtypes.Header
	if err := json.Unmarshal(raw, &header); err != nil {
		return nil, err
	}

	// Unmarshal block body information
	var body *Body
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, err
	}

	return ethtypes.NewBlock(header, body.Transactions, []*ethtypes.Header{}, []*ethtypes.Receipt{}), nil
}

// BlockByHash returns the given full block.
func (ec *Client) BlockByHash(ctx context.Context, chainID *big.Int, hash ethcommon.Hash) (*ethtypes.Block, error) {
	c, err := ec.getRPC(chainID)
	if err != nil {
		return nil, err
	}

	// Perform RPC call
	var raw json.RawMessage
	err = c.CallContext(ctx, &raw, "eth_getBlockByHash", hash, true)
	if err != nil {
		return nil, err
	} else if len(raw) == 0 {
		return nil, eth.NotFound
	}

	return blockFromRaw(raw)
}

// BlockByNumber returns a block from the current canonical chain. If number is nil, the
// latest known block is returned.
//
// Note that loading full blocks requires two requests. Use HeaderByNumber
// if you don't need all transactions or uncle headers.
func (ec *Client) BlockByNumber(ctx context.Context, chainID, number *big.Int) (*ethtypes.Block, error) {
	c, err := ec.getRPC(chainID)
	if err != nil {
		return nil, err
	}

	// Perform RPC call
	var raw json.RawMessage
	err = c.CallContext(ctx, &raw, "eth_getBlockByNumber", toBlockNumArg(number), true)
	if err != nil {
		return nil, err
	} else if len(raw) == 0 {
		return nil, eth.NotFound
	}

	return blockFromRaw(raw)
}

// HeaderByHash returns the block header with the given hash.
func (ec *Client) HeaderByHash(ctx context.Context, chainID *big.Int, hash ethcommon.Hash) (*ethtypes.Header, error) {
	c, err := ec.getRPC(chainID)
	if err != nil {
		return nil, err
	}

	var head *ethtypes.Header
	err = c.CallContext(ctx, &head, "eth_getBlockByHash", hash, false)
	if err == nil && head == nil {
		return nil, eth.NotFound
	}

	return head, err
}

// HeaderByNumber returns a block header from the current canonical chain. If number is
// nil, the latest known header is returned.
func (ec *Client) HeaderByNumber(ctx context.Context, chainID, number *big.Int) (*ethtypes.Header, error) {
	c, err := ec.getRPC(chainID)
	if err != nil {
		return nil, err
	}

	var head *ethtypes.Header
	err = c.CallContext(ctx, &head, "eth_getBlockByNumber", toBlockNumArg(number), false)
	if err == nil && head == nil {
		return nil, eth.NotFound
	}

	return head, err
}

// TransactionByHash returns the transaction with the given hash.
func (ec *Client) TransactionByHash(ctx context.Context, chainID *big.Int, hash ethcommon.Hash) (tx *ethtypes.Transaction, isPending bool, err error) {
	c, err := ec.getRPC(chainID)
	if err != nil {
		return nil, false, err
	}

	var raw json.RawMessage
	err = c.CallContext(ctx, &raw, "eth_getTransactionByHash", hash)
	if err != nil {
		return nil, false, err
	} else if len(raw) == 0 {
		return nil, false, eth.NotFound
	}

	if err := json.Unmarshal(raw, &tx); err != nil {
		return nil, false, err
	}

	if _, r, _ := tx.RawSignatureValues(); r == nil {
		return nil, false, fmt.Errorf("server returned transaction without signature")
	}

	// Unmarshal block body information
	var extra *txExtraInfo
	if err := json.Unmarshal(raw, &extra); err != nil {
		return nil, false, err
	}

	return tx, extra.BlockNumber == nil, nil
}

// TransactionReceipt returns the receipt of a transaction by transaction hash.
// Note that the receipt is not available for pending transactions.
func (ec *Client) TransactionReceipt(ctx context.Context, chainID *big.Int, txHash ethcommon.Hash) (*ethtypes.Receipt, error) {
	c, err := ec.getRPC(chainID)
	if err != nil {
		return nil, err
	}

	var r *ethtypes.Receipt
	err = c.CallContext(ctx, &r, "eth_getTransactionReceipt", txHash)
	if err == nil {
		if r == nil {
			return nil, eth.NotFound
		}
	}
	return r, err
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	return hexutil.EncodeBig(number)
}

// Networks return networks ID multi client is connected to
func (ec *Client) Networks(ctx context.Context) []*big.Int {
	networks := []*big.Int{}
	for _, c := range ec.rpcs {
		// Retrieve network version
		var version string
		if err := c.CallContext(ctx, &version, "net_version"); err != nil {
			continue
		}

		chain, ok := big.NewInt(0).SetString(version, 10)
		if !ok {
			continue
		}

		if chain != nil {
			networks = append(networks, chain)
		}
	}
	return networks
}

type Progress struct {
	StartingBlock hexutil.Uint64
	CurrentBlock  hexutil.Uint64
	HighestBlock  hexutil.Uint64
	PulledStates  hexutil.Uint64
	KnownStates   hexutil.Uint64
}

// SyncProgress retrieves the current progress of the sync algorithm. If there's
// no sync currently running, it returns nil.
func (ec *Client) SyncProgress(ctx context.Context, chainID *big.Int) (*eth.SyncProgress, error) {
	c, err := ec.getRPC(chainID)
	if err != nil {
		return nil, err
	}

	var raw json.RawMessage
	if err := c.CallContext(ctx, &raw, "eth_syncing"); err != nil {
		return nil, err
	}
	// Handle the possible response types
	var syncing bool
	if err := json.Unmarshal(raw, &syncing); err == nil {
		return nil, nil // Not syncing (always false)
	}

	var progress *Progress
	if err := json.Unmarshal(raw, &progress); err != nil {
		return nil, err
	}

	return &eth.SyncProgress{
		StartingBlock: uint64(progress.StartingBlock),
		CurrentBlock:  uint64(progress.CurrentBlock),
		HighestBlock:  uint64(progress.HighestBlock),
		PulledStates:  uint64(progress.PulledStates),
		KnownStates:   uint64(progress.KnownStates),
	}, nil
}

// State Access

// BalanceAt returns the wei balance of the given account.
// The block number can be nil, in which case the balance is taken from the latest known block.
func (ec *Client) BalanceAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, blockNumber *big.Int) (*big.Int, error) {
	c, err := ec.getRPC(chainID)
	if err != nil {
		return nil, err
	}

	var result hexutil.Big
	err = c.CallContext(ctx, &result, "eth_getBalance", account, toBlockNumArg(blockNumber))
	return (*big.Int)(&result), err
}

// StorageAt returns the value of key in the contract storage of the given account.
// The block number can be nil, in which case the value is taken from the latest known block.
func (ec *Client) StorageAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, key ethcommon.Hash, blockNumber *big.Int) ([]byte, error) {
	c, err := ec.getRPC(chainID)
	if err != nil {
		return nil, err
	}

	var result hexutil.Bytes
	err = c.CallContext(ctx, &result, "eth_getStorageAt", account, key, toBlockNumArg(blockNumber))
	return result, err
}

// CodeAt returns the contract code of the given account.
// The block number can be nil, in which case the code is taken from the latest known block.
func (ec *Client) CodeAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, blockNumber *big.Int) ([]byte, error) {
	c, err := ec.getRPC(chainID)
	if err != nil {
		return nil, err
	}

	var result hexutil.Bytes
	err = c.CallContext(ctx, &result, "eth_getCode", account, toBlockNumArg(blockNumber))
	return result, err
}

// NonceAt returns the account nonce of the given account.
// The block number can be nil, in which case the nonce is taken from the latest known block.
func (ec *Client) NonceAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, blockNumber *big.Int) (uint64, error) {
	c, err := ec.getRPC(chainID)
	if err != nil {
		return 0, err
	}

	var result hexutil.Uint64
	err = c.CallContext(ctx, &result, "eth_getTransactionCount", account, toBlockNumArg(blockNumber))
	return uint64(result), err
}

// Pending State

// PendingBalanceAt returns the wei balance of the given account in the pending state.
func (ec *Client) PendingBalanceAt(ctx context.Context, chainID *big.Int, account ethcommon.Address) (*big.Int, error) {
	c, err := ec.getRPC(chainID)
	if err != nil {
		return nil, err
	}

	var result hexutil.Big
	err = c.CallContext(ctx, &result, "eth_getBalance", account, "pending")
	return (*big.Int)(&result), err
}

// PendingStorageAt returns the value of key in the contract storage of the given account in the pending state.
func (ec *Client) PendingStorageAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, key ethcommon.Hash) ([]byte, error) {
	c, err := ec.getRPC(chainID)
	if err != nil {
		return nil, err
	}

	var result hexutil.Bytes
	err = c.CallContext(ctx, &result, "eth_getStorageAt", account, key, "pending")
	return result, err
}

// PendingCodeAt returns the contract code of the given account in the pending state.
func (ec *Client) PendingCodeAt(ctx context.Context, chainID *big.Int, account ethcommon.Address) ([]byte, error) {
	c, err := ec.getRPC(chainID)
	if err != nil {
		return nil, err
	}

	var result hexutil.Bytes
	err = c.CallContext(ctx, &result, "eth_getCode", account, "pending")
	return result, err
}

// PendingNonceAt returns the account nonce of the given account in the pending state.
// This is the nonce that should be used for the next transaction.
func (ec *Client) PendingNonceAt(ctx context.Context, chainID *big.Int, account ethcommon.Address) (uint64, error) {
	c, err := ec.getRPC(chainID)
	if err != nil {
		return 0, err
	}

	var result hexutil.Uint64
	err = c.CallContext(ctx, &result, "eth_getTransactionCount", account, "pending")
	return uint64(result), err
}

// Contract Calling

// CallContract executes a message call transaction, which is directly executed in the VM
// of the node, but never mined into the blockchain.
//
// blockNumber selects the block height at which the call runs. It can be nil, in which
// case the code is taken from the latest known block. Note that state from very old
// blocks might not be available.
func (ec *Client) CallContract(ctx context.Context, chainID *big.Int, msg *eth.CallMsg, blockNumber *big.Int) ([]byte, error) {
	c, err := ec.getRPC(chainID)
	if err != nil {
		return nil, err
	}

	var hex hexutil.Bytes
	err = c.CallContext(ctx, &hex, "eth_call", toCallArg(msg), toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}
	return hex, nil
}

// PendingCallContract executes a message call transaction using the EVM.
// The state seen by the contract call is the pending state.
func (ec *Client) PendingCallContract(ctx context.Context, chainID *big.Int, msg *eth.CallMsg) ([]byte, error) {
	c, err := ec.getRPC(chainID)
	if err != nil {
		return nil, err
	}

	var hex hexutil.Bytes
	err = c.CallContext(ctx, &hex, "eth_call", toCallArg(msg), "pending")
	if err != nil {
		return nil, err
	}
	return hex, nil
}

// SuggestGasPrice retrieves the currently suggested gas price to allow a timely
// execution of a transaction.
func (ec *Client) SuggestGasPrice(ctx context.Context, chainID *big.Int) (*big.Int, error) {
	c, err := ec.getRPC(chainID)
	if err != nil {
		return nil, err
	}

	var hex hexutil.Big
	err = c.CallContext(ctx, &hex, "eth_gasPrice")
	if err != nil {
		return nil, err
	}

	return (*big.Int)(&hex), nil
}

// EstimateGas tries to estimate the gas needed to execute a specific transaction based on
// the current pending state of the backend blockchain. There is no guarantee that this is
// the true gas limit requirement as other transactions may be added or removed by miners,
// but it should provide a basis for setting a reasonable default.
func (ec *Client) EstimateGas(ctx context.Context, chainID *big.Int, msg *eth.CallMsg) (uint64, error) {
	c, err := ec.getRPC(chainID)
	if err != nil {
		return 0, err
	}

	var hex hexutil.Uint64
	err = c.CallContext(ctx, &hex, "eth_estimateGas", toCallArg(msg))
	if err != nil {
		return 0, err
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
func (ec *Client) SendRawTransaction(ctx context.Context, chainID *big.Int, raw string) error {
	c, err := ec.getRPC(chainID)
	if err != nil {
		return err
	}

	return c.CallContext(ctx, nil, "eth_sendRawTransaction", raw)
}

// SendTransaction send transaction to Ethereum node
func (ec *Client) SendTransaction(ctx context.Context, chainID *big.Int, args *types.SendTxArgs) (txHash ethcommon.Hash, err error) {
	c, err := ec.getRPC(chainID)
	if err != nil {
		return ethcommon.Hash{}, err
	}

	err = c.CallContext(ctx, &txHash, "eth_sendTransaction", &args)
	if err != nil {
		return ethcommon.Hash{}, err
	}
	return txHash, nil
}

// SendRawPrivateTransaction send a raw transaction to a Ethreum node supporting privacy (e.g Quorum+Tessera node)
// TODO: to be implemented
func (ec *Client) SendRawPrivateTransaction(ctx context.Context, chainID *big.Int, raw string, args *types.PrivateArgs) (ethcommon.Hash, error) {
	return ethcommon.Hash{}, fmt.Errorf("%q is not implemented yet", "SendRawPrivateTransactionQuorum")
}
