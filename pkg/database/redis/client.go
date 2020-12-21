package redis

import (
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
)

// NonceManager manages nonce using an underlying redis cache
type Client struct {
	pool *redis.Pool
	conf *Config
}

func NewClient(pool *redis.Pool, conf *Config) *Client {
	return &Client{
		pool: pool,
		conf: conf,
	}
}

func (nm *Client) Load(key string) (value interface{}, ok bool, err error) {
	conn := nm.pool.Get()
	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			log.WithError(closeErr).Warn("could not close redis connection")
		}
	}()

	reply, err := conn.Do("GET", key)
	if err != nil {
		return reply, false, errors.FromError(err).SetComponent(component)
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
		return 0, false, FromRedisError(err).SetComponent(component)
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
		return false, false, FromRedisError(err).SetComponent(component)
	}

	return value, true, nil
}

func (nm *Client) Set(key string, value interface{}) error {
	conn := nm.pool.Get()
	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			log.WithError(closeErr).Warn("could not close redis connection")
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
			log.WithError(closeErr).Warn("could not close redis connection")
		}
	}()

	// Delete value
	_, err := conn.Do("DEL", key)
	if err != nil {
		return errors.FromError(err).SetComponent(component)
	}

	return nil
}

func (nm *Client) Incr(key string) error {
	conn := nm.pool.Get()
	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			log.WithError(closeErr).Warn("could not close redis connection")
		}
	}()

	_, err := conn.Do("INCR", key)
	if err != nil {
		return errors.FromError(err).SetComponent(component)
	}

	return nil
}

func (nm *Client) Ping() error {
	conn := nm.pool.Get()

	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			log.WithError(closeErr).Warn("could not close redis connection")
		}
	}()

	_, err := conn.Do("PING")
	if err != nil {
		return errors.FromError(err).SetComponent(component)
	}

	return nil
}
