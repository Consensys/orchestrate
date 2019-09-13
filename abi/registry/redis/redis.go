package redis

import (
	"time"
	remote "github.com/gomodule/redigo/redis"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
)

// PoolConfig is a place holder to configure the redis client
type PoolConfig struct {
	MaxIdle int
	MaxActive int
	MaxConnLifetime time.Duration
	IdleTimeout time.Duration
	Wait bool
	URI string
}

// NewPool creates a new redis pool
func NewPool(conf *PoolConfig) *remote.Pool {
	return &remote.Pool{
		// TODO Fine tune those parameters or make them accessible in config file
		MaxIdle:     conf.MaxIdle,
		MaxActive:   conf.MaxActive,
		MaxConnLifetime: conf.MaxConnLifetime,
		IdleTimeout: conf.IdleTimeout,
		Wait: conf.Wait,
		Dial:        func() (remote.Conn, error) { return dial("tcp", conf.URI) },
	}
}

func dial(network, address string, options ...remote.DialOption) (remote.Conn, error) {
	conn, err := remote.Dial(network, address, options...)
	if err != nil {
		return conn, errors.ConnectionError(err.Error())
	}
	return conn, nil
}

// Conn is a wrapper around a remote.Conn that handles internal errors
type Conn struct{ remote.Conn }

// Close terminates the connexion with the redis store
func (conn *Conn) Close() {
	conn.Conn.Close()
}

// Do sends a commands to the remote Redis instance
func (conn *Conn) Do(commandName string, args ...interface{}) (interface{}, error) {
	reply, err := conn.Conn.Do(commandName, args...)
	if err != nil {
		return reply, errors.ConnectionError(err.Error())
	}
	return reply, nil
}

// Get returns a stored byteslice stored on redis
func (conn *Conn) Get(key []byte) ([]byte, bool, error) {
	reply, err := conn.Do("GET", key)
	if err != nil {
		return []byte{}, false, err
	}

	if reply == nil {
		// No error is returned is returned if the abi is not stored.
		// This is higher level code's responsibility to deal with it
		return []byte{}, false, nil
	}

	res, err := remote.Bytes(reply, nil)
	if err != nil {
		return []byte{}, false, err
	}

	return res, true, nil
}

// Set value at a given key in the redis store
func (conn *Conn) Set(key, value []byte) error {
	_, err := conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	return nil
}

// LPush appends a stored list with a new element
func (conn *Conn) LPush(key, value []byte) (error) {
	_, err := conn.Do("LPUSH", key, value)
	if err != nil {
		return err
	}

	return nil
}

// LRange returns an entire list stored on Redis
func (conn *Conn) LRange(key []byte) ([][]byte, bool, error) {
	reply, err := conn.Do("LRANGE", key, 0, -1)
	if err != nil {
		return nil, false, err
	}

	if reply == nil {
		// No error is returned is returned if the abi is not stored.
		// This is higher level code's responsibility to deal with it
		return [][]byte{}, false, nil
	}

	res, err := remote.ByteSlices(reply, nil)
	if err != nil {
		return [][]byte{}, false, err
	}

	return res, true, nil
}

// Send writes a request in the redis buffer
func (conn *Conn) Send(commandName string, args ...interface{}) error {
	return conn.Conn.Send(commandName, args)
}

// Flush a batch of requests to the remote redis
func (conn *Conn) Flush() error {
	return conn.Conn.Flush()
}

// ReceiveBytes wait for the connexion respond with a []byte
func (conn *Conn) ReceiveBytes() ([]byte, bool, error) {
	reply, err := conn.Conn.Receive()
	if err != nil {
		return []byte{}, false, err
	}

	if reply == nil {
		// No error is returned is returned if the abi is not stored.
		// This is higher level code's responsibility to deal with it
		return []byte{}, false, nil
	}

	res, err := remote.Bytes(reply, nil)
	if err != nil {
		return []byte{}, false, err
	}

	return res, true, nil
}

// ReceiveByteSlices returns a pipelined [][]byte result
func (conn *Conn) ReceiveByteSlices() ([][]byte, bool, error) {
	reply, err := conn.Conn.Receive()
	if err != nil {
		return [][]byte{}, false, err
	}

	if reply == nil {
		// No error is returned is returned if the abi is not stored.
		// This is higher level code's responsibility to deal with it
		return [][]byte{}, false, nil
	}

	res, err := remote.ByteSlices(reply, nil)
	if err != nil {
		return [][]byte{}, false, err
	}

	return res, true, nil
}

// ReceiveCheck returns an error if the pipelined result is an error
func (conn *Conn) ReceiveCheck() (error) {
	_, err := conn.Conn.Receive()
	return err
}

// SendGet returns a stored byteslice stored on redis, but does not flush
func (conn *Conn) SendGet(key []byte) (error) {
	return conn.Send("GET", key)
}

// SendSet value at a given key in the redis store, but does not flush
func (conn *Conn) SendSet(key, value []byte) error {
	return conn.Send("SET", key, value)
}

// SendLPush appends a stored list with a new element, but does not flush
func (conn *Conn) SendLPush(key, value []byte) (error) {
	return conn.Send("LPUSH", key, value)
}

// SendLRange returns an entire list stored on Redis, but does not flush
func (conn *Conn) SendLRange(key []byte) (error) {
	return conn.Send("LRANGE", key, 0, -1)
}
