package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
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
func NewPool(address string, options ...redis.DialOption) *redis.Pool {
	return &redis.Pool{
		// TODO Fine tune those parameters or make them accessible in config file
		MaxIdle:     10000,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return Dial("tcp", address, options...) },
	}
}
