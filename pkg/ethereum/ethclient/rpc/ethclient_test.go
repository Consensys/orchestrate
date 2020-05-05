// +build unit

package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	eth "github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	backoffmock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/backoff/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/types"
	proto "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/ethereum"
	pkgUtils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

var testConfig = &pkgUtils.Config{
	Retry: &pkgUtils.RetryConfig{
		InitialInterval:     time.Millisecond,
		RandomizationFactor: 0.5,
		Multiplier:          1.5,
		MaxInterval:         time.Millisecond,
		MaxElapsedTime:      time.Millisecond,
	},
}

type mockRoundTripper struct{}

var skipPreCallRoundTrip bool

func (rt mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	if preCtx, ok := ctx.Value(testCtxKey("pre_call")).(context.Context); ok && !skipPreCallRoundTrip {
		skipPreCallRoundTrip = true
		ctx = preCtx
	}

	if err, ok := ctx.Value(testCtxKey("resp.error")).(error); ok {
		return nil, err
	}

	resp := &http.Response{}
	if statusCode, ok := ctx.Value(testCtxKey("resp.statusCode")).(int); ok {
		resp.StatusCode = statusCode
		resp.Status = http.StatusText(statusCode)
	}

	if body, ok := ctx.Value(testCtxKey("resp.body")).(io.ReadCloser); ok {
		resp.Body = body
	}

	return resp, nil
}

func newClient() *Client {
	newBackOff := func() backoff.BackOff { return pkgUtils.NewBackOff(testConfig) }
	ec := NewClient(newBackOff, &http.Client{
		Transport: mockRoundTripper{},
	})
	return ec
}

func TestProcessEthError(t *testing.T) {
	ec := newClient()

	// Nonce too low
	err := ec.processEthError(&JSONError{Message: "json-rpc: nonce too low"})
	assert.Equal(t, "BE001", errors.FromError(err).Hex(), "Error code should be correst")

	// Default
	err = ec.processEthError(&JSONError{Message: "json-rpc: failed"})
	assert.Equal(t, "BE000", errors.FromError(err).Hex(), "Error code should be correst")
}

type testCtxKey string

func newContext(err error, statusCode int, body io.ReadCloser) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, testCtxKey("resp.error"), err)
	ctx = context.WithValue(ctx, testCtxKey("resp.statusCode"), statusCode)
	ctx = context.WithValue(ctx, testCtxKey("resp.body"), body)
	return ctx
}

func TestDo(t *testing.T) {
	ec := newClient()

	// Test 1: with error
	ctx := newContext(fmt.Errorf("test-error"), 0, nil)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "test-endpoint", nil)
	_, err := ec.do(req)

	assert.Error(t, err, "#1 do should return error")
	assert.Contains(t, err.Error(), "test-error", "#1 Error message should be correct")

	// Test 2: with Status Code 201
	ctx = newContext(nil, 201, nil)
	req, _ = http.NewRequestWithContext(ctx, http.MethodPost, "test-endpoint", nil)
	_, err = ec.do(req)

	assert.NoError(t, err, "#3 do should not return error")
}

func makeRespBody(result interface{}, errMsg string) io.ReadCloser {
	respMsg := &JSONRpcMessage{}
	if result != nil {
		if b, ok := result.([]byte); ok {
			respMsg.Result = json.RawMessage(b)
		} else {
			jsonResult, _ := json.Marshal(result)
			respMsg.Result = json.RawMessage(jsonResult)
		}
	}

	if errMsg != "" {
		respMsg.Error = &JSONError{
			Message: errMsg,
		}
	}
	b, _ := json.Marshal(respMsg)
	return ioutil.NopCloser(bytes.NewReader(b))
}

func TestCallWithRetry(t *testing.T) {
	ec := newClient()
	var raw json.RawMessage

	// Test 1: Connection error, should retry
	bckoff := &backoffmock.MockBackoff{}
	ctx := newContext(fmt.Errorf("test-error"), 503, nil)
	err := ec.callWithRetry(ctx, func(context.Context) (*http.Request, error) {
		return ec.newJSONRpcRequestWithContext(ctx, "test-endpoint", "test_method")
	}, processResult(&raw), bckoff)
	assert.Error(t, err, "#1 TestCallWithRetry should  error")
	assert.True(t, bckoff.HasRetried(), "#1 Should have retried")

	// Test 2: not found error, should retry
	bckoff = &backoffmock.MockBackoff{}
	ctx = newContext(nil, 404, makeRespBody([]byte{}, ""))
	ctx = utils.RetryNotFoundError(ctx, true)
	err = ec.callWithRetry(ctx, func(context.Context) (*http.Request, error) {
		return ec.newJSONRpcRequestWithContext(ctx, "test-endpoint", "test_method")
	}, processResult(&raw), bckoff)
	assert.Error(t, err, "#2 TestCallWithRetry should  error")
	assert.True(t, bckoff.HasRetried(), "#2 Should have retried")

	// Test 3: invalid response body, should not retry
	bckoff = &backoffmock.MockBackoff{}
	ctx = newContext(nil, 200, makeRespBody([]byte(`"%@`), ""))
	err = ec.callWithRetry(ctx, func(context.Context) (*http.Request, error) {
		return ec.newJSONRpcRequestWithContext(ctx, "test-endpoint", "test_method")
	}, processResult(&raw), bckoff)
	assert.Error(t, err, "#3 TestCallWithRetry should  error")
	assert.False(t, bckoff.HasRetried(), "#3 Should not have retried")

	// Test 4: invalid response body with error status, should retry
	bckoff = &backoffmock.MockBackoff{}
	ctx = newContext(nil, 400, makeRespBody([]byte(`"%@`), ""))
	err = ec.callWithRetry(ctx, func(context.Context) (*http.Request, error) {
		return ec.newJSONRpcRequestWithContext(ctx, "test-endpoint", "test_method")
	}, processResult(&raw), bckoff)
	assert.Error(t, err, "#4 TestCallWithRetry should  error")
	assert.True(t, bckoff.HasRetried(), "#4 Should have retried")
}

func TestBlockByHash(t *testing.T) {
	ec := newClient()

	// Test 1 with Error
	ctx := newContext(fmt.Errorf("test-error"), 0, nil)
	_, err := ec.BlockByHash(ctx, "test-endpoint", ethcommon.HexToHash(""))
	assert.Error(t, err, "#1 BlockByHash should  error")

	// Test 2: empty block response
	blockEnc := ethcommon.FromHex("f90260f901f9a083cafc574e1f51ba9dc0568fc617a08ea2429fb384059c972f13b19fa1c8dd55a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347948888f1f195afa192cfee860698584c030f4c9db1a0ef1552a40b7165c3cd773806b9e0c165b75356e0314bf0706f279c729f51e017a05fe50b260da6308036625b850b5d6ced6d0a9f814c0688bc91ffb7b7a3a54b67a0bc37d79753ad738a6dac4921e57392f145d8887476de3f783dfa7edae9283e52b90100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008302000001832fefd8825208845506eb0780a0bd4472abb6659ebe3ee06ee4d7b72a00a9f4d001caca51342001075469aff49888a13a5a8c8f2bb1c4f861f85f800a82c35094095e7baea6a6c7c4c2dfeb977efac326af552d870a801ba09bea4c4daac7c7c52e093e6a4c35dbbcf8856f1af7b059ba20253e70848d094fa08a8fae537ce25ed8cb5af9adac3f141af69bd515bd2ba031522df09b97dd72b1c0")
	var expectedBlock ethtypes.Block
	if err = rlp.DecodeBytes(blockEnc, &expectedBlock); err != nil {
		t.Fatal("decode error: ", err)
	}

	ctx = newContext(nil, 200, makeRespBody(*expectedBlock.Header(), ""))
	respBlock, err := ec.BlockByHash(ctx, "test-endpoint", ethcommon.HexToHash(""))
	assert.NoError(t, err, "#2 BlockByHash should not error")
	assert.Equal(t, expectedBlock.NumberU64(), respBlock.NumberU64(), "#2 BlockByHash block number should match")
	assert.Equal(t, expectedBlock.ParentHash().Hex(), respBlock.ParentHash().Hex(), "#2 BlockByHash parent hash should match")

	// Test 3: empty block response
	ctx = newContext(nil, 200, makeRespBody(nil, ""))
	_, err = ec.BlockByHash(ctx, "test-endpoint", ethcommon.HexToHash(""))
	assert.Error(t, err, "#3 BlockByHash should error")
	assert.True(t, errors.IsNotFoundError(err), "#3 BlockByHash error should be not found")

	// Test 4: null block response
	ctx = newContext(nil, 200, makeRespBody([]byte(`null`), ""))
	_, err = ec.BlockByHash(ctx, "test-endpoint", ethcommon.HexToHash(""))
	assert.Error(t, err, "#4 BlockByHash should error")
	assert.True(t, errors.IsNotFoundError(err), "#4 BlockByHash error should be not found")
}

func TestBlockByNumber(t *testing.T) {
	ec := newClient()

	// Test 1 with Error
	ctx := newContext(fmt.Errorf("test-error"), 0, nil)
	_, err := ec.BlockByNumber(ctx, "test-endpoint", nil)
	assert.Error(t, err, "#1 BlockByNumber should  error")

	// Test 2: empty block response
	blockEnc := ethcommon.FromHex("f90260f901f9a083cafc574e1f51ba9dc0568fc617a08ea2429fb384059c972f13b19fa1c8dd55a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347948888f1f195afa192cfee860698584c030f4c9db1a0ef1552a40b7165c3cd773806b9e0c165b75356e0314bf0706f279c729f51e017a05fe50b260da6308036625b850b5d6ced6d0a9f814c0688bc91ffb7b7a3a54b67a0bc37d79753ad738a6dac4921e57392f145d8887476de3f783dfa7edae9283e52b90100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008302000001832fefd8825208845506eb0780a0bd4472abb6659ebe3ee06ee4d7b72a00a9f4d001caca51342001075469aff49888a13a5a8c8f2bb1c4f861f85f800a82c35094095e7baea6a6c7c4c2dfeb977efac326af552d870a801ba09bea4c4daac7c7c52e093e6a4c35dbbcf8856f1af7b059ba20253e70848d094fa08a8fae537ce25ed8cb5af9adac3f141af69bd515bd2ba031522df09b97dd72b1c0")
	var expectedBlock ethtypes.Block
	if err = rlp.DecodeBytes(blockEnc, &expectedBlock); err != nil {
		t.Fatal("decode error: ", err)
	}

	ctx = newContext(nil, 200, makeRespBody(*expectedBlock.Header(), ""))
	respBlock, err := ec.BlockByNumber(ctx, "test-endpoint", nil)
	assert.NoError(t, err, "#2 BlockByNumber should not error")
	assert.Equal(t, expectedBlock.NumberU64(), respBlock.NumberU64(), "Block number should match")
	assert.Equal(t, expectedBlock.ParentHash().Hex(), respBlock.ParentHash().Hex(), "Parent hash should match")

	// Test 3: empty block response
	ctx = newContext(nil, 200, makeRespBody(nil, ""))
	_, err = ec.BlockByNumber(ctx, "test-endpoint", nil)
	assert.Error(t, err, "#3 BlockByNumber should error")
	assert.True(t, errors.IsNotFoundError(err), "#3 BlockByNumber error should be not found")

	// Test 4: null block response
	ctx = newContext(nil, 200, makeRespBody([]byte(`null`), ""))
	_, err = ec.BlockByNumber(ctx, "test-endpoint", nil)
	assert.Error(t, err, "#4 BlockByNumber should error")
	assert.True(t, errors.IsNotFoundError(err), "#4 BlockByNumber error should be not found")
}

func TestHeaderByHash(t *testing.T) {
	ec := newClient()

	// Test 1 with Error
	ctx := newContext(fmt.Errorf("test-error"), 0, nil)
	_, err := ec.HeaderByHash(ctx, "test-endpoint", ethcommon.HexToHash(""))
	assert.Error(t, err, "#1 HeaderByHash should  error")

	// Test 2: empty block response
	blockEnc := ethcommon.FromHex("f90260f901f9a083cafc574e1f51ba9dc0568fc617a08ea2429fb384059c972f13b19fa1c8dd55a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347948888f1f195afa192cfee860698584c030f4c9db1a0ef1552a40b7165c3cd773806b9e0c165b75356e0314bf0706f279c729f51e017a05fe50b260da6308036625b850b5d6ced6d0a9f814c0688bc91ffb7b7a3a54b67a0bc37d79753ad738a6dac4921e57392f145d8887476de3f783dfa7edae9283e52b90100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008302000001832fefd8825208845506eb0780a0bd4472abb6659ebe3ee06ee4d7b72a00a9f4d001caca51342001075469aff49888a13a5a8c8f2bb1c4f861f85f800a82c35094095e7baea6a6c7c4c2dfeb977efac326af552d870a801ba09bea4c4daac7c7c52e093e6a4c35dbbcf8856f1af7b059ba20253e70848d094fa08a8fae537ce25ed8cb5af9adac3f141af69bd515bd2ba031522df09b97dd72b1c0")
	var expectedBlock ethtypes.Block
	if err = rlp.DecodeBytes(blockEnc, &expectedBlock); err != nil {
		t.Fatal("decode error: ", err)
	}

	ctx = newContext(nil, 200, makeRespBody(*expectedBlock.Header(), ""))
	respHeader, err := ec.HeaderByHash(ctx, "test-endpoint", ethcommon.HexToHash(""))
	assert.NoError(t, err, "#2 HeaderByHash should not error")
	assert.Equal(t, expectedBlock.ParentHash().Hex(), respHeader.ParentHash.Hex(), "#2 HeaderByHash parent hash should match")

	// Test 3: empty response
	ctx = newContext(nil, 200, makeRespBody(nil, ""))
	_, err = ec.HeaderByHash(ctx, "test-endpoint", ethcommon.Hash{})
	assert.Error(t, err, "#3 HeaderByHash should error")
	assert.True(t, errors.IsNotFoundError(err), "#3 HeaderByHash error should be not found")

	// Test 4: null block response
	ctx = newContext(nil, 200, makeRespBody([]byte(`null`), ""))
	_, err = ec.HeaderByHash(ctx, "test-endpoint", ethcommon.Hash{})
	assert.Error(t, err, "#4 HeaderByHash should error")
	assert.True(t, errors.IsNotFoundError(err), "#4 HeaderByHash error should be not found")
}

func TestHeaderByNumber(t *testing.T) {
	ec := newClient()

	// Test 1 with Error
	ctx := newContext(fmt.Errorf("test-error"), 0, nil)
	_, err := ec.HeaderByNumber(ctx, "test-endpoint", nil)
	assert.Error(t, err, "#1 HeaderByNumber should  error")

	// Test 2: empty block response
	blockEnc := ethcommon.FromHex("f90260f901f9a083cafc574e1f51ba9dc0568fc617a08ea2429fb384059c972f13b19fa1c8dd55a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347948888f1f195afa192cfee860698584c030f4c9db1a0ef1552a40b7165c3cd773806b9e0c165b75356e0314bf0706f279c729f51e017a05fe50b260da6308036625b850b5d6ced6d0a9f814c0688bc91ffb7b7a3a54b67a0bc37d79753ad738a6dac4921e57392f145d8887476de3f783dfa7edae9283e52b90100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008302000001832fefd8825208845506eb0780a0bd4472abb6659ebe3ee06ee4d7b72a00a9f4d001caca51342001075469aff49888a13a5a8c8f2bb1c4f861f85f800a82c35094095e7baea6a6c7c4c2dfeb977efac326af552d870a801ba09bea4c4daac7c7c52e093e6a4c35dbbcf8856f1af7b059ba20253e70848d094fa08a8fae537ce25ed8cb5af9adac3f141af69bd515bd2ba031522df09b97dd72b1c0")
	var expectedBlock ethtypes.Block
	if err = rlp.DecodeBytes(blockEnc, &expectedBlock); err != nil {
		t.Fatal("decode error: ", err)
	}

	ctx = newContext(nil, 200, makeRespBody(*expectedBlock.Header(), ""))
	respHeader, err := ec.HeaderByNumber(ctx, "test-endpoint", nil)
	assert.NoError(t, err, "#2 HeaderByNumber should not error")
	assert.Equal(t, expectedBlock.ParentHash().Hex(), respHeader.ParentHash.Hex(), "#2 HeaderByHash parent hash should match")

	// Test 3: empty response
	ctx = newContext(nil, 200, makeRespBody(nil, ""))
	_, err = ec.HeaderByNumber(ctx, "test-endpoint", nil)
	assert.Error(t, err, "#3 HeaderByNumber should error")
	assert.True(t, errors.IsNotFoundError(err), "#3 HeaderByNumber error should be not found")

	// Test 4: null header response
	ctx = newContext(nil, 200, makeRespBody([]byte(`null`), ""))
	_, err = ec.HeaderByNumber(ctx, "test-endpoint", nil)
	assert.Error(t, err, "#4 HeaderByNumber should error")
	assert.True(t, errors.IsNotFoundError(err), "#4 HeaderByNumber error should be not found")
}

type transactionResp struct {
	Nonce    hexutil.Uint64     `json:"nonce"    gencodec:"required"`
	GasPrice *hexutil.Big       `json:"gasPrice" gencodec:"required"`
	Gas      hexutil.Uint64     `json:"gas"      gencodec:"required"`
	To       *ethcommon.Address `json:"to"       rlp:"nil"` // nil means contract creation
	Value    *hexutil.Big       `json:"value"    gencodec:"required"`
	Data     hexutil.Bytes      `json:"input"    gencodec:"required"`

	// Signature values
	V *hexutil.Big `json:"v" gencodec:"required"`
	R *hexutil.Big `json:"r" gencodec:"required"`
	S *hexutil.Big `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash        *ethcommon.Hash `json:"hash" rlp:"-"`
	BlockNumber *string         `json:"blockNumber,omitempty"`
}

func newTxResp(tx *ethtypes.Transaction, blockNumber string) *transactionResp {
	hash := tx.Hash()
	resp := &transactionResp{
		Nonce:    hexutil.Uint64(tx.Nonce()),
		GasPrice: (*hexutil.Big)(tx.GasPrice()),
		Gas:      hexutil.Uint64(tx.Gas()),
		To:       tx.To(),
		Value:    (*hexutil.Big)(tx.Value()),
		Data:     tx.Data(),
		Hash:     &hash,
	}

	v, r, s := tx.RawSignatureValues()
	resp.V = (*hexutil.Big)(v)
	resp.R = (*hexutil.Big)(r)
	resp.S = (*hexutil.Big)(s)

	if blockNumber != "" {
		resp.BlockNumber = &blockNumber
	}

	return resp
}

func TestTransactionByHash(t *testing.T) {
	ec := newClient()

	// Test 1 with Error
	ctx := newContext(fmt.Errorf("test-error"), 0, nil)
	_, _, err := ec.TransactionByHash(ctx, "test-endpoint", ethcommon.HexToHash(""))
	assert.Error(t, err, "#1 TransactionByHash should  error")

	expectedTx, _ := ethtypes.NewTransaction(
		3,
		ethcommon.HexToAddress("b94f5374fce5edbc8e2a8697c15331677e6ebf0b"),
		big.NewInt(10),
		2000,
		big.NewInt(1),
		ethcommon.FromHex("5544"),
	).WithSignature(
		ethtypes.HomesteadSigner{},
		ethcommon.Hex2Bytes("98ff921201554726367d2be8c804a7ff89ccf285ebc57dff8ae4c44b9c19ac4a8887321be575c8095f789dd4c743dfe42c1820f9231f98a962b210e3ac2452a301"),
	)

	// Test 2 transaction with invalid transaction and no block info
	txRest := newTxResp(expectedTx, "")
	ctx = newContext(nil, 200, makeRespBody(txRest, ""))
	tx, isPending, err := ec.TransactionByHash(ctx, "test-endpoint", ethcommon.HexToHash(""))
	assert.NoError(t, err, "#2 TransactionByHash should not  error")
	assert.True(t, isPending, "#2 TransactionByHash tx should be pending")
	assert.Equal(t, expectedTx.Hash().Hex(), tx.Hash().Hex(), "#2 TransactionByHash tx should have correct hash")

	// Test 3 transaction with invalid transaction and block info
	txRest = newTxResp(expectedTx, "0x9")
	ctx = newContext(nil, 200, makeRespBody(txRest, ""))
	tx, isPending, err = ec.TransactionByHash(ctx, "test-endpoint", ethcommon.HexToHash(""))
	assert.NoError(t, err, "#3 TransactionByHash should not  error")
	assert.False(t, isPending, "#3 TransactionByHash tx should not be pending")
	assert.Equal(t, expectedTx.Hash().Hex(), tx.Hash().Hex(), "#3 TransactionByHash tx should have correct hash")

	// Test 4: null tx response
	ctx = newContext(nil, 200, makeRespBody([]byte(`null`), ""))
	_, _, err = ec.TransactionByHash(ctx, "test-endpoint", ethcommon.HexToHash(""))
	assert.Error(t, err, "#4 TransactionByHash should error")
	assert.True(t, errors.IsNotFoundError(err), "#4 TransactionByHash error should be not found")
}

type receiptResp struct {
	CumulativeGasUsed hexutil.Uint64  `json:"cumulativeGasUsed" gencodec:"required"`
	Bloom             ethtypes.Bloom  `json:"logsBloom"         gencodec:"required"`
	Logs              []*ethtypes.Log `json:"logs"              gencodec:"required"`
	TxHash            ethcommon.Hash  `json:"transactionHash" gencodec:"required"`
	GasUsed           hexutil.Uint64  `json:"gasUsed" gencodec:"required"`
}

func newReceiptResp(r *ethtypes.Receipt) *receiptResp {
	return &receiptResp{
		CumulativeGasUsed: hexutil.Uint64(r.CumulativeGasUsed),
		Bloom:             r.Bloom,
		Logs:              r.Logs,
		TxHash:            r.TxHash,
		GasUsed:           hexutil.Uint64(r.GasUsed),
	}
}

func TestTransactionReceipt(t *testing.T) {
	ec := newClient()

	// Test 1 with Error
	ctx := newContext(fmt.Errorf("test-error"), 0, nil)
	_, err := ec.TransactionReceipt(ctx, "test-endpoint", ethcommon.HexToHash(""))
	assert.Error(t, err, "#1 TestTransactionReceipt should  error")

	// Test 2 transaction with invalid transaction and no block info
	expectedReceipt := &ethtypes.Receipt{
		Status:            0,
		CumulativeGasUsed: 1000,
		Logs: []*ethtypes.Log{
			{
				Address: ethcommon.BytesToAddress([]byte{0x11}),
				Topics:  []ethcommon.Hash{ethcommon.HexToHash("dead"), ethcommon.HexToHash("beef")},
				Data:    []byte{0x01, 0x00, 0xff},
			},
			{
				Address: ethcommon.BytesToAddress([]byte{0x01, 0x11}),
				Topics:  []ethcommon.Hash{ethcommon.HexToHash("dead"), ethcommon.HexToHash("beef")},
				Data:    []byte{0x01, 0x00, 0xff},
			},
		},
		ContractAddress: ethcommon.BytesToAddress([]byte{0x01, 0x11, 0x11}),
		GasUsed:         111111,
	}
	expectedReceipt.Bloom = ethtypes.CreateBloom(ethtypes.Receipts{expectedReceipt})

	receiptResp := newReceiptResp(expectedReceipt)
	ctx = newContext(nil, 200, makeRespBody(receiptResp, ""))
	receipt, err := ec.TransactionReceipt(ctx, "test-endpoint", ethcommon.HexToHash(""))
	assert.NoError(t, err, "#2 TransactionReceipt should not  error")
	assert.Equal(t, expectedReceipt.CumulativeGasUsed, receipt.CumulativeGasUsed, "#2 TransactionReceipt receipt should have correct cumulative gas used")

	// Test 3: null receipt response
	ctx = newContext(nil, 200, makeRespBody([]byte(`null`), ""))
	_, err = ec.TransactionReceipt(ctx, "test-endpoint", ethcommon.HexToHash(""))
	assert.Error(t, err, "#4 TransactionReceipt should error")
	assert.True(t, errors.IsNotFoundError(err), "#4 TransactionReceipt error should be not found")
}

func TestPrivateTransactionReceipt(t *testing.T) {
	ec := newClient()

	ethReceipt := &ethtypes.Receipt{
		Status:            0,
		CumulativeGasUsed: 1000,
		Logs: []*ethtypes.Log{
			{
				Address: ethcommon.BytesToAddress([]byte{0x11}),
				Topics:  []ethcommon.Hash{ethcommon.HexToHash("dead"), ethcommon.HexToHash("beef")},
				Data:    []byte{0x01, 0x00, 0xff},
			},
			{
				Address: ethcommon.BytesToAddress([]byte{0x01, 0x11}),
				Topics:  []ethcommon.Hash{ethcommon.HexToHash("dead"), ethcommon.HexToHash("beef")},
				Data:    []byte{0x01, 0x00, 0xff},
			},
		},
		ContractAddress: ethcommon.BytesToAddress([]byte{0x01, 0x11, 0x11}),
		GasUsed:         111111,
	}
	ethReceipt.Bloom = ethtypes.CreateBloom(ethtypes.Receipts{ethReceipt})

	privReceipt := &privateReceipt{
		Status: "0x1",
		Output: "0x12123",
		Logs: []*proto.Log{
			{
				Address: ethcommon.BytesToAddress([]byte{0x11}).String(),
				Topics:  []string{ethcommon.HexToHash("0x12123").String(), ethcommon.HexToHash("0x12123").String()},
				Data:    string([]byte{0x01, 0x00, 0xff}),
			},
			{
				Address: ethcommon.BytesToAddress([]byte{0x01, 0x11}).String(),
				Topics:  []string{ethcommon.HexToHash("0x12123").String(), ethcommon.HexToHash("0x12123").String()},
				Data:    string([]byte{0x01, 0x00, 0xff}),
			},
		},
		PrivateFor:  []string{"PrivateFor"},
		PrivateFrom: "PrivateFrom",
	}

	ctx := newContext(nil, 200, makeRespBody(privReceipt, ""))
	// First tx receipt to fetch is the Public receipt
	ctx = context.WithValue(ctx, testCtxKey("pre_call"),
		newContext(nil, 200, makeRespBody(newReceiptResp(ethReceipt), "")))

	receipt, err := ec.PrivateTransactionReceipt(ctx, "test-endpoint", ethcommon.HexToHash(""))
	assert.NoError(t, err, "TransactionReceipt should not  error")
	assert.Equal(t, ethReceipt.CumulativeGasUsed, receipt.CumulativeGasUsed, "TransactionReceipt receipt should have correct cumulative gas used")
	assert.Equal(t, privReceipt.Output, receipt.Output, "TransactionReceipt receipt should have correct priv tx output")
	assert.Equal(t, uint64(0x1), receipt.Status, "TransactionReceipt receipt should have correct priv tx status")
}

func TestSyncProgress(t *testing.T) {
	ec := newClient()

	// Test 1 with Error
	ctx := newContext(fmt.Errorf("test-error"), 0, nil)
	_, err := ec.SyncProgress(ctx, "test-endpoint")
	assert.Error(t, err, "#1 SyncProgress should  error")

	// Test 2 with bool response
	ctx = newContext(nil, 200, makeRespBody(true, ""))
	prgrss, err := ec.SyncProgress(ctx, "test-endpoint")
	assert.NoError(t, err, "#2 SyncProgress should not error")
	assert.Nil(t, prgrss, "#2 SyncProgress progress should be nil")

	// Test 3 with sync progress response
	expectedProgress := &eth.SyncProgress{
		StartingBlock: 10000,
	}
	ctx = newContext(nil, 200, makeRespBody(&Progress{StartingBlock: hexutil.Uint64(expectedProgress.StartingBlock)}, ""))
	prgrss, err = ec.SyncProgress(ctx, "test-endpoint")
	assert.NoError(t, err, "#3 SyncProgress should not error")
	assert.Equal(t, expectedProgress.StartingBlock, prgrss.StartingBlock, "#3 SyncProgress should nbe correct")
}

func TestBalanceAt(t *testing.T) {
	ec := newClient()

	// Test 1 with Error
	ctx := newContext(fmt.Errorf("test-error"), 0, nil)
	_, err := ec.BalanceAt(ctx, "test-endpoint", ethcommon.Address{}, nil)
	assert.Error(t, err, "#1 BalanceAt should  error")

	// Test 2 without error
	expectedBalance := big.NewInt(1000)
	ctx = newContext(nil, 200, makeRespBody((*hexutil.Big)(expectedBalance), ""))
	balance, err := ec.BalanceAt(ctx, "test-endpoint", ethcommon.Address{}, nil)
	assert.NoError(t, err, "#3 BalanceAt should not error")
	assert.Equal(t, expectedBalance.Text(10), balance.Text(10), "#3 BalanceAt balance should be correct")
}

func TestStorageAt(t *testing.T) {
	ec := newClient()

	// Test 1 with Error
	ctx := newContext(fmt.Errorf("test-error"), 0, nil)
	_, err := ec.StorageAt(ctx, "test-endpoint", ethcommon.Address{}, ethcommon.Hash{}, nil)
	assert.Error(t, err, "#1 StorageAt should  error")

	// Test 2 without error
	expectedStorage := hexutil.MustDecode("0xabcd")
	ctx = newContext(nil, 200, makeRespBody(hexutil.Bytes(expectedStorage), ""))
	storage, err := ec.StorageAt(ctx, "test-endpoint", ethcommon.Address{}, ethcommon.Hash{}, nil)
	assert.NoError(t, err, "#3 StorageAt should not error")
	assert.Equal(t, hexutil.Encode(expectedStorage), hexutil.Encode(storage), "#3 StorageAt storage should be correct")
}

func TestCodeAt(t *testing.T) {
	ec := newClient()

	// Test 1 with Error
	ctx := newContext(fmt.Errorf("test-error"), 0, nil)
	_, err := ec.CodeAt(ctx, "test-endpoint", ethcommon.Address{}, nil)
	assert.Error(t, err, "#1 CodeAt should  error")

	// Test 2 without error
	expectedCode := hexutil.MustDecode("0xabcd")
	ctx = newContext(nil, 200, makeRespBody(hexutil.Bytes(expectedCode), ""))
	code, err := ec.CodeAt(ctx, "test-endpoint", ethcommon.Address{}, nil)
	assert.NoError(t, err, "#3 CodeAt should not error")
	assert.Equal(t, hexutil.Encode(expectedCode), hexutil.Encode(code), "#3 CodeAt code should be correct")
}

func TestNonceAt(t *testing.T) {
	ec := newClient()

	// Test 1 with Error
	ctx := newContext(fmt.Errorf("test-error"), 0, nil)
	_, err := ec.NonceAt(ctx, "test-endpoint", ethcommon.Address{}, nil)
	assert.Error(t, err, "#1 NonceAt should  error")

	// Test 2 without error
	expectedNonce := uint64(17)
	ctx = newContext(nil, 200, makeRespBody(hexutil.Uint64(expectedNonce), ""))
	nonce, err := ec.NonceAt(ctx, "test-endpoint", ethcommon.Address{}, nil)
	assert.NoError(t, err, "#3 NonceAt should not error")
	assert.Equal(t, expectedNonce, nonce, "#3 NonceAt nonce should be correct")
}

func TestPendingBalanceAt(t *testing.T) {
	ec := newClient()

	// Test 1 with Error
	ctx := newContext(fmt.Errorf("test-error"), 0, nil)
	_, err := ec.PendingBalanceAt(ctx, "test-endpoint", ethcommon.Address{})
	assert.Error(t, err, "#1 TestPendingBalanceAt should  error")

	// Test 2 without error
	expectedBalance := big.NewInt(1000)
	ctx = newContext(nil, 200, makeRespBody((*hexutil.Big)(expectedBalance), ""))
	balance, err := ec.PendingBalanceAt(ctx, "test-endpoint", ethcommon.Address{})
	assert.NoError(t, err, "#3 TestPendingBalanceAt should not error")
	assert.Equal(t, expectedBalance.Text(10), balance.Text(10), "#3 TestPendingBalanceAt balance should be correct")
}

func TestPendingStorageAt(t *testing.T) {
	ec := newClient()

	// Test 1 with Error
	ctx := newContext(fmt.Errorf("test-error"), 0, nil)
	_, err := ec.PendingStorageAt(ctx, "test-endpoint", ethcommon.Address{}, ethcommon.Hash{})
	assert.Error(t, err, "#1 PendingStorageAt should  error")

	// Test 2 without error
	expectedStorage := hexutil.MustDecode("0xabcd")
	ctx = newContext(nil, 200, makeRespBody(hexutil.Bytes(expectedStorage), ""))
	storage, err := ec.PendingStorageAt(ctx, "test-endpoint", ethcommon.Address{}, ethcommon.Hash{})
	assert.NoError(t, err, "#3 PendingStorageAt should not error")
	assert.Equal(t, hexutil.Encode(expectedStorage), hexutil.Encode(storage), "#3 PendingStorageAt storage should be correct")
}

func TestPendingCodeAt(t *testing.T) {
	ec := newClient()

	// Test 1 with Error
	ctx := newContext(fmt.Errorf("test-error"), 0, nil)
	_, err := ec.PendingCodeAt(ctx, "test-endpoint", ethcommon.Address{})
	assert.Error(t, err, "#1 PendingCodeAt should  error")

	// Test 2 without error
	expectedCode := hexutil.MustDecode("0xabcd")
	ctx = newContext(nil, 200, makeRespBody(hexutil.Bytes(expectedCode), ""))
	code, err := ec.PendingCodeAt(ctx, "test-endpoint", ethcommon.Address{})
	assert.NoError(t, err, "#3 PendingCodeAt should not error")
	assert.Equal(t, hexutil.Encode(expectedCode), hexutil.Encode(code), "#3 PendingCodeAt code should be correct")
}

func TestPendingNonceAt(t *testing.T) {
	ec := newClient()

	// Test 1 with Error
	ctx := newContext(fmt.Errorf("test-error"), 0, nil)
	_, err := ec.PendingNonceAt(ctx, "test-endpoint", ethcommon.Address{})
	assert.Error(t, err, "#1 PendingNonceAt should  error")

	// Test 2 without error
	expectedNonce := uint64(17)
	ctx = newContext(nil, 200, makeRespBody(hexutil.Uint64(expectedNonce), ""))
	nonce, err := ec.PendingNonceAt(ctx, "test-endpoint", ethcommon.Address{})
	assert.NoError(t, err, "#3 PendingNonceAt should not error")
	assert.Equal(t, expectedNonce, nonce, "#3 PendingNonceAt nonce should be correct")
}

func TestCallContract(t *testing.T) {
	ec := newClient()

	// Test 1 with Error
	ctx := newContext(fmt.Errorf("test-error"), 0, nil)
	_, err := ec.CallContract(ctx, "test-endpoint", &eth.CallMsg{}, nil)
	assert.Error(t, err, "#1 CallContract should  error")

	// Test 2 without error
	expectedContract := hexutil.MustDecode("0xabcd")
	ctx = newContext(nil, 200, makeRespBody(hexutil.Bytes(expectedContract), ""))
	contract, err := ec.CallContract(ctx, "test-endpoint", &eth.CallMsg{}, nil)
	assert.NoError(t, err, "#3 CallContract should not error")
	assert.Equal(t, hexutil.Encode(expectedContract), hexutil.Encode(contract), "#3 CallContract code should be correct")
}

func TestPendingCallContract(t *testing.T) {
	ec := newClient()

	// Test 1 with Error
	ctx := newContext(fmt.Errorf("test-error"), 0, nil)
	_, err := ec.PendingCallContract(ctx, "test-endpoint", &eth.CallMsg{})
	assert.Error(t, err, "#1 PendingCallContract should  error")

	// Test 2 without error
	expectedContract := hexutil.MustDecode("0xabcd")
	ctx = newContext(nil, 200, makeRespBody(hexutil.Bytes(expectedContract), ""))
	contract, err := ec.PendingCallContract(ctx, "test-endpoint", &eth.CallMsg{})
	assert.NoError(t, err, "#3 PendingCallContract should not error")
	assert.Equal(t, hexutil.Encode(expectedContract), hexutil.Encode(contract), "#3 PendingCallContract code should be correct")
}

func TestSuggestGasPrice(t *testing.T) {
	ec := newClient()

	// Test 1 with Error
	ctx := newContext(fmt.Errorf("test-error"), 0, nil)
	_, err := ec.SuggestGasPrice(ctx, "test-endpoint")
	assert.Error(t, err, "#1 SuggestGasPrice should  error")

	// Test 2 without error
	expectedPrice := big.NewInt(1000)
	ctx = newContext(nil, 200, makeRespBody((*hexutil.Big)(expectedPrice), ""))
	price, err := ec.SuggestGasPrice(ctx, "test-endpoint")
	assert.NoError(t, err, "#3 SuggestGasPrice should not error")
	assert.Equal(t, expectedPrice.Text(10), price.Text(10), "#3 SuggestGasPrice balance should be correct")
}

func TestEstimateGas(t *testing.T) {
	ec := newClient()

	// Test 1 with Error
	ctx := newContext(fmt.Errorf("test-error"), 0, nil)
	_, err := ec.EstimateGas(ctx, "test-endpoint", &eth.CallMsg{})
	assert.Error(t, err, "#1 EstimateGas should  error")

	// Test 2 without error
	expectedGas := uint64(17)
	ctx = newContext(nil, 200, makeRespBody(hexutil.Uint64(expectedGas), ""))
	gas, err := ec.EstimateGas(ctx, "test-endpoint", &eth.CallMsg{})
	assert.NoError(t, err, "#3 EstimateGas should not error")
	assert.Equal(t, expectedGas, gas, "#3 EstimateGas nonce should be correct")
}

func TestSendRawTransaction(t *testing.T) {
	ec := newClient()

	// Test 1 with Error
	ctx := newContext(fmt.Errorf("test-error"), 0, nil)
	err := ec.SendRawTransaction(ctx, "test-endpoint", "")
	assert.Error(t, err, "#1 SendRawTransaction should  error")

	// Test 2 without error
	ctx = newContext(nil, 200, makeRespBody("", ""))
	err = ec.SendRawTransaction(ctx, "test-endpoint", "")
	assert.NoError(t, err, "#2 SendRawTransaction should not error")

	// Test 3 with Nonce Too Low error
	ctx = newContext(nil, 200, makeRespBody("", "Nonce too low"))
	err = ec.SendRawTransaction(ctx, "test-endpoint", "")
	assert.Error(t, err, "#2 SendRawTransaction should not error")
}

func TestSendTransaction(t *testing.T) {
	ec := newClient()

	// Test 1 with Error
	ctx := newContext(fmt.Errorf("test-error"), 0, nil)
	_, err := ec.SendTransaction(ctx, "test-endpoint", &types.SendTxArgs{})
	assert.Error(t, err, "#1 SendTransaction should  error")
}

func TestSendQuorumRawPrivateTransaction(t *testing.T) {
	ec := newClient()

	// Test 1 with Error
	ctx := newContext(fmt.Errorf("test-error"), 0, nil)
	_, err := ec.SendQuorumRawPrivateTransaction(ctx, "test-endpoint", "", nil)
	assert.Error(t, err, "#1 SendQuorumRawPrivateTransaction should  error")
}

func TestSendRawPrivateTransaction(t *testing.T) {
	ec := newClient()

	// Test 1 with Error
	ctx := newContext(fmt.Errorf("test-error"), 0, nil)
	_, err := ec.SendRawPrivateTransaction(ctx, "test-endpoint", "")
	assert.Error(t, err, "#1 SendRawPrivateTransaction should  error")
}

func TestNetwork(t *testing.T) {
	ec := newClient()

	// Test 1 with Error
	ctx := newContext(fmt.Errorf("test-error"), 0, nil)
	_, err := ec.Network(ctx, "test-endpoint")
	assert.Error(t, err, "#1 Network should  error")

	// Test 2 without error
	ctx = newContext(nil, 200, makeRespBody("1234", ""))
	chain, err := ec.Network(ctx, "test-endpoint")
	assert.NoError(t, err, "#2 Network should not error")
	assert.Equal(t, uint64(1234), chain.Uint64(), "#2 Chain id should match")

	// Test 3 without encoding format
	ctx = newContext(nil, 200, makeRespBody("%/", ""))
	_, err = ec.Network(ctx, "test-endpoint")
	assert.Error(t, err, "#3 Network should error")
}
