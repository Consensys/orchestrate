package base

import (
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/keystore/session"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/keystore/wallet"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/secretstore/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
)

// KeyStore olds the methods of the interfaces BaseKeyStore
type KeyStore struct {
	SecretStore services.SecretStore
}

// NewKeyStore construct a BaseKeyStore from a client
func NewKeyStore(secretStore services.SecretStore) *KeyStore {
	return &KeyStore{
		SecretStore: secretStore,
	}
}

// SignTx returns a signed transaction. It is perfectly equivalent to SignTx
func (s *KeyStore) SignTx(chain *common.Chain, a ethcommon.Address, tx *ethtypes.Transaction) (raw []byte, hash *ethcommon.Hash, err error) {
	// Creates a new signing session
	sess := session.MakeTxSignature(s.SecretStore)
	err = sess.SetWallet(&a)
	if err != nil {
		return []byte{}, nil, err
	}
	err = sess.SetChain(chain)
	if err != nil {
		return []byte{}, nil, err
	}
	err = sess.SetTx(tx)
	if err != nil {
		return []byte{}, nil, err
	}

	// Run signing session
	err = sess.Run()
	if err != nil {
		return []byte{}, nil, err
	}

	return sess.Raw, sess.Hash, nil
}

// SignMsg returns a signed message and its hash
func (s *KeyStore) SignMsg(a ethcommon.Address, msg string) (rsv []byte, hash *ethcommon.Hash, err error) {
	return []byte{}, nil, fmt.Errorf("not implemented yet")
}

// SignRawHash returns a signed raw hash
func (s *KeyStore) SignRawHash(
	a ethcommon.Address,
	hash []byte,
) (rsv []byte, err error) {

	return []byte{}, fmt.Errorf("not implemented yet")
}

// GenerateWallet create and stores a new wallet in the vault
func (s *KeyStore) GenerateWallet() (add *ethcommon.Address, err error) {
	w := wallet.NewWallet(s.SecretStore)
	err = w.Generate()
	if err != nil {
		return nil, err
	}

	err = w.Store()
	if err != nil {
		return nil, err
	}

	return w.Address(), nil
}

// ImportPrivateKey adds a private key in the vault
// TODO: this is Unsafe and should be removed soon
func (s *KeyStore) ImportPrivateKey(priv string) (err error) {

	w := wallet.NewWallet(s.SecretStore)
	err = w.FromPrivateKey(priv)
	if err != nil {
		return err
	}

	err = w.Store()
	if err != nil {
		return err
	}

	return nil
}
