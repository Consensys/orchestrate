package mock

import (
	"crypto/ecdsa"
	"fmt"
	"sync"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
)

// KeyStore holds a pool of private keys in memory
type KeyStore struct {
	pKeys   map[string]*ecdsa.PrivateKey
	signers *sync.Map
}

// NewKeyStore creates a new static signer
func NewKeyStore() *KeyStore {
	return &KeyStore{
		signers: &sync.Map{},
		pKeys:   make(map[string]*ecdsa.PrivateKey),
	}
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

func (s *KeyStore) getSigner(c *common.Chain) ethtypes.Signer {
	signer, _ := s.signers.LoadOrStore(c.Id, MakeSigner(c))
	return signer.(ethtypes.Signer)
}

// SignTx sign transaction on context
func (s *KeyStore) SignTx(chain *common.Chain, a ethcommon.Address, tx *ethtypes.Transaction) (raw []byte, hash *ethcommon.Hash, err error) {
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

// SignMsg returns a signed message and its hash
func (s *KeyStore) SignMsg(a ethcommon.Address, msg string) (rsv []byte, hash *ethcommon.Hash, err error) {
	return []byte{}, nil, fmt.Errorf("Not implemented yet")
}

// SignRawHash returns a signed raw hash
func (s *KeyStore) SignRawHash(
	a ethcommon.Address,
	hash []byte,
) (rsv []byte, err error) {

	return []byte{}, fmt.Errorf("Not implemented yet")
}

// GenerateWallet create and stores a new wallet in the vault
func (s *KeyStore) GenerateWallet() (add *ethcommon.Address, err error) {
	prv, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	addr := ethcrypto.PubkeyToAddress(prv.PublicKey)

	// Register private key
	s.pKeys[addr.Hex()] = prv

	return &addr, nil
}

// ImportPrivateKey adds a private key in the vault
func (s *KeyStore) ImportPrivateKey(priv string) (err error) {
	prv, err := ethcrypto.HexToECDSA(priv)
	if err != nil {
		return err
	}
	a := ethcrypto.PubkeyToAddress(prv.PublicKey)
	s.pKeys[a.Hex()] = prv

	return nil
}
