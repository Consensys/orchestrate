package signature

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
)

type ethECDSA struct{}

// EthECDSA is an implementation of ethereum's digital signature algorithm
var (
	EthECDSA = ethECDSA{}
)

func (e ethECDSA) Sign(digest []byte, priv *ecdsa.PrivateKey) ([]byte, error) {
	rsv, err := crypto.Sign(digest, priv)
	if err != nil {
		return []byte{}, errors.InternalError(err.Error()).SetComponent(component)
	}
	// crypto.Sign returns the signature in the form (r, s, v)
	// But we need to format it as v + 27, r, s for the precompile Ecrecover to work
	vrs := make([]byte, 65)
	// The conversion uint8 -> byte is implicit and the linter mark it as unnecessary to specify
	vrs[0] = uint8(int(rsv[64])) + 27
	// Copy r and s in a single command as their order do not change
	copy(vrs[1:], rsv[:64])

	return vrs, nil
}
