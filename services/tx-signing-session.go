package services

import (
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
	//"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
	"github.com/ethereum/go-ethereum/common"
	aws "gitlab.com/ConsenSys/client/fr/core-stack/infra/aws-secret-manager.git/aws"
	"fmt"
)

// TxSignatureSession holds all the logic allowing the signature of an ethereum transaction
type TxSignatureSession struct { 
	wallet *aws.Wallet
	chain *types.Chain
	tx *ethtypes.Transaction
	signedRaw []byte
	txHash *common.Hash
}

// MakeTxSignature create a new tx signature session from address
func MakeTxSignature() *TxSignatureSession {

	return &TxSignatureSession{
		wallet: nil,
		chain: nil,
		tx: nil,
		txHash: nil,
		signedRaw: nil,
	}
}

// SetWallet sets the wallet to the provided address
func (sess *TxSignatureSession) SetWallet(address *common.Address) error {

	wallet, err := aws.GetWallet(address)
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
		return fmt.Errorf("Chain has not been set")
	}

	if sess.chain.IsEIP155 {
		// We copy chain ID to ensure pointer can be safely used elsewhere
		id := new(big.Int)
		id.Set(sess.chain.ID)
		signer = ethtypes.NewEIP155Signer(id)

	} else {
		signer = ethtypes.HomesteadSigner{}
	}
	return signer
}

// Run : once all the element of the session have been set, 
// it assigns the signed transaction and the txhash
func (sess *TxSignatureSession) Run() error {

	signer := sess.getSigner()

	t, err := ethtypes.SignTx(tx, signer, sess.Wallet.priv)
	if err != nil {
		return err
	}

	// Set raw transaction
	signedRaw, err = rlp.EncodeToBytes(t)
	if err != nil {
		// TODO: handle error
		return err
	}

	sess.signedRaw = signedRaw
	sess.txHash = t.Hash()
	return nil
}


