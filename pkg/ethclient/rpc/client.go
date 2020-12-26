package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/containous/traefik/v2/pkg/log"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient/utils"
)

type ProcessResultFunc func(result json.RawMessage) error

// Client is a connector to Ethereum blockchains that uses Geth rpc client
type Client struct {
	client *http.Client

	// Pool for backoffs
	pool *sync.Pool

	idCounter uint32
}

// NewClient creates a new MultiClient
func NewClient(newBackOff func() backoff.BackOff, client *http.Client) *Client {
	return &Client{
		client: client,
		pool: &sync.Pool{
			New: func() interface{} { return newBackOff() },
		},
		idCounter: 0,
	}
}

func (ec *Client) Call(ctx context.Context, endpoint string, processResult func(result json.RawMessage) error, method string, args ...interface{}) error {
	bckoff := backoff.WithContext(ec.pool.Get().(backoff.BackOff), ctx)
	defer ec.pool.Put(bckoff)

	return ec.callWithRetry(ctx, func(context.Context) (*http.Request, error) {
		return ec.newJSONRpcRequestWithContext(ctx, endpoint, method, args...)
	}, processResult, bckoff)
}

func (ec *Client) Pool() *sync.Pool {
	return ec.pool
}

func (ec *Client) HTTPClient() *http.Client {
	return ec.client
}

func (ec *Client) callWithRetry(ctx context.Context, reqBuilder func(context.Context) (*http.Request, error),
	processResult func(result json.RawMessage) error, bckoff backoff.BackOff) error {
	return backoff.RetryNotify(
		func() error {
			// Every request we generate a new object
			req, err := reqBuilder(ctx)
			if err != nil {
				return backoff.Permanent(err)
			}

			e := ec.call(req, processResult)
			switch {
			case e == nil:
				return nil
			// Capture NotFoundData RPC Error and replace by InvalidParameterError to prevent 404 response
			case errors.IsNotFoundError(e):
				return backoff.Permanent(errors.InvalidParameterError(e.Error()))
			// Retry on timeout of temporally out of order AND on eth connection errors
			case errors.IsConnectionError(e) && utils.ShouldRetryConnectionError(ctx):
				return e
			default:
				return backoff.Permanent(e)
			}
		},
		bckoff,
		func(e error, duration time.Duration) {
			log.FromContext(ctx).
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

	var respMsg utils.JSONRpcMessage
	if resp.Body != nil {
		defer func() {
			err = resp.Body.Close()
			if err != nil {
				log.FromContext(req.Context()).
					WithError(err).
					Warn("could not close request body")
			}
		}()

		err = json.NewDecoder(resp.Body).Decode(&respMsg)
	}

	switch {
	case err == nil && respMsg.Error != nil:
		return ec.processEthError(respMsg.Error)
	case err == nil && len(respMsg.Result) == 0:
		return errors.NotFoundError("data not found")
	case resp.StatusCode == 404:
		return errors.NotFoundError("url %s not found", req.URL)
	case resp.StatusCode < 200 || resp.StatusCode >= 300:
		return errors.EthConnectionError("%v (code=%v)", resp.Status, resp.StatusCode)
	case err != nil:
		return errors.EncodingError("invalid RPC response")
	default:
		return processResult(respMsg.Result)
	}
}

func (ec *Client) newJSONRpcMessage(method string, args ...interface{}) (*utils.JSONRpcMessage, error) {
	id := ec.nextID()
	msg := &utils.JSONRpcMessage{
		Method:  method,
		Version: "2.0",
		ID:      id,
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
		log.FromContext(req.Context()).WithError(err).Warn("connection error")
		rerr, ok := err.(*url.Error)
		// We consider these two error types as recoverable
		if ok && (rerr.Timeout() || rerr.Temporary()) {
			return nil, errors.ServiceConnectionError(err.Error())
		}

		return nil, errors.ConnectionError(err.Error())
	}

	return resp, nil
}

func (ec *Client) nextID() json.RawMessage {
	id := atomic.AddUint32(&ec.idCounter, 1)
	return strconv.AppendUint(nil, uint64(id), 10)
}

func (ec *Client) processEthError(err *utils.JSONError) error {
	if strings.Contains(err.Message, "nonce too low") || strings.Contains(err.Message, "Nonce too low") || strings.Contains(err.Message, "Incorrect nonce") {
		return errors.NonceTooLowError("code: %d - message: %s", err.Code, err.Message)
	}
	return errors.EthereumError("code: %d - message: %s", err.Code, err.Message)
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

func processBlockResult(header **ethtypes.Header, body **Body) ProcessResultFunc {
	return func(result json.RawMessage) error {
		var raw json.RawMessage
		err := utils.ProcessResult(&raw)(result)
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
