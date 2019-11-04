package redis

import (
	"time"

	remote "github.com/gomodule/redigo/redis"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

// PoolConfig is a place holder to configure the redis client
type PoolConfig struct {
	MaxIdle         int
	MaxActive       int
	MaxConnLifetime time.Duration
	IdleTimeout     time.Duration
	Wait            bool
	URI             string
}

// NewPool creates a new redis pool
func NewPool(conf *PoolConfig, dialFunc DialFunc) *remote.Pool {
	return &remote.Pool{
		// TODO Fine tune those parameters or make them accessible in config file
		MaxIdle:         conf.MaxIdle,
		MaxActive:       conf.MaxActive,
		MaxConnLifetime: conf.MaxConnLifetime,
		IdleTimeout:     conf.IdleTimeout,
		Wait:            conf.Wait,
		Dial:            func() (remote.Conn, error) { return dialFunc("tcp", conf.URI) },
	}
}

// DialFunc is a function alias for function used by the pool to dial redis
type DialFunc func(network, address string, options ...remote.DialOption) (remote.Conn, error)

// Dial is the regular redis dialer
func Dial(network, address string, options ...remote.DialOption) (remote.Conn, error) {
	conn, err := remote.Dial(network, address, options...)
	if err != nil {
		return conn, errors.ConnectionError(err.Error())
	}
	return conn, nil
}

// Conn is a wrapper around a remote.Conn that handles internal errors
type Conn struct{ remote.Conn }

// Close terminates the connection  with the redis store
func (conn *Conn) Close() {
	_ = conn.Conn.Close()
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
func (conn *Conn) Get(key []byte) (byteslice []byte, ok bool, err error) {
	reply, err := conn.Do("GET", key)
	if err != nil {
		return []byte{}, false, err
	}

	if reply == nil {
		// No error is returned if the abi is not stored.
		// This is higher level code's responsibility to deal with it
		return []byte{}, false, nil
	}

	res, err := remote.Bytes(reply, nil)
	if err != nil {
		return []byte{}, false, err
	}

	if len(res) == 0 {
		// No error is returned if the abi is not stored.
		// This is higher level code's responsibility to deal with it
		return []byte{}, false, nil
	}

	return res, true, nil
}

// Set value at a given key in the redis store
func (conn *Conn) Set(key, value []byte) error {
	_, err := conn.Do("SET", key, value)
	return err
}

// RPush appends a stored list with a new element
func (conn *Conn) RPush(key, value []byte) error {
	_, err := conn.Do("RPUSH", key, value)
	return err
}

// LRange returns an entire list stored on Redis
func (conn *Conn) LRange(key []byte) (list [][]byte, ok bool, err error) {
	reply, err := conn.Do("LRANGE", key, 0, -1)
	if err != nil {
		return nil, false, err
	}

	if reply == nil {
		// No error is returned if the abi is not stored.
		// This is higher level code's responsibility to deal with it
		return [][]byte{}, false, nil
	}

	res, err := remote.ByteSlices(reply, nil)
	if err != nil {
		return [][]byte{}, false, err
	}

	if len(res) == 0 {
		// No error is returned if the abi is not stored.
		// This is higher level code's responsibility to deal with it
		return [][]byte{}, false, nil
	}

	return res, true, nil
}

// Send writes a request in the redis buffer
func (conn *Conn) Send(commandName string, args ...interface{}) error {
	err := conn.Conn.Send(commandName, args...)
	if err != nil {
		return errors.ConnectionError(err.Error())
	}
	return nil
}

// Flush a batch of requests to the remote redis
func (conn *Conn) Flush() error {
	err := conn.Conn.Flush()
	if err != nil {
		return errors.ConnectionError(err.Error())
	}
	return nil
}

// ReceiveBytes wait for the connection respond with a []byte
func (conn *Conn) ReceiveBytes() (bytes []byte, ok bool, err error) {
	reply, err := conn.Conn.Receive()
	if err != nil {
		return []byte{}, false, err
	}

	if reply == nil {
		// No error is returned if the abi is not stored.
		// This is higher level code's responsibility to deal with it
		return []byte{}, false, nil
	}

	res, err := remote.Bytes(reply, nil)
	if err != nil {
		return []byte{}, false, err
	}

	if len(res) == 0 {
		// No error is returned if the abi is not stored.
		// This is higher level code's responsibility to deal with it
		return []byte{}, false, nil
	}

	return res, true, nil
}

// ReceiveByteSlices returns a pipelined [][]byte result
func (conn *Conn) ReceiveByteSlices() (byteSlices [][]byte, ok bool, err error) {
	reply, err := conn.Conn.Receive()
	if err != nil {
		return [][]byte{}, false, err
	}

	if reply == nil {
		// No error is returned if the abi is not stored.
		// This is higher level code's responsibility to deal with it
		return [][]byte{}, false, nil
	}

	res, err := remote.ByteSlices(reply, nil)
	if err != nil {
		return [][]byte{}, false, err
	}

	if len(res) == 0 {
		// No error is returned if the abi is not stored.
		// This is higher level code's responsibility to deal with it
		return [][]byte{}, false, nil
	}

	return res, true, nil
}

// ReceiveCheck returns an error if the pipelined result is an error
func (conn *Conn) ReceiveCheck() error {
	_, err := conn.Conn.Receive()
	return err
}

// SendGet returns a stored byteslice stored on redis, but does not flush
func (conn *Conn) SendGet(key []byte) error {
	return conn.Send("GET", key)
}

// SendSet value at a given key in the redis store, but does not flush
func (conn *Conn) SendSet(key, value []byte) error {
	return conn.Send("SET", key, value)
}

// SendRPush appends a stored list with a new element, but does not flush
func (conn *Conn) SendRPush(key, value []byte) error {
	return conn.Send("RPUSH", key, value)
}

// SendLRange returns an entire list stored on Redis, but does not flush
func (conn *Conn) SendLRange(key []byte) error {
	return conn.Send("LRANGE", key, 0, -1)
}
