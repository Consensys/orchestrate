package ethereum

import (
	"crypto/ecdsa"
	"fmt"
	"sync"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
)

// StaticSigner holds a pool of private keys in memory
type StaticSigner struct {
	pKeys   map[string]*ecdsa.PrivateKey
	signers *sync.Map
}

// NewStaticSigner creates a new static signer
func NewStaticSigner(pKeys []string) *StaticSigner {
	s := &StaticSigner{
		signers: &sync.Map{},
		pKeys:   make(map[string]*ecdsa.PrivateKey),
	}
	for _, pKey := range pKeys {
		prv, err := ethcrypto.HexToECDSA(pKey)
		if err != nil {
			panic(err)
		}
		a := ethcrypto.PubkeyToAddress(prv.PublicKey)
		s.pKeys[a.Hex()] = prv
	}
	return s
}

// MakeSigner creates a new signer
func MakeSigner(c *common.Chain) ethtypes.Signer {
	var signer ethtypes.Signer
	if c.IsEIP155 {
		// We copy chain ID to ensure pointer can be safely used elsewhere
		signer = ethtypes.NewEIP155Signer(c.ID())
	} else {
		signer = ethtypes.HomesteadSigner{}
	}
	return signer
}

func (s *StaticSigner) getSigner(c *common.Chain) ethtypes.Signer {
	signer, _ := s.signers.LoadOrStore(c.Id, MakeSigner(c))
	return signer.(ethtypes.Signer)
}

// Sign sign transaction on context
func (s *StaticSigner) Sign(chain *common.Chain, a ethcommon.Address, tx *ethtypes.Transaction) (raw []byte, hash *ethcommon.Hash, err error) {
	signer := s.getSigner(chain)

	prv, ok := s.pKeys[a.Hex()]
	if !ok {
		return []byte{}, nil, fmt.Errorf("No Private Key for account %q", a.Hex())
	}

	t, err := ethtypes.SignTx(tx, signer, prv)
	if err != nil {
		return []byte{}, nil, err
	}

	// Set raw transaction
	raw, err = rlp.EncodeToBytes(t)
	if err != nil {
		// TODO: handle error
		return []byte{}, nil, err
	}
	h := t.Hash()

	return raw, &h, nil
}
