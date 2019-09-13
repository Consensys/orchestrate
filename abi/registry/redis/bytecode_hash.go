package redis

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const byteCodeHashPrefix = "ByteCodeHashPrefix"

// ByteCodeHashModel is a zero object gathering methods to look up a abis in redis
type ByteCodeHashModel struct{}

// ByteCodeHash returns is sugar to manage bytecode hashes
var ByteCodeHash = &ByteCodeHashModel{}

// Key serializes a lookup key for bytecode hash stored on redis
func (*ByteCodeHashModel) Key(name, tag string) []byte {
	prefixBytes := []byte(byteCodeHashPrefix)
	// Allocate memory to build the key
	res := make([]byte, 0, len(prefixBytes)+len(name)+len(tag))
	res = append(res, prefixBytes...)
	res = append(res, name...)
	res = append(res, tag...)
	return res
}

// Get returns a serialized contract from its corresponding bytecode hash
func (b *ByteCodeHashModel) Get(conn *Conn, name, tag string) (ethcommon.Hash, bool, error) {
	bytes, ok, err := conn.Get(b.Key(name, tag))
	if err != nil || !ok {
		return ethcommon.Hash{}, false, err
	}

	return ethcommon.BytesToHash(bytes), true, nil
}

// Set stores a bytecode hash in the registry
func (b *ByteCodeHashModel) Set(conn *Conn, name, tag string, byteCodeHash ethcommon.Hash) error {
	return conn.Set(b.Key(name, tag), byteCodeHash[:])
}
