package services

import (
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
	//"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
	"github.com/ethereum/go-ethereum/common"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"fmt"
)

// AWSKeyStore implements the services.TxSigner interface
type AWSKeyStore struct {
	client *secretsmanager.SecretsManager
}

//NewAWSKeyStore construct a AWSKeyStore from 
func NewAWSKeyStore(client *secretsmanager.SecretsManager) *AWSKeyStore {
	return &AWSKeyStore{
		client: client,
	}
}

// Sign is the old methods of the keystore to sign a transaction it is perfectly equivalent to SignTx
func (s *AWSKeyStore) Sign(chain *types.Chain, a common.Address, tx *ethtypes.Transaction) (raw []byte, hash *common.Hash, err error) {
	return s.SignTx(chain, a, tx)
}

// SignTx returns a signed transaction. It is perfectly equivalent to SignTx
func (s *AWSKeyStore) SignTx(
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
func (s *AWSKeyStore) SignMsg(
	a common.Address, 
	msg string,
	) (rsv []byte, hash *common.Hash, err error) {

		return []byte{}, nil, fmt.Errorf("Not implemented yet")
}

// SignRawHash returns a signed raw hash
func (s *AWSKeyStore) SignRawHash(
	a common.Address, 
	hash []byte,
) (rsv []byte, err error) {

	return []byte{}, fmt.Errorf("Not implemented yet")
}