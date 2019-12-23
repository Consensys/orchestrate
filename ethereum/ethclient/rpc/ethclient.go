package rpc

import (
	"context"
	"math/big"
	"sync"

	eth "github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

// Client is a connector to Ethereum blockchains that uses Geth rpc client
type Client struct {
	mux  *sync.RWMutex
	urls map[string]string

	clientV2 *ClientV2
}

// NewClient creates a new MultiClient
func NewClient(conf *Config) *Client {
	return &Client{
		mux:      &sync.RWMutex{},
		urls:     make(map[string]string),
		clientV2: NewClientV2(conf),
	}
}

func (ec *Client) Dial(ctx context.Context, rawurl string) (*big.Int, error) {
	// Retrieve network version
	var version string
	if err := ec.clientV2.Call(ctx, rawurl, processResult(&version), "net_version"); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	chainID, ok := big.NewInt(0).SetString(version, 10)
	if !ok {
		return nil, errors.InvalidFormatError("invalid chain identifier %v", chainID)
	}

	// Register client
	ec.mux.Lock()
	ec.urls[chainID.Text(10)] = rawurl
	ec.mux.Unlock()

	return chainID, nil
}

func (ec *Client) getURL(chainID *big.Int) (string, error) {
	ec.mux.RLock()
	defer ec.mux.RUnlock()
	c, ok := ec.urls[chainID.Text(10)]
	if !ok {
		return "", errors.EthConnectionError("no RPC connection registered for chain %q", chainID.Text(10))
	}
	return c, nil
}

func (ec *Client) BlockByHash(ctx context.Context, chainID *big.Int, hash ethcommon.Hash) (*ethtypes.Block, error) {
	url, err := ec.getURL(chainID)
	if err != nil {
		return nil, err
	}
	return ec.clientV2.BlockByHash(ctx, url, hash)
}

func (ec *Client) BlockByNumber(ctx context.Context, chainID, number *big.Int) (*ethtypes.Block, error) {
	url, err := ec.getURL(chainID)
	if err != nil {
		return nil, err
	}
	return ec.clientV2.BlockByNumber(ctx, url, number)
}

func (ec *Client) HeaderByHash(ctx context.Context, chainID *big.Int, hash ethcommon.Hash) (*ethtypes.Header, error) {
	url, err := ec.getURL(chainID)
	if err != nil {
		return nil, err
	}
	return ec.clientV2.HeaderByHash(ctx, url, hash)
}

func (ec *Client) HeaderByNumber(ctx context.Context, chainID, number *big.Int) (*ethtypes.Header, error) {
	url, err := ec.getURL(chainID)
	if err != nil {
		return nil, err
	}
	return ec.clientV2.HeaderByNumber(ctx, url, number)
}

func (ec *Client) TransactionByHash(ctx context.Context, chainID *big.Int, hash ethcommon.Hash) (tx *ethtypes.Transaction, isPending bool, err error) {
	url, err := ec.getURL(chainID)
	if err != nil {
		return nil, false, err
	}
	return ec.clientV2.TransactionByHash(ctx, url, hash)
}

func (ec *Client) TransactionReceipt(ctx context.Context, chainID *big.Int, txHash ethcommon.Hash) (*ethtypes.Receipt, error) {
	url, err := ec.getURL(chainID)
	if err != nil {
		return nil, err
	}
	return ec.clientV2.TransactionReceipt(ctx, url, txHash)
}

func (ec *Client) Networks(ctx context.Context) (networks []*big.Int) {
	ec.mux.RLock()
	defer ec.mux.RUnlock()
	for _, url := range ec.urls {
		// Retrieve network version
		var version string
		if err := ec.clientV2.Call(ctx, url, processResult(&version), "net_version"); err != nil {
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

func (ec *Client) SyncProgress(ctx context.Context, chainID *big.Int) (*eth.SyncProgress, error) {
	url, err := ec.getURL(chainID)
	if err != nil {
		return nil, err
	}
	return ec.clientV2.SyncProgress(ctx, url)
}

func (ec *Client) BalanceAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, blockNumber *big.Int) (*big.Int, error) {
	url, err := ec.getURL(chainID)
	if err != nil {
		return nil, err
	}
	return ec.clientV2.BalanceAt(ctx, url, account, blockNumber)
}

func (ec *Client) StorageAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, key ethcommon.Hash, blockNumber *big.Int) ([]byte, error) {
	url, err := ec.getURL(chainID)
	if err != nil {
		return nil, err
	}
	return ec.clientV2.StorageAt(ctx, url, account, key, blockNumber)
}

func (ec *Client) CodeAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, blockNumber *big.Int) ([]byte, error) {
	url, err := ec.getURL(chainID)
	if err != nil {
		return nil, err
	}
	return ec.clientV2.CodeAt(ctx, url, account, blockNumber)
}

func (ec *Client) NonceAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, blockNumber *big.Int) (uint64, error) {
	url, err := ec.getURL(chainID)
	if err != nil {
		return 0, err
	}
	return ec.clientV2.NonceAt(ctx, url, account, blockNumber)
}

func (ec *Client) PendingBalanceAt(ctx context.Context, chainID *big.Int, account ethcommon.Address) (*big.Int, error) {
	url, err := ec.getURL(chainID)
	if err != nil {
		return nil, err
	}
	return ec.clientV2.PendingBalanceAt(ctx, url, account)
}

func (ec *Client) PendingStorageAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, key ethcommon.Hash) ([]byte, error) {
	url, err := ec.getURL(chainID)
	if err != nil {
		return nil, err
	}
	return ec.clientV2.PendingStorageAt(ctx, url, account, key)
}

func (ec *Client) PendingCodeAt(ctx context.Context, chainID *big.Int, account ethcommon.Address) ([]byte, error) {
	url, err := ec.getURL(chainID)
	if err != nil {
		return nil, err
	}
	return ec.clientV2.PendingCodeAt(ctx, url, account)
}

func (ec *Client) PendingNonceAt(ctx context.Context, chainID *big.Int, account ethcommon.Address) (uint64, error) {
	url, err := ec.getURL(chainID)
	if err != nil {
		return 0, err
	}
	return ec.clientV2.PendingNonceAt(ctx, url, account)
}

func (ec *Client) CallContract(ctx context.Context, chainID *big.Int, msg *eth.CallMsg, blockNumber *big.Int) ([]byte, error) {
	url, err := ec.getURL(chainID)
	if err != nil {
		return nil, err
	}
	return ec.clientV2.CallContract(ctx, url, msg, blockNumber)
}

func (ec *Client) PendingCallContract(ctx context.Context, chainID *big.Int, msg *eth.CallMsg) ([]byte, error) {
	url, err := ec.getURL(chainID)
	if err != nil {
		return nil, err
	}
	return ec.clientV2.PendingCallContract(ctx, url, msg)
}

func (ec *Client) SuggestGasPrice(ctx context.Context, chainID *big.Int) (*big.Int, error) {
	url, err := ec.getURL(chainID)
	if err != nil {
		return nil, err
	}
	return ec.clientV2.SuggestGasPrice(ctx, url)
}

func (ec *Client) EstimateGas(ctx context.Context, chainID *big.Int, msg *eth.CallMsg) (uint64, error) {
	url, err := ec.getURL(chainID)
	if err != nil {
		return 0, err
	}
	return ec.clientV2.EstimateGas(ctx, url, msg)
}

func (ec *Client) SendRawTransaction(ctx context.Context, chainID *big.Int, raw string) error {
	url, err := ec.getURL(chainID)
	if err != nil {
		return err
	}
	return ec.clientV2.SendRawTransaction(ctx, url, raw)
}

func (ec *Client) SendTransaction(ctx context.Context, chainID *big.Int, args *types.SendTxArgs) (txHash ethcommon.Hash, err error) {
	url, err := ec.getURL(chainID)
	if err != nil {
		return ethcommon.Hash{}, err
	}
	return ec.clientV2.SendTransaction(ctx, url, args)
}

func (ec *Client) SendQuorumRawPrivateTransaction(ctx context.Context, chainID *big.Int, signedTxHash []byte, privateFor []string) (ethcommon.Hash, error) {
	url, err := ec.getURL(chainID)
	if err != nil {
		return ethcommon.Hash{}, err
	}
	return ec.clientV2.SendQuorumRawPrivateTransaction(ctx, url, signedTxHash, privateFor)
}

func (ec *Client) SendRawPrivateTransaction(ctx context.Context, chainID *big.Int, raw []byte, args *types.PrivateArgs) (ethcommon.Hash, error) {
	url, err := ec.getURL(chainID)
	if err != nil {
		return ethcommon.Hash{}, err
	}
	return ec.clientV2.SendRawPrivateTransaction(ctx, url, raw, args)
}
