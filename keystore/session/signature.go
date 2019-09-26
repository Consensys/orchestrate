package session

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
	pkgchain "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/multi-vault.git/keystore/crypto"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/multi-vault.git/keystore/wallet"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/multi-vault.git/secretstore/services"
	"golang.org/x/crypto/sha3"
)

// SigningSession holds all the logic allowing the signature of an ethereum transaction
type SigningSession struct {
	secretStore services.SecretStore
	wallet      *wallet.Wallet
	chain       *pkgchain.Chain
}

// NewSigningSession create a new tx signature session from address
func NewSigningSession(secretStore services.SecretStore) *SigningSession {
	return &SigningSession{
		secretStore: secretStore,
	}
}

// SetWallet sets the wallet to the provided address
func (sess *SigningSession) SetWallet(address *ethcommon.Address) error {
	w := wallet.NewWallet(sess.secretStore)
	err := w.Load(address)
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	sess.wallet = w
	return nil
}

// SetChain is a setter for the chain used in the signed process
func (sess *SigningSession) SetChain(netChain *pkgchain.Chain) error {
	sess.chain = netChain
	return nil
}

// getSigner is internal function that returns an object used during the process
func (sess *SigningSession) getSigner() (ethtypes.Signer, error) {
	var signer ethtypes.Signer
	if sess.chain == nil {
		return nil, errors.FailedPreconditionError("chain not set").ExtendComponent(component)
	}

	// We copy chain ID to ensure pointer can be safely used elsewhere
	signer = ethtypes.NewEIP155Signer(new(big.Int).Set(sess.chain.ID()))

	return signer, nil
}

// ExecuteForTx : once all the element of the session have been set,
// it returns the signed transaction and the txhash
func (sess *SigningSession) ExecuteForTx(tx *ethtypes.Transaction) ([]byte, *ethcommon.Hash, error) {
	signer, err := sess.getSigner()
	if err != nil {
		return []byte{}, nil, errors.CryptoOperationError(err.Error()).SetComponent(component)
	}

	t, err := ethtypes.SignTx(tx, signer, sess.wallet.Priv())
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
func (sess *SigningSession) ExecuteForMsg(msg []byte, dsa crypto.DSA) ([]byte, *ethcommon.Hash, error) {
	rawHash := ethcrypto.Keccak256(msg)
	signature, err := dsa.Sign(rawHash, sess.wallet.Priv())
	if err != nil {
		return []byte{}, nil, errors.CryptoOperationError(err.Error()).SetComponent(component)
	}
	Hash := ethcommon.BytesToHash(rawHash)
	return signature, &Hash, nil
}

// ExecuteForEEATx : once all the element of the session have been set,
// it assigns the signed transaction and the txhash
// Calculates a transaction hash for sending private transactions for EEA extension
func (sess *SigningSession) ExecuteForEEATx(tx *ethtypes.Transaction, privateArgs *types.PrivateArgs) ([]byte, *ethcommon.Hash, error) {
	// Should sign using the EIP155 signer
	signer, err := sess.getSigner()
	if err != nil {
		return []byte{}, nil, errors.CryptoOperationError(err.Error()).SetComponent(component)
	}

	hash := privateTxHash(tx, privateArgs, sess.chain)

	sig, err := ethcrypto.Sign(hash[:], sess.wallet.Priv())
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
func (sess *SigningSession) ExecuteForTesseraTx(tx *ethtypes.Transaction) ([]byte, *ethcommon.Hash, error) {
	// Transactions for Tessera should be signed using Homestead signer
	t, err := ethtypes.SignTx(tx, ethtypes.HomesteadSigner{}, sess.wallet.Priv())
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

	txHash := rlpHash([]interface{}{
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

	return signedRaw, &txHash, nil
}

func privateTxHash(tx *ethtypes.Transaction, privateArgs *types.PrivateArgs, chain *pkgchain.Chain) ethcommon.Hash {
	hash := rlpHash([]interface{}{
		tx.Nonce(),
		tx.GasPrice(),
		tx.Gas(),
		tx.To(),
		tx.Value(),
		tx.Data(),
		chain.ID(),
		uint(0),
		uint(0),
		privateArgs.PrivateFrom,
		privateArgs.PrivateFor,
		privateArgs.PrivateTxType,
	})
	return hash
}

func encodePrivateTx(tx *ethtypes.Transaction, privateArgs *types.PrivateArgs) []byte {
	v, r, s := tx.RawSignatureValues()
	rplEncoding, _ := rlpEncode([]interface{}{
		tx.Nonce(),
		tx.GasPrice(),
		tx.Gas(),
		tx.To(),
		tx.Value(),
		tx.Data(),
		v,
		r,
		s,
		privateArgs.PrivateFrom,
		privateArgs.PrivateFor,
		privateArgs.PrivateTxType,
	})
	return rplEncoding
}

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
