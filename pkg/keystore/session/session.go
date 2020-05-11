package session

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/account"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore/crypto"
)

const component = "keystore.session"

// signingSession holds all the logic allowing the signature of an ethereum transaction
type signingSession struct {
	account account.Account
	chain   *big.Int
}

// NewSigningSession create a new tx signature session from address
func NewSigningSession() SigningSession {
	return &signingSession{}
}

// SetAccount sets the account to the provided address
func (sess *signingSession) SetAccount(acc account.Account) error {
	sess.account = acc
	return nil
}

// SetChain is a setter for the chain used in the signed process
func (sess *signingSession) SetChain(netChain *big.Int) error {
	sess.chain = netChain
	return nil
}

// getSigner is internal function that returns an object used during the process
func (sess *signingSession) getSigner() (ethtypes.Signer, error) {
	var signer ethtypes.Signer
	if sess.chain == nil {
		return nil, errors.FailedPreconditionError("chain not set").ExtendComponent(component)
	}

	// We copy chain UUID to ensure pointer can be safely used elsewhere
	signer = ethtypes.NewEIP155Signer(new(big.Int).Set(sess.chain))

	return signer, nil
}

// ExecuteForTx : once all the element of the session have been set,
// it returns the signed transaction and the txhash
func (sess *signingSession) ExecuteForTx(tx *ethtypes.Transaction) ([]byte, *ethcommon.Hash, error) {
	signer, err := sess.getSigner()
	if err != nil {
		return []byte{}, nil, errors.CryptoOperationError(err.Error()).SetComponent(component)
	}

	t, err := ethtypes.SignTx(tx, signer, sess.account.Priv())
	if err != nil {
		return []byte{}, nil, errors.CryptoOperationError(err.Error()).SetComponent(component)
	}

	// Set raw transaction
	signedRaw, err := rlp.EncodeToBytes(t)
	if err != nil {
		return []byte{}, nil, errors.CryptoOperationError(err.Error()).SetComponent(component)
	}

	txHash := t.Hash()
	return signedRaw, &txHash, nil
}

// ExecuteForMsg : once all the element of the session have been set,
// it returns the message signature and the txhash
func (sess *signingSession) ExecuteForMsg(msg []byte, dsa crypto.DSA) ([]byte, *ethcommon.Hash, error) {
	rawHash := ethcrypto.Keccak256(msg)
	signature, err := dsa.Sign(rawHash, sess.account.Priv())
	if err != nil {
		return []byte{}, nil, errors.CryptoOperationError(err.Error()).SetComponent(component)
	}
	Hash := ethcommon.BytesToHash(rawHash)
	return signature, &Hash, nil
}

// ExecuteForEEATx : once all the element of the session have been set,
// it assigns the signed transaction and the txhash
// Calculates a transaction hash for sending private transactions for EEA extension
func (sess *signingSession) ExecuteForEEATx(tx *ethtypes.Transaction, privateArgs *types.PrivateArgs) ([]byte, *ethcommon.Hash, error) {
	// Should sign using the EIP155 signer
	signer, err := sess.getSigner()
	if err != nil {
		return []byte{}, nil, errors.CryptoOperationError(err.Error()).SetComponent(component)
	}

	hash, err := privateTxHash(tx, privateArgs, sess.chain)
	if err != nil {
		return nil, nil, err
	}

	sig, err := ethcrypto.Sign(hash[:], sess.account.Priv())
	if err != nil {
		return []byte{}, nil, errors.CryptoOperationError(err.Error()).SetComponent(component)
	}

	t, err := tx.WithSignature(signer, sig)
	if err != nil {
		return []byte{}, nil, errors.CryptoOperationError(err.Error()).SetComponent(component)
	}

	rplEncoding := encodePrivateTx(t, privateArgs)

	return rplEncoding, &hash, nil
}

// ExecuteForTesseraTx : once all the element of the session have been set,
// it assigns the signed transaction and the txhash
// Signs a transaction for Tessera private enclave
func (sess *signingSession) ExecuteForTesseraTx(tx *ethtypes.Transaction) ([]byte, *ethcommon.Hash, error) {
	// Transactions for Tessera should be signed using Homestead signer
	t, err := ethtypes.SignTx(tx, ethtypes.HomesteadSigner{}, sess.account.Priv())
	if err != nil {
		return []byte{}, nil, errors.CryptoOperationError(err.Error()).SetComponent(component)
	}

	// Get raw signed transaction
	signedRaw, err := rlp.EncodeToBytes(t)
	if err != nil {
		return []byte{}, nil, errors.CryptoOperationError(err.Error()).SetComponent(component)
	}

	// TODO: This does not match a hash returned by Tessera
	// Tessera transactions should be deterministic and we should be able to generate them locally

	v, r, s := t.RawSignatureValues()
	privateV := calculatePrivateV(v)

	txHash, err := rlpHash([]interface{}{
		tx.Nonce(),
		tx.GasPrice(),
		tx.Gas(),
		tx.To(),
		tx.Value(),
		tx.Data(),
		privateV,
		r,
		s,
	})
	if err != nil {
		return nil, nil, err
	}

	return signedRaw, &txHash, nil
}
