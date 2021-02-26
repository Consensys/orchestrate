package redis

import (
	"time"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/tls"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/gomodule/redigo/redis"
)

// Conn is a wrapper around a redis.Conn that handles internal errors
type Conn struct {
	redis.Conn
}

// Dial connects to the Redis server
func Dial(network, address string, options ...redis.DialOption) (redis.Conn, error) {
	conn, err := redis.Dial(network, address, options...)
	if err != nil {
		return conn, errors.ConnectionError(err.Error())
	}
	return Conn{conn}, nil
}

func (conn Conn) Do(commandName string, args ...interface{}) (interface{}, error) {
	reply, err := conn.Conn.Do(commandName, args...)
	if err != nil {
		return reply, errors.RedisConnectionError(err.Error())
	}
	return reply, nil
}

// Creates a new redis pool
func NewPool(cfg *Config) (*redis.Pool, error) {
	options := []redis.DialOption{}
	if cfg.Database != databaseDefault {
		options = append(options, redis.DialDatabase(cfg.Database))
	}

	if cfg.User != "" {
		options = append(options, redis.DialUsername(cfg.User))
	}

	if cfg.Password != "" {
		options = append(options, redis.DialPassword(cfg.Password))
	}

	if cfg.TLS != nil {
		c, err := tls.NewConfig(cfg.TLS)
		if err != nil {
			return nil, err
		}

		options = append(options, redis.DialTLSConfig(c), redis.DialUseTLS(true))
		log.WithoutContext().Debug("Redis TLS is enabled")
	}

	return &redis.Pool{
		// TODO Fine tune those parameters or make them accessible in config file
		MaxIdle:     10000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return Dial("tcp", cfg.URL(), options...)
		},
	}, nil
}
