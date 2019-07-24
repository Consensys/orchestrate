package session

import (
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/types"
	"golang.org/x/crypto/sha3"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/multi-vault.git/keystore/wallet"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/multi-vault.git/secretstore/services"
)

// TxSignatureSession holds all the logic allowing the signature of an ethereum transaction
type TxSignatureSession struct {
	secretStore services.SecretStore
	wallet      *wallet.Wallet
	chain       *chain.Chain
	tx          *ethtypes.Transaction
	privateArgs *types.PrivateArgs
	Raw         []byte
	Hash        *ethcommon.Hash
}

// MakeTxSignature create a new tx signature session from address
func MakeTxSignature(secretStore services.SecretStore) *TxSignatureSession {
	return &TxSignatureSession{
		secretStore: secretStore,
	}
}

// SetWallet sets the wallet to the provided address
func (sess *TxSignatureSession) SetWallet(address *ethcommon.Address) error {

	w := wallet.NewWallet(sess.secretStore)
	err := w.Load(address)
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	sess.wallet = w
	return nil
}

// SetChain is a setter for the chain used in the signed process
func (sess *TxSignatureSession) SetChain(netChain *chain.Chain) {
	sess.chain = netChain
}

// SetTx setter for the Tx to sign in the session
func (sess *TxSignatureSession) SetTx(tx *ethtypes.Transaction) {
	sess.tx = tx
}

// SetPrivateArgs setter for the arguments of the Tx to sign in the session
func (sess *TxSignatureSession) SetPrivateArgs(privateArgs *types.PrivateArgs) {
	sess.privateArgs = privateArgs
}

// getSigner is internal function that returns an object used during the process
func (sess *TxSignatureSession) getSigner() (ethtypes.Signer, error) {
	var signer ethtypes.Signer
	if sess.chain == nil {
		return nil, errors.FailedPreconditionError("chain not set").ExtendComponent(component)
	}

	// We copy chain ID to ensure pointer can be safely used elsewhere
	signer = ethtypes.NewEIP155Signer(new(big.Int).Set(sess.chain.ID()))

	return signer, nil
}

// Run : once all the element of the session have been set,
// it assigns the signed transaction and the txhash
func (sess *TxSignatureSession) Run() (err error) {
	signer, err := sess.getSigner()
	if err != nil {
		return err
	}

	t, err := ethtypes.SignTx(sess.tx, signer, sess.wallet.Priv())
	if err != nil {
		return errors.CryptoOperationError(err.Error()).SetComponent(component)
	}

	// Set raw transaction
	signedRaw, err := rlp.EncodeToBytes(t)
	if err != nil {
		return errors.InternalError(err.Error()).SetComponent(component)
	}

	txHash := t.Hash()
	sess.Raw, sess.Hash = signedRaw, &txHash

	return nil
}

// SignPrivateEEATransaction : once all the element of the session have been set,
// it assigns the signed transaction and the txhash
// Calculates a transaction hash for sending private transactions for EEA extension
func (sess *TxSignatureSession) SignPrivateEEATransaction() (err error) {
	// Should sign using the EIP155 signer
	signer, err := sess.getSigner()
	if err != nil {
		return err
	}

	hash := privateTxHash(sess)

	sig, err := crypto.Sign(hash[:], sess.wallet.Priv())
	if err != nil {
		return err
	}

	t, err := sess.tx.WithSignature(signer, sig)
	if err != nil {
		return err
	}

	rplEncoding := encodePrivateTx(t, sess)

	sess.Raw = rplEncoding
	sess.Hash = &hash
	return nil
}

func privateTxHash(sess *TxSignatureSession) ethcommon.Hash {
	hash := rlpHash([]interface{}{
		sess.tx.Nonce(),
		sess.tx.GasPrice(),
		sess.tx.Gas(),
		sess.tx.To(),
		sess.tx.Value(),
		sess.tx.Data(),
		sess.chain.ID(),
		uint(0),
		uint(0),
		sess.privateArgs.PrivateFrom,
		sess.privateArgs.PrivateFor,
		sess.privateArgs.PrivateTxType,
	})
	return hash
}

func encodePrivateTx(t *ethtypes.Transaction, sess *TxSignatureSession) []byte {
	v, r, s := t.RawSignatureValues()
	rplEncoding, _ := rlpEncode([]interface{}{
		sess.tx.Nonce(),
		sess.tx.GasPrice(),
		sess.tx.Gas(),
		sess.tx.To(),
		sess.tx.Value(),
		sess.tx.Data(),
		v,
		r,
		s,
		sess.privateArgs.PrivateFrom,
		sess.privateArgs.PrivateFor,
		sess.privateArgs.PrivateTxType,
	})
	return rplEncoding
}

// SignPrivateTesseraTransaction : once all the element of the session have been set,
// it assigns the signed transaction and the txhash
// Signs a transaction for Tessera private enclave
func (sess *TxSignatureSession) SignPrivateTesseraTransaction() (err error) {
	// Transactions for Tessera should be signed using Homestead signer
	t, err := ethtypes.SignTx(sess.tx, ethtypes.HomesteadSigner{}, sess.wallet.Priv())
	if err != nil {
		return err
	}

	// Get raw signed transaction
	signedRaw, err := rlp.EncodeToBytes(t)
	if err != nil {
		return err
	}

	txHash := hashTesseraTx(t, sess)
	sess.Raw = signedRaw
	sess.Hash = &txHash
	return nil
}

func hashTesseraTx(t *ethtypes.Transaction, sess *TxSignatureSession) ethcommon.Hash {
	v, r, s := t.RawSignatureValues()
	privateV := calculatePrivateV(v)
	txHash := rlpHash([]interface{}{
		sess.tx.Nonce(),
		sess.tx.GasPrice(),
		sess.tx.Gas(),
		sess.tx.To(),
		sess.tx.Value(),
		sess.tx.Data(),
		privateV,
		r,
		s,
	})
	return txHash
}

// Tessera replaces the "V" value with 37 or 38 if it was 27 or 28 respectively
func calculatePrivateV(v *big.Int) *big.Int {
	if v.Int64() == 27 {
		return big.NewInt(37)
	}
	return big.NewInt(38)
}

func rlpHash(object interface{}) (hash ethcommon.Hash) {
	hashAlgo := sha3.NewLegacyKeccak256()
	_ = rlp.Encode(hashAlgo, object)
	hashAlgo.Sum(hash[:0])
	return hash
}

func rlpEncode(object interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(object)
}
