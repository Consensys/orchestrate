package keystore

import (
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/aws-secret-manager.git/secretstore"
	"fmt"
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
	chain *types.Chain, 
	a common.Address, 
	tx *ethtypes.Transaction,
	) (raw []byte, hash *common.Hash, err error) {

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

	return  sess.signedRaw, sess.txHash, nil
}

// SignMsg returns a signed message and its hash
func (s *BaseKeyStore) SignMsg(
	a common.Address, 
	msg string,
	) (rsv []byte, hash *common.Hash, err error) {

		return []byte{}, nil, fmt.Errorf("Not implemented yet")
}

// SignRawHash returns a signed raw hash
func (s *BaseKeyStore) SignRawHash(
	a common.Address, 
	hash []byte,
) (rsv []byte, err error) {

	return []byte{}, fmt.Errorf("Not implemented yet")
}

// GenerateWallet create and stores a new wallet in the vault
func (s* BaseKeyStore) GenerateWallet() (add *common.Address, err error) {

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