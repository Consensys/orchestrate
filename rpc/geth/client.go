package geth

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
	eth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/logger"
)

// Client is a wrapper around Geth rpc client supporting Backoff retry strategy
type Client struct {
	rpc *rpc.Client

	pool *sync.Pool
	conf *Config
}

// Dial creates a new client for the given URL.
func Dial(rawurl string) (*Client, error) {
	return DialFromConfig(rawurl, NewConfig())
}

// Dial creates a new client for the given URL.
//
// The currently supported URL schemes are "http", "https", "ws" and "wss". If rawurl is a
// file name with no URL scheme, a local socket connection is established using UNIX
// domain sockets on supported platforms and named pipes on Windows. If you want to
// configure transport options, use DialHTTP, DialWebsocket or DialIPC instead.
//
// For websocket connections, the origin is set to the local host name.
//
// The client reconnects automatically if the connection is lost.
func DialFromConfig(rawurl string, conf *Config) (*Client, error) {
	return DialContext(context.Background(), rawurl, conf)
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

// DialContext creates a new RPC client, just like Dial.

// The context is used to cancel or time out the initial connection establishment. It does
// not affect subsequent interactions with the client.
func DialContext(ctx context.Context, rawurl string, conf *Config) (*Client, error) {
	c, err := rpc.DialContext(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	return &Client{
		rpc:  c,
		conf: conf,
		pool: &sync.Pool{
			New: func() interface{} { return NewBackOff(conf) },
		},
	}, nil
}

// Close closes the client, aborting any in-flight requests.
func (c *Client) Close() {
	c.rpc.Close()
}

// CallContext performs a JSON-RPC call with the given arguments. If the context is
// canceled before the call has successfully returned, CallContext returns immediately.
//
// The result must be a pointer so that package json can unmarshal into it. You
// can also pass nil, in which case the result is ignored.
func (c *Client) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	bckoff := backoff.WithContext(c.pool.Get().(backoff.BackOff), ctx)
	defer c.pool.Put(bckoff)

	return backoff.RetryNotify(
		func() error {
			var raw json.RawMessage
			log.Debugf("calling method %s(%+v)", method, args)
			err := c.rpc.CallContext(ctx, &raw, method, args...)
			if err != nil {
				log.Debugf("failed to call %s(%+v)", method, args)
				return err
			} else if len(raw) == 0 {
				log.Debugf("%s(%+v) returned empty result", method, args)
				return eth.NotFound
			}

			if err := json.Unmarshal(raw, &result); err != nil {
				log.Debugf("failed to parse the result of call %s(%+v)", method, args)
				return err
			}

			return nil
		},
		bckoff,
		func(err error, duration time.Duration) {
			logger.GetLogEntry(ctx).
				WithError(err).
				WithFields(log.Fields{
					"method": method,
				}).Warnf("eth-client: error on JSON-RPC call, retrying in %v...", duration)
		},
	)
}
