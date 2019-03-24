package ethclient

import (
	"context"
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/trace"
)

// EthClient embed a go-ethereum ethclient so we can implement SendRawTransaction
type EthClient struct {
	*ethclient.Client
	rpc *rpc.Client
}

// NewClient creates a client that uses the given RPC client.
func NewClient(c *rpc.Client) *EthClient {
	ec := ethclient.NewClient(c)
	return &EthClient{ec, c}
}

// Dial connects a client to the given URL.
func Dial(rawurl string) (*EthClient, error) {
	return DialContext(context.Background(), rawurl)
}

// DialContext connects a client to the given URL.
func DialContext(ctx context.Context, rawurl string) (*EthClient, error) {
	c, err := rpc.DialContext(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	return NewClient(c), nil
}

// PrivateArgs are arguments to provide to an Ethereum client supporting privacy client (such as Quorum) when sendng a transaction
type PrivateArgs struct {
	// Quorum Fields
	PrivateFrom   string   `json:"privateFrom"`
	PrivateFor    []string `json:"privateFor"`
	PrivateTxType string   `json:"restriction"`
}

// SendTxArgs are arguments to provide to jsonRPC call `eth_sendTransaction`
type SendTxArgs struct {
	// From address in case of a non raw transaction
	From ethcommon.Address `json:"from"`

	// Main transaction attributes
	To       *ethcommon.Address `json:"to"`
	Gas      *hexutil.Uint64    `json:"gas"`
	GasPrice *hexutil.Big       `json:"gasPrice"`
	Value    *hexutil.Big       `json:"value"`
	Nonce    *hexutil.Uint64    `json:"nonce"`

	// We accept "data" and "input" for backwards-compatibility reasons. "input" is the
	// newer name and should be preferred
	Data  hexutil.Bytes `json:"data"`
	Input hexutil.Bytes `json:"input"`

	// Private field
	PrivateArgs
}

// Call2PrivateArgs creates PrivateArgs from a call object
func Call2PrivateArgs(call *common.Call) *PrivateArgs {
	var args PrivateArgs
	args.PrivateFrom = call.GetQuorum().GetPrivateFrom()
	args.PrivateFor = call.GetQuorum().GetPrivateFor()
	args.PrivateTxType = call.GetQuorum().GetPrivateTxType()
	return &args
}

// Trace2SendTxArgs creates SendTxArgs from a trace
func Trace2SendTxArgs(tr *trace.Trace) *SendTxArgs {
	args := SendTxArgs{
		From:        tr.GetSender().Address(),
		GasPrice:    (*hexutil.Big)(tr.GetTx().GetTxData().GasPriceBig()),
		Value:       (*hexutil.Big)(tr.GetTx().GetTxData().ValueBig()),
		Data:        hexutil.Bytes(tr.GetTx().GetTxData().DataBytes()),
		Input:       hexutil.Bytes(tr.GetTx().GetTxData().DataBytes()),
		PrivateArgs: *(Call2PrivateArgs(tr.GetCall())),
	}

	if gas := tr.GetTx().GetTxData().GetGas(); gas != 0 {
		args.Gas = (*hexutil.Uint64)(&gas)
	}

	if nonce := tr.GetTx().GetTxData().GetNonce(); nonce != 0 {
		args.Nonce = (*hexutil.Uint64)(&nonce)
	}

	if tr.GetTx().GetTxData().GetTo() != "" {
		to := tr.GetTx().GetTxData().ToAddress()
		args.To = &to
	}

	return &args
}

// SendRawTransaction allows to send a raw transaction
func (ec *EthClient) SendRawTransaction(ctx context.Context, raw string) error {
	return ec.rpc.CallContext(ctx, nil, "eth_sendRawTransaction", raw)
}

// SendTransaction send transaction to Ethereum node
func (ec *EthClient) SendTransaction(ctx context.Context, args *SendTxArgs) (txHash ethcommon.Hash, err error) {
	err = ec.rpc.CallContext(ctx, &txHash, "eth_sendTransaction", &args)
	if err != nil {
		return ethcommon.Hash{}, err
	}
	return txHash, nil
}

// SendRawPrivateTransaction send a raw transaction to a Ethreum node supporting privacy (e.g Quorum+Tessera node)
// TODO: to be implemented
func (ec *EthClient) SendRawPrivateTransaction(ctx context.Context, raw string, q *PrivateArgs) (ethcommon.Hash, error) {
	return ethcommon.Hash{}, fmt.Errorf("%q is not implemented yet", "SendRawPrivateTransactionQuorum")
}
