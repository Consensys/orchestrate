package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/containous/traefik/v2/pkg/log"
	eth "github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/types"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

// ClientV2 is a connector to Ethereum blockchains that uses Geth rpc client
type ClientV2 struct {
	client *http.Client

	// Pool for backoffs
	pool *sync.Pool

	conf *Config

	idCounter uint32
}

// NewBackOff creates a new Exponential backoff
func NewBackOff(conf *Config) backoff.BackOff {
	return &backoff.ExponentialBackOff{
		InitialInterval:     conf.Retry.InitialInterval,
		RandomizationFactor: conf.Retry.RandomizationFactor,
		Multiplier:          conf.Retry.Multiplier,
		MaxInterval:         conf.Retry.MaxInterval,
		MaxElapsedTime:      conf.Retry.MaxElapsedTime,
		Clock:               backoff.SystemClock,
	}
}

// NewClientV2 creates a new MultiClient
func NewClientV2(conf *Config) *ClientV2 {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialer.DialContext,
			MaxIdleConnsPerHost:   200,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	return &ClientV2{
		client: client,
		conf:   conf,
		pool: &sync.Pool{
			New: func() interface{} { return NewBackOff(conf) },
		},
		idCounter: 0,
	}
}

func (ec *ClientV2) Call(ctx context.Context, endpoint string, result interface{}, method string, args ...interface{}) error {
	req, err := ec.newJSONRpcRequestWithContext(ctx, endpoint, method, args...)
	if err != nil {
		return err
	}

	resp, err := ec.do(req)
	if err != nil {
		return err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	var respMsg JSONRpcMessage
	if err := json.NewDecoder(resp.Body).Decode(&respMsg); err != nil {
		return errors.EncodingError(err.Error())
	}

	switch {
	case respMsg.Error != nil:
		return ec.processEthError(respMsg.Error)
	case len(respMsg.Result) == 0:
		return errors.NotFoundError("not found")
	default:
		err := json.Unmarshal(respMsg.Result, &result)
		if err != nil {
			return errors.EncodingError(err.Error())
		}
		return nil
	}
}

type JSONRpcMessage struct {
	Version string          `json:"jsonrpc,omitempty"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Error   *JSONError      `json:"error,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
}

type JSONError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (ec *ClientV2) newJSONRpcMessage(method string, args ...interface{}) (*JSONRpcMessage, error) {
	msg := &JSONRpcMessage{
		Method:  method,
		Version: "2.0",
		ID:      ec.nextID(),
	}
	if args != nil {
		var err error
		if msg.Params, err = json.Marshal(args); err != nil {
			return nil, errors.EncodingError(err.Error())
		}
	}
	return msg, nil
}

func (ec *ClientV2) newJSONRpcRequestWithContext(ctx context.Context, endpoint, method string, args ...interface{}) (*http.Request, error) {
	// Create RPC message
	msg, err := ec.newJSONRpcMessage(method, args...)
	if err != nil {
		return nil, err
	}

	// Marshal body
	body, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// Set headers for JSON-RPC request
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (ec *ClientV2) do(req *http.Request) (resp *http.Response, err error) {
	bckoff := backoff.WithContext(ec.pool.Get().(backoff.BackOff), req.Context())
	defer ec.pool.Put(bckoff)

	err = backoff.RetryNotify(
		func() error {
			resp, err = ec.client.Do(req)
			switch {
			case err != nil:
				return err
			case resp.StatusCode < 200, resp.StatusCode >= 300:
				return fmt.Errorf("%v (code=%v)", resp.Status, resp.StatusCode)
			default:
				return nil
			}
		},
		bckoff,
		func(err error, duration time.Duration) {
			log.FromContext(req.Context()).
				WithError(err).
				WithFields(logrus.Fields{
					"json-rpc.url": req.URL.String(),
				}).Warnf("eth-client: JSON-RPC call failed, retrying in %v...", duration)
		},
	)
	if err != nil {
		return nil, errors.EthConnectionError(err.Error())
	}
	return resp, nil
}

func (ec *ClientV2) nextID() json.RawMessage {
	id := atomic.AddUint32(&ec.idCounter, 1)
	return strconv.AppendUint(nil, uint64(id), 10)
}

func (ec *ClientV2) processEthError(err *JSONError) error {
	if strings.Contains(err.Message, "nonce too low") || strings.Contains(err.Message, "Nonce too low") {
		return errors.NonceTooLowError(err.Message)
	}
	return errors.EthereumError(err.Message)
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
func (ec *ClientV2) BlockByHash(ctx context.Context, endpoint string, hash ethcommon.Hash) (*ethtypes.Block, error) {
	// Perform RPC call
	var raw json.RawMessage
	err := ec.Call(ctx, endpoint, &raw, "eth_getBlockByHash", hash, true)
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
func (ec *ClientV2) BlockByNumber(ctx context.Context, endpoint string, number *big.Int) (*ethtypes.Block, error) {
	// Perform RPC call
	var raw json.RawMessage
	err := ec.Call(ctx, endpoint, &raw, "eth_getBlockByNumber", toBlockNumArg(number), true)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return blockFromRaw(raw)
}

// HeaderByHash returns the block header with the given hash.
func (ec *ClientV2) HeaderByHash(ctx context.Context, endpoint string, hash ethcommon.Hash) (*ethtypes.Header, error) {
	var head *ethtypes.Header
	err := ec.Call(ctx, endpoint, &head, "eth_getBlockByHash", hash, false)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return head, nil
}

// HeaderByNumber returns a block header from the current canonical chain. If number is
// nil, the latest known header is returned.
func (ec *ClientV2) HeaderByNumber(ctx context.Context, endpoint string, number *big.Int) (*ethtypes.Header, error) {
	var head *ethtypes.Header
	err := ec.Call(ctx, endpoint, &head, "eth_getBlockByNumber", toBlockNumArg(number), false)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return head, nil
}

// TransactionByHash returns the transaction with the given hash.
func (ec *ClientV2) TransactionByHash(ctx context.Context, endpoint string, hash ethcommon.Hash) (tx *ethtypes.Transaction, isPending bool, err error) {
	var raw json.RawMessage
	err = ec.Call(ctx, endpoint, &raw, "eth_getTransactionByHash", hash)
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
func (ec *ClientV2) TransactionReceipt(ctx context.Context, endpoint string, txHash ethcommon.Hash) (*ethtypes.Receipt, error) {
	var r *ethtypes.Receipt
	err := ec.Call(ctx, endpoint, &r, "eth_getTransactionReceipt", txHash)
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

// SyncProgress retrieves the current progress of the sync algorithm. If there's
// no sync currently running, it returns nil.
func (ec *ClientV2) SyncProgress(ctx context.Context, endpoint string) (*eth.SyncProgress, error) {
	var raw json.RawMessage
	if err := ec.Call(ctx, endpoint, &raw, "eth_syncing"); err != nil {
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
func (ec *ClientV2) BalanceAt(ctx context.Context, endpoint string, account ethcommon.Address, blockNumber *big.Int) (*big.Int, error) {
	var result hexutil.Big
	err := ec.Call(ctx, endpoint, &result, "eth_getBalance", account, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return (*big.Int)(&result), nil
}

// StorageAt returns the value of key in the contract storage of the given account.
// The block number can be nil, in which case the value is taken from the latest known block.
func (ec *ClientV2) StorageAt(ctx context.Context, endpoint string, account ethcommon.Address, key ethcommon.Hash, blockNumber *big.Int) ([]byte, error) {
	var result hexutil.Bytes
	err := ec.Call(ctx, endpoint, &result, "eth_getStorageAt", account, key, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return result, nil
}

// CodeAt returns the contract code of the given account.
// The block number can be nil, in which case the code is taken from the latest known block.
func (ec *ClientV2) CodeAt(ctx context.Context, endpoint string, account ethcommon.Address, blockNumber *big.Int) ([]byte, error) {
	var result hexutil.Bytes
	err := ec.Call(ctx, endpoint, &result, "eth_getCode", account, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return result, nil
}

// NonceAt returns the account nonce of the given account.
// The block number can be nil, in which case the nonce is taken from the latest known block.
func (ec *ClientV2) NonceAt(ctx context.Context, endpoint string, account ethcommon.Address, blockNumber *big.Int) (uint64, error) {
	var result hexutil.Uint64
	err := ec.Call(ctx, endpoint, &result, "eth_getTransactionCount", account, toBlockNumArg(blockNumber))
	if err != nil {
		return 0, errors.FromError(err).ExtendComponent(component)
	}
	return uint64(result), nil
}

// Pending State

// PendingBalanceAt returns the wei balance of the given account in the pending state.
func (ec *ClientV2) PendingBalanceAt(ctx context.Context, endpoint string, account ethcommon.Address) (*big.Int, error) {
	var result hexutil.Big
	err := ec.Call(ctx, endpoint, &result, "eth_getBalance", account, "pending")
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return (*big.Int)(&result), nil
}

// PendingStorageAt returns the value of key in the contract storage of the given account in the pending state.
func (ec *ClientV2) PendingStorageAt(ctx context.Context, endpoint string, account ethcommon.Address, key ethcommon.Hash) ([]byte, error) {
	var result hexutil.Bytes
	err := ec.Call(ctx, endpoint, &result, "eth_getStorageAt", account, key, "pending")
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return result, nil
}

// PendingCodeAt returns the contract code of the given account in the pending state.
func (ec *ClientV2) PendingCodeAt(ctx context.Context, endpoint string, account ethcommon.Address) ([]byte, error) {
	var result hexutil.Bytes
	err := ec.Call(ctx, endpoint, &result, "eth_getCode", account, "pending")
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return result, nil
}

// PendingNonceAt returns the account nonce of the given account in the pending state.
// This is the nonce that should be used for the next transaction.
func (ec *ClientV2) PendingNonceAt(ctx context.Context, endpoint string, account ethcommon.Address) (uint64, error) {
	var result hexutil.Uint64
	err := ec.Call(ctx, endpoint, &result, "eth_getTransactionCount", account, "pending")
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
func (ec *ClientV2) CallContract(ctx context.Context, endpoint string, msg *eth.CallMsg, blockNumber *big.Int) ([]byte, error) {
	var hex hexutil.Bytes
	err := ec.Call(ctx, endpoint, &hex, "eth_call", toCallArg(msg), toBlockNumArg(blockNumber))
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return hex, nil
}

// PendingCallContract executes a message call transaction using the EVM.
// The state seen by the contract call is the pending state.
func (ec *ClientV2) PendingCallContract(ctx context.Context, endpoint string, msg *eth.CallMsg) ([]byte, error) {
	var hex hexutil.Bytes
	err := ec.Call(ctx, endpoint, &hex, "eth_call", toCallArg(msg), "pending")
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return hex, nil
}

// SuggestGasPrice retrieves the currently suggested gas price to allow a timely
// execution of a transaction.
func (ec *ClientV2) SuggestGasPrice(ctx context.Context, endpoint string) (*big.Int, error) {
	var hex hexutil.Big
	err := ec.Call(ctx, endpoint, &hex, "eth_gasPrice")
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return (*big.Int)(&hex), nil
}

// EstimateGas tries to estimate the gas needed to execute a specific transaction based on
// the current pending state of the backend blockchain. There is no guarantee that this is
// the true gas limit requirement as other transactions may be added or removed by miners,
// but it should provide a basis for setting a reasonable default.
func (ec *ClientV2) EstimateGas(ctx context.Context, endpoint string, msg *eth.CallMsg) (uint64, error) {
	var hex hexutil.Uint64
	err := ec.Call(ctx, endpoint, &hex, "eth_estimateGas", toCallArg(msg))
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
func (ec *ClientV2) SendRawTransaction(ctx context.Context, endpoint, raw string) error {
	err := ec.Call(ctx, endpoint, nil, "eth_sendRawTransaction", raw)
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}
	return nil
}

// SendTransaction send transaction to an Ethereum node
func (ec *ClientV2) SendTransaction(ctx context.Context, endpoint string, args *types.SendTxArgs) (txHash ethcommon.Hash, err error) {
	err = ec.Call(ctx, endpoint, &txHash, "eth_sendTransaction", &args)
	if err != nil {
		return ethcommon.Hash{}, errors.FromError(err).ExtendComponent(component)
	}
	return txHash, nil
}

func (ec *ClientV2) SendQuorumRawPrivateTransaction(ctx context.Context, endpoint string, signedTxHash []byte, privateFor []string) (ethcommon.Hash, error) {
	rawTxHashHex := hexutil.Encode(signedTxHash)
	privateForParam := map[string]interface{}{
		"privateFor": privateFor,
	}
	var hash string
	err := ec.Call(ctx, endpoint, &hash, "eth_sendRawPrivateTransaction", rawTxHashHex, privateForParam)
	if err != nil {
		return ethcommon.Hash{}, errors.FromError(err).ExtendComponent(component)
	}
	return ethcommon.HexToHash(hash), nil
}

// SendRawPrivateTransaction send a raw transaction to an Ethereum node supporting EEA extension
func (ec *ClientV2) SendRawPrivateTransaction(ctx context.Context, endpoint string, raw []byte, args *types.PrivateArgs) (ethcommon.Hash, error) {
	// Send a raw signed transactions using EEA extension method
	// Method documentation here: https://besu.hyperledger.org/en/stable/Reference/API-Methods/#eea_sendrawtransaction
	var hash string
	err := ec.Call(ctx, endpoint, &hash, "eea_sendRawTransaction", hexutil.Encode(raw))
	if err != nil {
		return ethcommon.Hash{}, errors.FromError(err).ExtendComponent(component)
	}
	return ethcommon.HexToHash(hash), nil
}
