package keystore

import (
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/aws-secret-manager.git/secretstore"
	"fmt"
	"math/big"
)

// TxSignatureSession holds all the logic allowing the signature of an ethereum transaction
type TxSignatureSession struct {
	secretStore secretstore.SecretStore
	wallet *Wallet
	chain *types.Chain
	tx *ethtypes.Transaction
	signedRaw []byte
	txHash *common.Hash
}

// MakeTxSignature create a new tx signature session from address
func MakeTxSignature(secretStore secretstore.SecretStore) *TxSignatureSession {

	return &TxSignatureSession{
		secretStore: secretStore,
	}
}

// SetWallet sets the wallet to the provided address
func (sess *TxSignatureSession) SetWallet(address *common.Address) error {

	wallet := NewWallet(sess.secretStore)
	err := wallet.Load(address)
	if err != nil {
		return fmt.Errorf("Could not retrieve private key for address : " + err.Error())
	}

	sess.wallet = wallet
	return nil
}

// SetChain is a setter for the chain used in the signed process
func (sess *TxSignatureSession) SetChain(chain *types.Chain) error {
	sess.chain = chain
	return nil
}

// SetTx setter for the Tx to sign in the session
func (sess *TxSignatureSession) SetTx(tx *ethtypes.Transaction) error {
	sess.tx = tx
	return nil
}

// getSigner is internal function that returns an object used during the process
func (sess *TxSignatureSession) getSigner() (ethtypes.Signer, error) {

	var signer ethtypes.Signer
	if sess.chain == nil {
		return nil, fmt.Errorf("Chain has not been set")
	}

	if sess.chain.IsEIP155 {
		// We copy chain ID to ensure pointer can be safely used elsewhere
		id := new(big.Int)
		id.Set(sess.chain.ID)
		signer = ethtypes.NewEIP155Signer(id)

	} else {
		signer = ethtypes.HomesteadSigner{}
	}
	return signer, nil
}

// Run : once all the element of the session have been set, 
// it assigns the signed transaction and the txhash
func (sess *TxSignatureSession) Run() (err error) {

	signer, err := sess.getSigner()
	if err != nil {
		return err
	}

	t, err := ethtypes.SignTx(sess.tx, signer, sess.wallet.GetPriv())
	if err != nil {
		return err
	}

	// Set raw transaction
	signedRaw, err := rlp.EncodeToBytes(t)
	if err != nil {
		return err
	}

	txHash := t.Hash()
	sess.signedRaw = signedRaw
	sess.txHash = &txHash
	return nil
}


