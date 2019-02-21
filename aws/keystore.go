package aws

import (
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
	//"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
	"github.com/ethereum/go-ethereum/common"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"fmt"
)

// KeyStore implements the services.TxSigner interface
type KeyStore struct {
	client *secretsmanager.SecretsManager
}

//NewKeyStore construct a KeyStore from 
func NewKeyStore(client *secretsmanager.SecretsManager) *KeyStore {
	return &KeyStore{
		client: client,
	}
}

// Sign is the old methods of the keystore to sign a transaction it is perfectly equivalent to SignTx
func (s *KeyStore) Sign(chain *types.Chain, a common.Address, tx *ethtypes.Transaction) (raw []byte, hash *common.Hash, err error) {
	return s.SignTx(chain, a, tx)
}

// SignTx returns a signed transaction. It is perfectly equivalent to SignTx
func (s *KeyStore) SignTx(
	chain *types.Chain, 
	a common.Address, 
	tx *ethtypes.Transaction,
	) (raw []byte, hash *common.Hash, err error) {

	sess := MakeTxSignature(s.client)
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
func (s *KeyStore) SignMsg(
	a common.Address, 
	msg string,
	) (rsv []byte, hash *common.Hash, err error) {

		return []byte{}, nil, fmt.Errorf("Not implemented yet")
}

// SignRawHash returns a signed raw hash
func (s *KeyStore) SignRawHash(
	a common.Address, 
	hash []byte,
) (rsv []byte, err error) {

	return []byte{}, fmt.Errorf("Not implemented yet")
}

// GenerateWallet returns an ethereum address, corresponding to a newly created wallet
func (s* KeyStore) GenerateWallet() (add *common.Address, err error) {

	wal, err := GenerateWallet(s.client)
	if err != nil {
		return nil, err
	}

	return wal.GetAddress(), nil
}