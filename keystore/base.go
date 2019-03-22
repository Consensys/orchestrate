package keystore

import (
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/secretstore"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
)

//BaseKeyStore olds the methods of the interfaces BaseKeyStore
type BaseKeyStore struct {
	SecretStore secretstore.SecretStore
}

//NewBaseKeyStore construct a BaseKeyStore from a client
func NewBaseKeyStore(secretStore secretstore.SecretStore) *BaseKeyStore {

	return &BaseKeyStore{
		SecretStore: secretStore,
	}
}

// SignTx returns a signed transaction. It is perfectly equivalent to SignTx
func (s *BaseKeyStore) SignTx(
	chain *common.Chain,
	a ethcommon.Address,
	tx *ethtypes.Transaction,
) (raw []byte, hash *ethcommon.Hash, err error) {

	sess := MakeTxSignature(s.SecretStore)
	err = sess.SetWallet(&a)
	if err != nil {
		return []byte{}, nil, err
	}

	sess.SetChain(chain)
	sess.SetTx(tx)

	err = sess.Run()
	if err != nil {
		return []byte{}, nil, err
	}

	return sess.signedRaw, sess.txHash, nil
}

// SignMsg returns a signed message and its hash
func (s *BaseKeyStore) SignMsg(
	a ethcommon.Address,
	msg string,
) (rsv []byte, hash *ethcommon.Hash, err error) {

	return []byte{}, nil, fmt.Errorf("Not implemented yet")
}

// SignRawHash returns a signed raw hash
func (s *BaseKeyStore) SignRawHash(
	a ethcommon.Address,
	hash []byte,
) (rsv []byte, err error) {

	return []byte{}, fmt.Errorf("Not implemented yet")
}

// GenerateWallet create and stores a new wallet in the vault
func (s *BaseKeyStore) GenerateWallet() (add *ethcommon.Address, err error) {

	wallet := NewWallet(s.SecretStore)
	err = wallet.Generate()
	if err != nil {
		return nil, err
	}

	err = wallet.Store()
	if err != nil {
		return nil, err
	}

	return wallet.GetAddress(), nil
}

// ImportPrivateKey adds a private key in the vault
func (s *BaseKeyStore) ImportPrivateKey(priv string) (err error) {

	wallet := NewWallet(s.SecretStore)
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

// RegisterPkeys allow to register private keys on key store
func (s *BaseKeyStore) RegisterPkeys(pkeys []string) (err error) {
	for _, pkey := range pkeys {
		err = s.ImportPrivateKey(pkey)
		if err != nil {
			return err
		}
	}
	return nil
}
