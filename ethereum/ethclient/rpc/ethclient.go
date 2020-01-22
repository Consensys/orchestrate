package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"math/big"
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
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/types"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

type ProcessResultFunc func(result json.RawMessage) error

// Client is a connector to Ethereum blockchains that uses Geth rpc client
type Client struct {
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

// NewClient creates a new MultiClient
func NewClient(conf *Config, client *http.Client) *Client {
	return &Client{
		client: client,
		conf:   conf,
		pool: &sync.Pool{
			New: func() interface{} { return NewBackOff(conf) },
		},
		idCounter: 0,
	}
}

func (ec *Client) Call(ctx context.Context, endpoint string, processResult func(result json.RawMessage) error, method string, args ...interface{}) error {
	req, err := ec.newJSONRpcRequestWithContext(ctx, endpoint, method, args...)
	if err != nil {
		return err
	}

	bckoff := backoff.WithContext(ec.pool.Get().(backoff.BackOff), req.Context())
	defer ec.pool.Put(bckoff)

	return ec.callWithRetry(req, processResult, bckoff)
}

func (ec *Client) callWithRetry(req *http.Request, processResult func(result json.RawMessage) error, bckoff backoff.BackOff) error {
	return backoff.RetryNotify(
		func() error {
			e := ec.call(req, processResult)
			switch {
			case e == nil:
				return nil
			case errors.IsConnectionError(e),
				errors.IsNotFoundError(e) && utils.ShouldRetryNotFoundError(req.Context()):
				return e
			default:
				return backoff.Permanent(e)
			}
		},
		bckoff,
		func(e error, duration time.Duration) {
			// Reset body before retrying
			req.Body, _ = req.GetBody()
			log.FromContext(req.Context()).
				WithError(e).
				Warnf("eth-client: JSON-RPC call failed, retrying in %v...", duration)
		},
	)
}

func (ec *Client) call(req *http.Request, processResult ProcessResultFunc) error {
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
		return processResult(respMsg.Result)
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

func (ec *Client) newJSONRpcMessage(method string, args ...interface{}) (*JSONRpcMessage, error) {
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

func (ec *Client) newJSONRpcRequestWithContext(ctx context.Context, endpoint, method string, args ...interface{}) (*http.Request, error) {
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

func (ec *Client) do(req *http.Request) (*http.Response, error) {
	resp, err := ec.client.Do(req)
	if err != nil {
		return nil, errors.EthConnectionError(err.Error())
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.EthConnectionError("%v (code=%v)", resp.Status, resp.StatusCode)
	}

	return resp, nil
}

func (ec *Client) nextID() json.RawMessage {
	id := atomic.AddUint32(&ec.idCounter, 1)
	return strconv.AppendUint(nil, uint64(id), 10)
}

func (ec *Client) processEthError(err *JSONError) error {
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
	Hash         ethcommon.Hash          `json:"hash"`
	Transactions []*ethtypes.Transaction `json:"transactions"`
	UncleHashes  []ethcommon.Hash        `json:"uncles"`
}

func processResult(v interface{}) ProcessResultFunc {
	return func(result json.RawMessage) error {
		err := json.Unmarshal(result, &v)
		if err != nil {
			return errors.EncodingError(err.Error())
		}

		return nil
	}
}

func processBlockResult(header **ethtypes.Header, body **Body) ProcessResultFunc {
	return func(result json.RawMessage) error {
		var raw json.RawMessage
		err := processResult(&raw)(result)
		if err != nil {
			return err
		}

		if len(raw) == 0 {
			// Block was not found
			return errors.NotFoundError("block not found")
		}

		// Unmarshal block header information
		if err := encoding.Unmarshal(raw, header); err != nil {
			return errors.FromError(err).ExtendComponent(component)
		}

		// Unmarshal block body information
		if err := encoding.Unmarshal(raw, body); err != nil {
			return errors.FromError(err).ExtendComponent(component)
		}

		return nil
	}
}

// BlockByHash returns the given full block.
func (ec *Client) BlockByHash(ctx context.Context, endpoint string, hash ethcommon.Hash) (*ethtypes.Block, error) {
	// Perform RPC call
	var header *ethtypes.Header
	var body *Body
	err := ec.Call(ctx, endpoint, processBlockResult(&header, &body), "eth_getBlockByHash", hash, true)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return ethtypes.NewBlock(header, body.Transactions, []*ethtypes.Header{}, []*ethtypes.Receipt{}), nil
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

	return ethtypes.NewBlock(header, body.Transactions, []*ethtypes.Header{}, []*ethtypes.Receipt{}), nil
}

func processHeaderResult(head **ethtypes.Header) ProcessResultFunc {
	return func(result json.RawMessage) error {
		err := processResult(head)(result)
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

func processTxResult(tx **ethtypes.Transaction, extra **txExtraInfo) ProcessResultFunc {
	return func(result json.RawMessage) error {
		var raw json.RawMessage
		err := processResult(&raw)(result)
		if err != nil {
			return err
		}

		if len(raw) == 0 {
			// Block was not found
			return errors.NotFoundError("transaction not found")
		}

		if err := encoding.Unmarshal(raw, tx); err != nil {
			return errors.FromError(err)
		}

		if _, r, _ := (*tx).RawSignatureValues(); r == nil {
			return errors.DataCorruptedError("transaction without signature")
		}

		// Unmarshal block body information
		if err := encoding.Unmarshal(raw, extra); err != nil {
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

func processReceiptResult(receipt **ethtypes.Receipt) ProcessResultFunc {
	return func(result json.RawMessage) error {
		err := processResult(&receipt)(result)
		if err != nil {
			return err
		}

		if *receipt == nil {
			// Receipt was not found
			return errors.NotFoundError("receipt not found")
		}

		return nil
	}
}

// TransactionReceipt returns the receipt of a transaction by transaction hash.
// Note that the receipt is not available for pending transactions.
func (ec *Client) TransactionReceipt(ctx context.Context, endpoint string, txHash ethcommon.Hash) (*ethtypes.Receipt, error) {
	var r *ethtypes.Receipt
	err := ec.Call(ctx, endpoint, processReceiptResult(&r), "eth_getTransactionReceipt", txHash)
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

func processProgressResult(progress **Progress) ProcessResultFunc {
	return func(result json.RawMessage) error {
		var raw json.RawMessage
		err := processResult(&raw)(result)
		if err != nil {
			return err
		}

		var syncing bool
		if err = encoding.Unmarshal(raw, &syncing); err == nil {
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
	if err := ec.Call(ctx, endpoint, processResult(&version), "net_version"); err != nil {
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
	err := ec.Call(ctx, endpoint, processResult(&balance), "eth_getBalance", account, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return (*big.Int)(&balance), nil
}

// StorageAt returns the value of key in the contract storage of the given account.
// The block number can be nil, in which case the value is taken from the latest known block.
func (ec *Client) StorageAt(ctx context.Context, endpoint string, account ethcommon.Address, key ethcommon.Hash, blockNumber *big.Int) ([]byte, error) {
	var storage hexutil.Bytes
	err := ec.Call(ctx, endpoint, processResult(&storage), "eth_getStorageAt", account, key, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return storage, nil
}

// CodeAt returns the contract code of the given account.
// The block number can be nil, in which case the code is taken from the latest known block.
func (ec *Client) CodeAt(ctx context.Context, endpoint string, account ethcommon.Address, blockNumber *big.Int) ([]byte, error) {
	var code hexutil.Bytes
	err := ec.Call(ctx, endpoint, processResult(&code), "eth_getCode", account, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return code, nil
}

// NonceAt returns the account nonce of the given account.
// The block number can be nil, in which case the nonce is taken from the latest known block.
func (ec *Client) NonceAt(ctx context.Context, endpoint string, account ethcommon.Address, blockNumber *big.Int) (uint64, error) {
	var nonce hexutil.Uint64
	err := ec.Call(ctx, endpoint, processResult(&nonce), "eth_getTransactionCount", account, toBlockNumArg(blockNumber))
	if err != nil {
		return 0, errors.FromError(err).ExtendComponent(component)
	}
	return uint64(nonce), nil
}

// Pending State

// PendingBalanceAt returns the wei balance of the given account in the pending state.
func (ec *Client) PendingBalanceAt(ctx context.Context, endpoint string, account ethcommon.Address) (*big.Int, error) {
	var balance hexutil.Big
	err := ec.Call(ctx, endpoint, processResult(&balance), "eth_getBalance", account, "pending")
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return (*big.Int)(&balance), nil
}

// PendingStorageAt returns the value of key in the contract storage of the given account in the pending state.
func (ec *Client) PendingStorageAt(ctx context.Context, endpoint string, account ethcommon.Address, key ethcommon.Hash) ([]byte, error) {
	var storage hexutil.Bytes
	err := ec.Call(ctx, endpoint, processResult(&storage), "eth_getStorageAt", account, key, "pending")
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return storage, nil
}

// PendingCodeAt returns the contract code of the given account in the pending state.
func (ec *Client) PendingCodeAt(ctx context.Context, endpoint string, account ethcommon.Address) ([]byte, error) {
	var code hexutil.Bytes
	err := ec.Call(ctx, endpoint, processResult(&code), "eth_getCode", account, "pending")
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return code, nil
}

// PendingNonceAt returns the account nonce of the given account in the pending state.
// This is the nonce that should be used for the next transaction.
func (ec *Client) PendingNonceAt(ctx context.Context, endpoint string, account ethcommon.Address) (uint64, error) {
	var nonce hexutil.Uint64
	err := ec.Call(ctx, endpoint, processResult(&nonce), "eth_getTransactionCount", account, "pending")
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
	err := ec.Call(ctx, endpoint, processResult(&hex), "eth_call", toCallArg(msg), toBlockNumArg(blockNumber))
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return hex, nil
}

// PendingCallContract executes a message call transaction using the EVM.
// The state seen by the contract call is the pending state.
func (ec *Client) PendingCallContract(ctx context.Context, endpoint string, msg *eth.CallMsg) ([]byte, error) {
	var hex hexutil.Bytes
	err := ec.Call(ctx, endpoint, processResult(&hex), "eth_call", toCallArg(msg), "pending")
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return hex, nil
}

// SuggestGasPrice retrieves the currently suggested gas price to allow a timely
// execution of a transaction.
func (ec *Client) SuggestGasPrice(ctx context.Context, endpoint string) (*big.Int, error) {
	var hex hexutil.Big
	err := ec.Call(ctx, endpoint, processResult(&hex), "eth_gasPrice")
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return (*big.Int)(&hex), nil
}

// EstimateGas tries to estimate the gas needed to execute a specific transaction based on
// the current pending state of the backend blockchain. There is no guarantee that this is
// the true gas limit requirement as other transactions may be added or removed by miners,
// but it should provide a basis for setting a reasonable default.
func (ec *Client) EstimateGas(ctx context.Context, endpoint string, msg *eth.CallMsg) (uint64, error) {
	var hex hexutil.Uint64
	err := ec.Call(ctx, endpoint, processResult(&hex), "eth_estimateGas", toCallArg(msg))
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
func (ec *Client) SendRawTransaction(ctx context.Context, endpoint, raw string) error {
	err := ec.Call(ctx, endpoint, processResult(nil), "eth_sendRawTransaction", raw)
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}
	return nil
}

// SendTransaction send transaction to an Ethereum node
func (ec *Client) SendTransaction(ctx context.Context, endpoint string, args *types.SendTxArgs) (txHash ethcommon.Hash, err error) {
	err = ec.Call(ctx, endpoint, processResult(&txHash), "eth_sendTransaction", &args)
	if err != nil {
		return ethcommon.Hash{}, errors.FromError(err).ExtendComponent(component)
	}
	return txHash, nil
}

func (ec *Client) SendQuorumRawPrivateTransaction(ctx context.Context, endpoint string, signedTxHash []byte, privateFor []string) (ethcommon.Hash, error) {
	rawTxHashHex := hexutil.Encode(signedTxHash)
	privateForParam := map[string]interface{}{
		"privateFor": privateFor,
	}
	var hash string
	err := ec.Call(ctx, endpoint, processResult(&hash), "eth_sendRawPrivateTransaction", rawTxHashHex, privateForParam)
	if err != nil {
		return ethcommon.Hash{}, errors.FromError(err).ExtendComponent(component)
	}
	return ethcommon.HexToHash(hash), nil
}

// SendRawPrivateTransaction send a raw transaction to an Ethereum node supporting EEA extension
func (ec *Client) SendRawPrivateTransaction(ctx context.Context, endpoint string, raw []byte, args *types.PrivateArgs) (ethcommon.Hash, error) {
	// Send a raw signed transactions using EEA extension method
	// Method documentation here: https://besu.hyperledger.org/en/stable/Reference/API-Methods/#eea_sendrawtransaction
	var hash string
	err := ec.Call(ctx, endpoint, processResult(&hash), "eea_sendRawTransaction", hexutil.Encode(raw))
	if err != nil {
		return ethcommon.Hash{}, errors.FromError(err).ExtendComponent(component)
	}
	return ethcommon.HexToHash(hash), nil
}
