package hashicorps

import (
	"github.com/hashicorp/vault/api"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
	"github.com/ethereum/go-ethereum/common"
)

//KeyStore olds the methods of the interfaces KeyStore
type KeyStore struct {
	client *api.Client
}

//NewKeyStore construct a AWSKeyStore from 
func NewKeyStore(client *api.Client) *KeyStore {
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

// Init fetched the root-token
func Unseal() {
	if 
}