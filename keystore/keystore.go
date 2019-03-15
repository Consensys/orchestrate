package keystore

import (
	"github.com/hashicorp/vault/api"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"fmt"
)

//BasicKeyStore olds the methods of the interfaces BasicKeyStore
type BasicKeyStore struct {
	SecretStore *SecretStore
}

//NewBasicKeyStore construct a BasicKeyStore from a client 
func NewBasicKeyStore(secretStore *SecretStore) *BasicKeyStore {
	return &BasicKeyStore{
		SecretStore: secretStore
	}
}

// SignTx returns a signed transaction. It is perfectly equivalent to SignTx
func (s *BasicKeyStore) SignTx(
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
func (s *BasicKeyStore) SignMsg(
	a common.Address, 
	msg string,
	) (rsv []byte, hash *common.Hash, err error) {

		return []byte{}, nil, fmt.Errorf("Not implemented yet")
}

// SignRawHash returns a signed raw hash
func (s *BasicKeyStore) SignRawHash(
	a common.Address, 
	hash []byte,
) (rsv []byte, err error) {

	return []byte{}, fmt.Errorf("Not implemented yet")
}

// GenerateWallet create and stores a new wallet in the vault
func (s* BasicKeyStore) GenerateWallet() (add *common.Address, err error) {

	wal, err := NewWallet(sess.secretStore).Generate()
	if err != nil {
		return nil, err
	}

	err = wal.Store()
	if err != nil {
		return nil, err
	}

	return wal.GetAddress(), nil
}