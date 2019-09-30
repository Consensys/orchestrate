package redis

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/protobuf/proto"

	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/types/abi"
)

const artifactPrefix = "ArtifactPrefix"

// ArtifactModel is a zero object gathering methods to look up a abis in redis
type ArtifactModel struct{}

// Artifact returns a model object
var Artifact = &ArtifactModel{}

// Key serializes a lookup key for contract artifacts stored in redis
func (*ArtifactModel) Key(byteCodeHash ethcommon.Hash) []byte {
	// Allocate memory to build the key
	res := make([]byte, 0, len(artifactPrefix)+len(byteCodeHash))
	res = append(res, artifactPrefix...)
	res = append(res, byteCodeHash[:]...)
	return res
}

// Get returns a serialized contract from its corresponding bytecode hash
func (a *ArtifactModel) Get(conn *Conn, byteCodeHash ethcommon.Hash) (*abi.Contract, bool, error) {
	marshalledArtifact, ok, err := conn.Get(a.Key(byteCodeHash))
	if err != nil || !ok {
		// The check is redundant because !ok => err != nil
		return nil, false, err
	}

	contract := &abi.Contract{}
	if err = proto.Unmarshal(marshalledArtifact, contract); err != nil {
		return nil, false, errors.FromError(err).ExtendComponent(component)
	}

	return contract, true, nil
}

// Set stores an artifact in the registry
func (a *ArtifactModel) Set(conn *Conn, byteCodeHash ethcommon.Hash, contract proto.Message) error {
	marshalledArtifact, err := proto.Marshal(contract)
	if err != nil {
		return err
	}

	return conn.Set(a.Key(byteCodeHash), marshalledArtifact)
}
