package redis

import (
	ethcommon "github.com/ethereum/go-ethereum/common"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
)

const deployedByteCodeHashPrefix = "DeployedByteCodeHashPrefix"

// DeployedByteCodeHashModel is a zero object gathering methods to look up a abis in redis
type DeployedByteCodeHashModel struct{}

// DeployedByteCodeHash returns is sugar to return an abi object
var DeployedByteCodeHash = &DeployedByteCodeHashModel{}

// Key serializes a lookup key for an ABI stored on redis
func (*DeployedByteCodeHashModel) Key(chain *chain.Chain, address *ethereum.Account) []byte {
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
func (b *DeployedByteCodeHashModel) Get(conn *Conn, chain *chain.Chain, address *ethereum.Account) (ethcommon.Hash, error) {
	bytes, err := conn.Get(b.Key(chain, address))
	if err != nil {
		return ethcommon.Hash{}, err
	}

	if len(bytes) == 0 {
		return ethcommon.Hash{}, nil
	}

	return ethcommon.BytesToHash(bytes), nil
}

// Set stores an abi object in the registry
func (b *DeployedByteCodeHashModel) Set(conn *Conn, chain *chain.Chain, address *ethereum.Account, byteCodeHash ethcommon.Hash) error {
	return conn.Set(b.Key(chain, address), byteCodeHash[:])
}
