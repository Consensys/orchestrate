package ethclient

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
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

// SendRawTransaction allows to send a raw transaction
func (ec *EthClient) SendRawTransaction(ctx context.Context, raw string) error {
	return ec.rpc.CallContext(ctx, nil, "eth_sendRawTransaction", raw)
}

// SendTxArgs are arguments to provide to JSONRPC call on Quorum
type SendTxArgs struct {
	From     common.Address  `json:"from"`
	To       *common.Address `json:"to"`
	Gas      hexutil.Uint64  `json:"gas"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
	Value    *hexutil.Big    `json:"value"`
	Nonce    hexutil.Uint64  `json:"nonce"`

	// We accept "data" and "input" for backwards-compatibility reasons. "input" is the
	// newer name and should be preferred by clients.
	Data  hexutil.Bytes `json:"data"`
	Input hexutil.Bytes `json:"input"`

	//Quorum
	PrivateFrom   string   `json:"privateFrom"`
	PrivateFor    []string `json:"privateFor"`
	PrivateTxType string   `json:"restriction"`
}

// Trace2SendTxArgs create SendTxQuorumArgs from a trace
func Trace2SendTxArgs(tr *trace.Trace) *SendTxArgs {
	var args SendTxArgs
	args.From = tr.GetSender().Address()
	if tr.GetTx().GetTxData().GetTo() != "" {
		to := tr.GetTx().GetTxData().ToAddress()
		args.To = &to
	}
	args.Gas = hexutil.Uint64(tr.GetTx().GetTxData().GetGas())
	args.GasPrice = (*hexutil.Big)(tr.GetTx().GetTxData().GasPriceBig())
	args.Value = (*hexutil.Big)(tr.GetTx().GetTxData().ValueBig())
	args.Nonce = hexutil.Uint64(tr.GetTx().GetTxData().GetNonce())
	args.Data = hexutil.Bytes(tr.GetTx().GetTxData().DataBytes())
	args.Input = hexutil.Bytes(tr.GetTx().GetTxData().DataBytes())
	args.PrivateFrom = tr.GetCall().GetQuorum().GetPrivateFrom()
	args.PrivateFor = tr.GetCall().GetQuorum().GetPrivateFor()
	args.PrivateTxType = tr.GetCall().GetQuorum().GetPrivateTxType()
	return &args
}

// SendPrivateTransactionQuorum send transaction to Quorum node
func (ec *EthClient) SendPrivateTransactionQuorum(ctx context.Context, args *SendTxArgs) (txHash common.Hash, err error) {
	err = ec.rpc.CallContext(ctx, &txHash, "eth_sendTransaction", &args)
	if err != nil {
		return common.Hash{}, err
	}

	return txHash, nil
}

// SendRawPrivateTransactionQuorum send a raw transaction to a Quorum node (only compatible if Quorum node uses Tessera)
// TODO: to be implemented
func (ec *EthClient) SendRawPrivateTransactionQuorum(ctx context.Context, args *SendTxArgs) (common.Hash, error) {
	return common.Hash{}, nil
}
