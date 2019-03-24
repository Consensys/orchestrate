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

// QuorumArgs are arguments to provide to a Quorum jsonRPC call
type QuorumArgs struct {
	// Quorum Fields
	PrivateFrom   string   `json:"privateFrom"`
	PrivateFor    []string `json:"privateFor"`
	PrivateTxType string   `json:"restriction"`
}

// SendTxArgs are arguments to provide to jsonRPC call eth_sendTransaction
type SendTxArgs struct {
	// From address in case of a non raw transaction
	From ethcommon.Address `json:"from"`

	// Main transaction attributes
	To       *ethcommon.Address `json:"to"`
	Gas      hexutil.Uint64     `json:"gas"`
	GasPrice *hexutil.Big       `json:"gasPrice"`
	Value    *hexutil.Big       `json:"value"`
	Nonce    hexutil.Uint64     `json:"nonce"`

	// We accept "data" and "input" for backwards-compatibility reasons. "input" is the
	// newer name and should be preferred
	Data  hexutil.Bytes `json:"data"`
	Input hexutil.Bytes `json:"input"`

	// Quorum Fields
	QuorumArgs
}

// Call2QuorumArgs creates QuorumArgs from a call object
func Call2QuorumArgs(call *common.Call) *QuorumArgs {
	var args QuorumArgs
	args.PrivateFrom = call.GetQuorum().GetPrivateFrom()
	args.PrivateFor = call.GetQuorum().GetPrivateFor()
	args.PrivateTxType = call.GetQuorum().GetPrivateTxType()
	return &args
}

// Trace2SendTxArgs creates SendTxArgs from a trace
func Trace2SendTxArgs(tr *trace.Trace) *SendTxArgs {
	args := SendTxArgs{
		From:       tr.GetSender().Address(),
		Gas:        hexutil.Uint64(tr.GetTx().GetTxData().GetGas()),
		GasPrice:   (*hexutil.Big)(tr.GetTx().GetTxData().GasPriceBig()),
		Value:      (*hexutil.Big)(tr.GetTx().GetTxData().ValueBig()),
		Nonce:      hexutil.Uint64(tr.GetTx().GetTxData().GetNonce()),
		Data:       hexutil.Bytes(tr.GetTx().GetTxData().DataBytes()),
		Input:      hexutil.Bytes(tr.GetTx().GetTxData().DataBytes()),
		QuorumArgs: *(Call2QuorumArgs(tr.GetCall())),
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

// SendPrivateTransactionQuorum send transaction to Quorum node
func (ec *EthClient) SendPrivateTransactionQuorum(ctx context.Context, args *SendTxArgs) (txHash ethcommon.Hash, err error) {
	err = ec.rpc.CallContext(ctx, &txHash, "eth_sendTransaction", &args)
	if err != nil {
		return ethcommon.Hash{}, err
	}
	return txHash, nil
}

// SendRawPrivateTransactionQuorum send a raw transaction to a Quorum node (only compatible if Quorum node uses Tessera)
// TODO: to be implemented
func (ec *EthClient) SendRawPrivateTransactionQuorum(ctx context.Context, raw string, q *QuorumArgs) (ethcommon.Hash, error) {
	return ethcommon.Hash{}, fmt.Errorf("%q is not implemented yet", "SendRawPrivateTransactionQuorum")
}
