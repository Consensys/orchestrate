package base

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/types"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/multi-vault.git/keystore/session"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/multi-vault.git/keystore/wallet"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/multi-vault.git/secretstore/services"
)

// KeyStore olds the methods of the interfaces BaseKeyStore
type KeyStore struct {
	SecretStore services.SecretStore
}

// NewKeyStore construct a BaseKeyStore from a client
func NewKeyStore(secretStore services.SecretStore) *KeyStore {
	return &KeyStore{
		SecretStore: secretStore,
	}
}

// SignTx returns a signed transaction. It is perfectly equivalent to SignTx
func (s *KeyStore) SignTx(netChain *chain.Chain, a ethcommon.Address, tx *ethtypes.Transaction) (raw []byte, hash *ethcommon.Hash, err error) {
	// Creates a new signing session
	sess := session.MakeTxSignature(s.SecretStore)
	err = sess.SetWallet(&a)
	if err != nil {
		return []byte{}, nil, errors.FromError(err).ExtendComponent(component)
	}
	sess.SetChain(netChain)
	sess.SetTx(tx)

	// Run signing session
	err = sess.Run()
	if err != nil {
		return []byte{}, nil, errors.FromError(err).ExtendComponent(component)
	}

	return sess.Raw, sess.Hash, nil
}

// SignPrivateEEATx signs a private transaction
func (s *KeyStore) SignPrivateEEATx(netChain *chain.Chain, a ethcommon.Address, tx *ethtypes.Transaction, privateArgs *types.PrivateArgs) (raw []byte, txHash *ethcommon.Hash, err error) {
	// Creates a new signing session
	sess := session.MakeTxSignature(s.SecretStore)
	err = sess.SetWallet(&a)
	if err != nil {
		return []byte{}, nil, errors.FromError(err).ExtendComponent(component)
	}
	sess.SetChain(netChain)
	sess.SetTx(tx)
	sess.SetPrivateArgs(privateArgs)

	// Run signing session
	err = sess.SignPrivateEEATransaction()
	if err != nil {
		return []byte{}, nil, err
	}

	return sess.Raw, sess.Hash, nil
}

// SignPrivateTesseraTx signs a private transaction using Tessera
func (s *KeyStore) SignPrivateTesseraTx(netChain *chain.Chain, a ethcommon.Address, tx *ethtypes.Transaction) (raw []byte, txHash *ethcommon.Hash, err error) {
	// Creates a new signing session
	sess := session.MakeTxSignature(s.SecretStore)
	err = sess.SetWallet(&a)
	if err != nil {
		return []byte{}, nil, err
	}
	sess.SetTx(tx)

	// Run signing session
	err = sess.SignPrivateTesseraTransaction()
	if err != nil {
		return []byte{}, nil, errors.FromError(err).ExtendComponent(component)
	}

	return sess.Raw, sess.Hash, nil
}

// SignMsg returns a signed message and its hash
func (s *KeyStore) SignMsg(a ethcommon.Address, msg string) (rsv []byte, hash *ethcommon.Hash, err error) {
	return []byte{}, nil, errors.FeatureNotSupportedError("SignMsg not implemented yet").SetComponent(component)
}

// SignRawHash returns a signed raw hash
func (s *KeyStore) SignRawHash(
	a ethcommon.Address,
	hash []byte,
) (rsv []byte, err error) {

	return []byte{}, errors.FeatureNotSupportedError("SignRawHash not implemented yet").SetComponent(component)
}

// GenerateWallet create and stores a new wallet in the vault
func (s *KeyStore) GenerateWallet() (add *ethcommon.Address, err error) {
	w := wallet.NewWallet(s.SecretStore)
	err = w.Generate()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	err = w.Store()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return w.Address(), nil
}

// ImportPrivateKey adds a private key in the vault
// TODO: this is Unsafe and should be removed soon
func (s *KeyStore) ImportPrivateKey(priv string) (err error) {

	w := wallet.NewWallet(s.SecretStore)
	err = w.FromPrivateKey(priv)
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	err = w.Store()
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	return nil
}
