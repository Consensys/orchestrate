package redis

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/protobuf/proto"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/abi"
)

const artifactPrefix = "ArtifactPrefix"

// ArtifactModel is a zero object gathering methods to look up a abis in redis
type ArtifactModel struct{}

// Artifact returns a model object
var Artifact = &ArtifactModel{}

// Key serializes a lookup key for an ABI stored on redis
func (*ArtifactModel) Key(byteCodeHash ethcommon.Hash) []byte {
	// Allocate memory to build the key
	res := make([]byte, 0, len(artifactPrefix)+len(byteCodeHash))
	res = append(res, artifactPrefix...)
	res = append(res, byteCodeHash[:]...)
	return res
}

// Get returns a serialized contract from its corresponding bytecode hash
func (a *ArtifactModel) Get(conn *Conn, byteCodeHash ethcommon.Hash) (*abi.Contract, error) {
	marshalledArtifact, err := conn.Get(a.Key(byteCodeHash))
	if err != nil {
		return nil, err
	}

	contract := &abi.Contract{}
	if err = proto.Unmarshal(marshalledArtifact, contract); err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return contract, nil
}

// Set stores an abi object in the registry
func (a *ArtifactModel) Set(conn *Conn, byteCodeHash ethcommon.Hash, contract *abi.Contract) error {
	marshalledArtifact, err := proto.Marshal(contract)
	if err != nil {
		return err
	}

	return conn.Set(a.Key(byteCodeHash), marshalledArtifact)
}
