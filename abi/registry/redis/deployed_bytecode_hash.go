package redis

import (
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
)

const deployedByteCodeHashPrefix = "DeployedByteCodeHashPrefix"

// DeployedByteCodeHashModel is a zero object gathering methods to look up a abis in redis
type DeployedByteCodeHashModel struct{}

// DeployedByteCodeHash returns is sugar to manage deployed bytecode hashes
var DeployedByteCodeHash = &DeployedByteCodeHashModel{}

// Key serializes a lookup key for deployed bytecode hash stored on redis
func (*DeployedByteCodeHashModel) Key(chain fmt.Stringer, address *ethereum.Account) []byte {
	prefixBytes := []byte(deployedByteCodeHashPrefix)
	chainBytes := chain.String()
	addressBytes := address.GetRaw()
	// Allocate memory to build the key
	res := make([]byte, 0, len(prefixBytes)+len(chainBytes)+len(addressBytes))
	res = append(res, prefixBytes...)
	res = append(res, chainBytes...)
	res = append(res, addressBytes...)
	return res
}

// Get returns a serialized contract from its corresponding bytecode hash
func (b *DeployedByteCodeHashModel) Get(conn *Conn, chain fmt.Stringer, address *ethereum.Account) (ethcommon.Hash, bool, error) {
	bytes, ok, err := conn.Get(b.Key(chain, address))
	if err != nil || !ok {
		return ethcommon.Hash{}, false, err
	}

	return ethcommon.BytesToHash(bytes), true, nil
}

// Set stores a deployed bytecode hash in the registry
func (b *DeployedByteCodeHashModel) Set(conn *Conn, chain fmt.Stringer, address *ethereum.Account, byteCodeHash ethcommon.Hash) error {
	return conn.Set(b.Key(chain, address), byteCodeHash[:])
}
