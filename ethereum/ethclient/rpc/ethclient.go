package rpc

import (
	"context"
	"encoding/json"
	"math/big"
	"sync"

	eth "github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/rpc/geth"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/types"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/errors"
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
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	// Retrieve network version
	var version string
	if err = c.CallContext(ctx, &version, "net_version"); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	chainID, ok := big.NewInt(0).SetString(version, 10)
	if !ok {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	// Register client
	ec.mux.Lock()
	ec.rpcs[chainID.Text(10)] = c
	ec.mux.Unlock()

	return chainID, nil
}

// Close client and all underlying Geth RPC client
func (ec *Client) Close() {
	log.Infof("Closing RPC clients. Number of clients %d", len(ec.rpcs))
	for _, c := range ec.rpcs {
		go c.Close()
	}
}

func (ec *Client) getRPC(chainID *big.Int) rpc.Client {
	c, ok := ec.rpcs[chainID.Text(10)]
	if ok {
		return c
	}

	nullClient := geth.CreateNullClient(chainID)
	ec.rpcs[chainID.String()] = nullClient
	return nullClient
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
	if err := encoding.Unmarshal(raw, &header); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	// Unmarshal block body information
	var body *Body
	if err := encoding.Unmarshal(raw, &body); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return ethtypes.NewBlock(header, body.Transactions, []*ethtypes.Header{}, []*ethtypes.Receipt{}), nil
}

// BlockByHash returns the given full block.
func (ec *Client) BlockByHash(ctx context.Context, chainID *big.Int, hash ethcommon.Hash) (*ethtypes.Block, error) {
	// Perform RPC call
	var raw json.RawMessage
	err := ec.getRPC(chainID).CallContext(ctx, &raw, "eth_getBlockByHash", hash, true)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return blockFromRaw(raw)
}

// BlockByNumber returns a block from the current canonical chain. If number is nil, the
// latest known block is returned.
//
// Note that loading full blocks requires two requests. Use HeaderByNumber
// if you don't need all transactions or uncle headers.
func (ec *Client) BlockByNumber(ctx context.Context, chainID, number *big.Int) (*ethtypes.Block, error) {
	// Perform RPC call
	var raw json.RawMessage
	err := ec.getRPC(chainID).CallContext(ctx, &raw, "eth_getBlockByNumber", toBlockNumArg(number), true)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return blockFromRaw(raw)
}

// HeaderByHash returns the block header with the given hash.
func (ec *Client) HeaderByHash(ctx context.Context, chainID *big.Int, hash ethcommon.Hash) (*ethtypes.Header, error) {
	var head *ethtypes.Header
	err := ec.getRPC(chainID).CallContext(ctx, &head, "eth_getBlockByHash", hash, false)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	if head == nil {
		return nil, errors.NotFoundError("not found").SetComponent(component)
	}

	return head, nil
}

// HeaderByNumber returns a block header from the current canonical chain. If number is
// nil, the latest known header is returned.
func (ec *Client) HeaderByNumber(ctx context.Context, chainID, number *big.Int) (*ethtypes.Header, error) {
	var head *ethtypes.Header
	err := ec.getRPC(chainID).CallContext(ctx, &head, "eth_getBlockByNumber", toBlockNumArg(number), false)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	if head == nil {
		return nil, errors.NotFoundError("not found").SetComponent(component)
	}

	return head, nil
}

// TransactionByHash returns the transaction with the given hash.
func (ec *Client) TransactionByHash(ctx context.Context, chainID *big.Int, hash ethcommon.Hash) (tx *ethtypes.Transaction, isPending bool, err error) {
	var raw json.RawMessage
	err = ec.getRPC(chainID).CallContext(ctx, &raw, "eth_getTransactionByHash", hash)
	if err != nil {
		return nil, false, errors.FromError(err).ExtendComponent(component)
	}
	if err := encoding.Unmarshal(raw, &tx); err != nil {
		return nil, false, errors.FromError(err).ExtendComponent(component)
	}

	if _, r, _ := tx.RawSignatureValues(); r == nil {
		return nil, false, errors.DataCorruptedError("transaction without signature").ExtendComponent(component)
	}

	// Unmarshal block body information
	var extra *txExtraInfo
	if err := encoding.Unmarshal(raw, &extra); err != nil {
		return nil, false, errors.FromError(err).ExtendComponent(component)
	}

	return tx, extra.BlockNumber == nil, nil
}

// TransactionReceipt returns the receipt of a transaction by transaction hash.
// Note that the receipt is not available for pending transactions.
func (ec *Client) TransactionReceipt(ctx context.Context, chainID *big.Int, txHash ethcommon.Hash) (*ethtypes.Receipt, error) {
	var r *ethtypes.Receipt
	err := ec.getRPC(chainID).CallContext(ctx, &r, "eth_getTransactionReceipt", txHash)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	if r == nil {
		return nil, errors.NotFoundError("not found").SetComponent(component)
	}

	return r, nil
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	return hexutil.EncodeBig(number)
}

// Networks return networks ID multi client is connected to
func (ec *Client) Networks(ctx context.Context) (networks []*big.Int) {
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
	var raw json.RawMessage
	if err := ec.getRPC(chainID).CallContext(ctx, &raw, "eth_syncing"); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	// Handle the possible response types
	var syncing bool
	if err := encoding.Unmarshal(raw, &syncing); err == nil {
		return nil, nil
	}

	var progress *Progress
	if err := encoding.Unmarshal(raw, &progress); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
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
	var result hexutil.Big
	err := ec.getRPC(chainID).CallContext(ctx, &result, "eth_getBalance", account, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return (*big.Int)(&result), nil
}

// StorageAt returns the value of key in the contract storage of the given account.
// The block number can be nil, in which case the value is taken from the latest known block.
func (ec *Client) StorageAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, key ethcommon.Hash, blockNumber *big.Int) ([]byte, error) {
	var result hexutil.Bytes
	err := ec.getRPC(chainID).CallContext(ctx, &result, "eth_getStorageAt", account, key, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return result, nil
}

// CodeAt returns the contract code of the given account.
// The block number can be nil, in which case the code is taken from the latest known block.
func (ec *Client) CodeAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, blockNumber *big.Int) ([]byte, error) {
	var result hexutil.Bytes
	err := ec.getRPC(chainID).CallContext(ctx, &result, "eth_getCode", account, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return result, nil
}

// NonceAt returns the account nonce of the given account.
// The block number can be nil, in which case the nonce is taken from the latest known block.
func (ec *Client) NonceAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, blockNumber *big.Int) (uint64, error) {
	var result hexutil.Uint64
	err := ec.getRPC(chainID).CallContext(ctx, &result, "eth_getTransactionCount", account, toBlockNumArg(blockNumber))
	if err != nil {
		return 0, errors.FromError(err).ExtendComponent(component)
	}
	return uint64(result), nil
}

// Pending State

// PendingBalanceAt returns the wei balance of the given account in the pending state.
func (ec *Client) PendingBalanceAt(ctx context.Context, chainID *big.Int, account ethcommon.Address) (*big.Int, error) {
	var result hexutil.Big
	err := ec.getRPC(chainID).CallContext(ctx, &result, "eth_getBalance", account, "pending")
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return (*big.Int)(&result), nil
}

// PendingStorageAt returns the value of key in the contract storage of the given account in the pending state.
func (ec *Client) PendingStorageAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, key ethcommon.Hash) ([]byte, error) {
	var result hexutil.Bytes
	err := ec.getRPC(chainID).CallContext(ctx, &result, "eth_getStorageAt", account, key, "pending")
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return result, nil
}

// PendingCodeAt returns the contract code of the given account in the pending state.
func (ec *Client) PendingCodeAt(ctx context.Context, chainID *big.Int, account ethcommon.Address) ([]byte, error) {
	var result hexutil.Bytes
	err := ec.getRPC(chainID).CallContext(ctx, &result, "eth_getCode", account, "pending")
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return result, nil
}

// PendingNonceAt returns the account nonce of the given account in the pending state.
// This is the nonce that should be used for the next transaction.
func (ec *Client) PendingNonceAt(ctx context.Context, chainID *big.Int, account ethcommon.Address) (uint64, error) {
	var result hexutil.Uint64
	err := ec.getRPC(chainID).CallContext(ctx, &result, "eth_getTransactionCount", account, "pending")
	if err != nil {
		return 0, errors.FromError(err).ExtendComponent(component)
	}
	return uint64(result), nil
}

// Contract Calling

// CallContract executes a message call transaction, which is directly executed in the VM
// of the node, but never mined into the blockchain.
//
// blockNumber selects the block height at which the call runs. It can be nil, in which
// case the code is taken from the latest known block. Note that state from very old
// blocks might not be available.
func (ec *Client) CallContract(ctx context.Context, chainID *big.Int, msg *eth.CallMsg, blockNumber *big.Int) ([]byte, error) {
	var hex hexutil.Bytes
	err := ec.getRPC(chainID).CallContext(ctx, &hex, "eth_call", toCallArg(msg), toBlockNumArg(blockNumber))
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return hex, nil
}

// PendingCallContract executes a message call transaction using the EVM.
// The state seen by the contract call is the pending state.
func (ec *Client) PendingCallContract(ctx context.Context, chainID *big.Int, msg *eth.CallMsg) ([]byte, error) {
	var hex hexutil.Bytes
	err := ec.getRPC(chainID).CallContext(ctx, &hex, "eth_call", toCallArg(msg), "pending")
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return hex, nil
}

// SuggestGasPrice retrieves the currently suggested gas price to allow a timely
// execution of a transaction.
func (ec *Client) SuggestGasPrice(ctx context.Context, chainID *big.Int) (*big.Int, error) {
	var hex hexutil.Big
	err := ec.getRPC(chainID).CallContext(ctx, &hex, "eth_gasPrice")
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return (*big.Int)(&hex), nil
}

// EstimateGas tries to estimate the gas needed to execute a specific transaction based on
// the current pending state of the backend blockchain. There is no guarantee that this is
// the true gas limit requirement as other transactions may be added or removed by miners,
// but it should provide a basis for setting a reasonable default.
func (ec *Client) EstimateGas(ctx context.Context, chainID *big.Int, msg *eth.CallMsg) (uint64, error) {
	var hex hexutil.Uint64
	err := ec.getRPC(chainID).CallContext(ctx, &hex, "eth_estimateGas", toCallArg(msg))
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
func (ec *Client) SendRawTransaction(ctx context.Context, chainID *big.Int, raw string) error {
	err := ec.getRPC(chainID).CallContext(ctx, nil, "eth_sendRawTransaction", raw)
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}
	return nil
}

// SendTransaction send transaction to an Ethereum node
func (ec *Client) SendTransaction(ctx context.Context, chainID *big.Int, args *types.SendTxArgs) (txHash ethcommon.Hash, err error) {
	log.WithFields(log.Fields{
		"nonce": args.Nonce.String(),
		"from":  args.From.Hex(),
	}).Info("sending a transaction")

	err = ec.getRPC(chainID).CallContext(ctx, &txHash, "eth_sendTransaction", &args)
	if err != nil {
		return ethcommon.Hash{}, errors.FromError(err).ExtendComponent(component)
	}
	return txHash, nil
}

func (ec *Client) SendQuorumRawPrivateTransaction(ctx context.Context, chainID *big.Int, signedTxHash []byte, privateFor []string) (ethcommon.Hash, error) {
	rawTxHashHex := hexutil.Encode(signedTxHash)
	privateForParam := map[string]interface{}{
		"privateFor": privateFor,
	}
	var hash string
	err := ec.getRPC(chainID).CallContext(ctx, &hash, "eth_sendRawPrivateTransaction", rawTxHashHex, privateForParam)
	if err != nil {
		return ethcommon.Hash{}, errors.FromError(err).ExtendComponent(component)
	}
	return ethcommon.HexToHash(hash), nil
}

// SendRawPrivateTransaction send a raw transaction to an Ethereum node supporting EEA extension
func (ec *Client) SendRawPrivateTransaction(ctx context.Context, chainID *big.Int, raw []byte, args *types.PrivateArgs) (ethcommon.Hash, error) {
	// Send a raw signed transactions using EEA extension method
	// Method documentation here: https://docs.pantheon.pegasys.tech/en/latest/Reference/Pantheon-API-Methods/#eea_sendrawtransaction
	var hash string
	err := ec.getRPC(chainID).CallContext(ctx, &hash, "eea_sendRawTransaction", hexutil.Encode(raw))
	if err != nil {
		return ethcommon.Hash{}, errors.FromError(err).ExtendComponent(component)
	}
	return ethcommon.HexToHash(hash), nil
}
