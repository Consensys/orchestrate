package redis

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const byteCodeHashPrefix = "ByteCodeHashPrefix"

// ByteCodeHashModel is a zero object gathering methods to look up a abis in redis
type ByteCodeHashModel struct{}

// ByteCodeHash returns is sugar to return an abi object
var ByteCodeHash = &ByteCodeHashModel{}

// Key serializes a lookup key for an ABI stored on redis
func (*ByteCodeHashModel) Key(name string, tag string) []byte {
	prefixBytes := []byte(byteCodeHashPrefix)
	// Allocate memory to build the key
	res := make([]byte, 0, len(prefixBytes)+len(name)+len(tag))
	res = append(res, prefixBytes...)
	res = append(res, name...)
	res = append(res, tag...)
	return res
}

// Get returns a serialized contract from its corresponding bytecode hash
func (b *ByteCodeHashModel) Get(conn *Conn, name string, tag string) (ethcommon.Hash, error) {
	bytes, err := conn.Get(b.Key(name, tag))
	if err != nil {
		return ethcommon.Hash{}, err
	}

	if len(bytes) == 0 {
		return ethcommon.Hash{}, nil
	}

	return ethcommon.BytesToHash(bytes), nil
}

// Set stores an abi object in the registry
func (b *ByteCodeHashModel) Set(conn *Conn, name string, tag string, byteCodeHash ethcommon.Hash) error {
	return conn.Set(b.Key(name, tag), byteCodeHash[:])
}

// FromAccountInstance
func (b *ByteCodeHashModel) FromAccountInstance()