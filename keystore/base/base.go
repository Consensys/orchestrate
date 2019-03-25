package base

import (
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/keystore/session"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/keystore/wallet"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/secretstore"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
)

// KeyStore olds the methods of the interfaces BaseKeyStore
type KeyStore struct {
	SecretStore secretstore.SecretStore
}

//NewKeyStore construct a BaseKeyStore from a client
func NewKeyStore(secretStore secretstore.SecretStore) *KeyStore {
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
	sess.SetChain(chain)
	sess.SetTx(tx)

	// Run signing session
	err = sess.Run()
	if err != nil {
		return []byte{}, nil, err
	}

	return sess.Raw, sess.Hash, nil
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
	wallet := wallet.NewWallet(s.SecretStore)
	err = wallet.Generate()
	if err != nil {
		return nil, err
	}

	err = wallet.Store()
	if err != nil {
		return nil, err
	}

	return wallet.Address(), nil
}

// ImportPrivateKey adds a private key in the vault
// TODO: this is Unsafe and should be removed soon
func (s *KeyStore) ImportPrivateKey(priv string) (err error) {

	wallet := wallet.NewWallet(s.SecretStore)
	err = wallet.FromPrivateKey(priv)
	if err != nil {
		return err
	}

	err = wallet.Store()
	if err != nil {
		return err
	}

	return nil
}
