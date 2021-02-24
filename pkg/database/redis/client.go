package redis

import (
	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/log"
	"github.com/gomodule/redigo/redis"
)

const cannotCloseConnErr = "could not close connection"

// NonceManager manages nonce using an underlying redis cache
type Client struct {
	pool   *redis.Pool
	conf   *Config
	logger *log.Logger
}

func NewClient(pool *redis.Pool, conf *Config) *Client {
	return &Client{
		pool:   pool,
		conf:   conf,
		logger: log.NewLogger().SetComponent(component).WithField("host", conf.URL()),
	}
}

func (nm *Client) Load(key string) (value interface{}, ok bool, err error) {
	conn := nm.pool.Get()
	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			nm.logger.WithError(closeErr).Error(cannotCloseConnErr)
		}
	}()

	reply, err := conn.Do("GET", key)
	if err != nil {
		return reply, false, errors.RedisConnectionError("failed to set value").AppendReason(err.Error())
	}

	if reply == nil {
		return nil, false, nil
	}

	return reply, true, nil
}

func (nm *Client) LoadUint64(key string) (value uint64, ok bool, err error) {
	// Load value
	reply, ok, err := nm.Load(key)
	if err != nil || !ok {
		return 0, false, err
	}

	// Format reply to Uint64
	value, err = redis.Uint64(reply, nil)
	if err != nil {
		return 0, false, parseRedisError(err, "failed to load UInt64 value")
	}

	return value, true, nil
}

func (nm *Client) LoadBool(key string) (value, ok bool, err error) {
	// Load value
	reply, ok, err := nm.Load(key)
	if err != nil || !ok {
		return false, false, err
	}

	// Format reply to Uint64
	value, err = redis.Bool(reply, nil)
	if err != nil {
		return false, false, parseRedisError(err, "failed to load Boolean value")
	}

	return value, true, nil
}

func (nm *Client) Set(key string, value interface{}) error {
	conn := nm.pool.Get()
	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			nm.logger.WithError(closeErr).Error(cannotCloseConnErr)
		}
	}()

	// Set value with expiration
	_, err := conn.Do("PSETEX", key, nm.conf.Expiration, value)
	if err != nil {
		return errors.FromError(err).SetComponent(component)
	}

	return nil
}

func (nm *Client) Delete(key string) error {
	conn := nm.pool.Get()
	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			nm.logger.WithError(closeErr).Warn("could not close connection")
		}
	}()

	// Delete value
	_, err := conn.Do("DEL", key)
	if err != nil {
		return parseRedisError(err, "failed to delete key")
	}

	return nil
}

func (nm *Client) Incr(key string) error {
	conn := nm.pool.Get()
	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			nm.logger.WithError(closeErr).Warn("could not close redis connection")
		}
	}()

	_, err := conn.Do("INCR", key)
	if err != nil {
		return parseRedisError(err, "failed to increment value")
	}

	return nil
}

func (nm *Client) Ping() error {
	conn := nm.pool.Get()

	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			nm.logger.WithError(closeErr).Warn("could not close redis connection")
		}
	}()

	_, err := conn.Do("PING")
	if err != nil {
		return parseRedisError(err, "failed to ping")
	}

	return nil
}
