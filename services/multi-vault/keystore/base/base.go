package base

import (
	"context"
	"math/big"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/types"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/keystore/account"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/keystore/crypto/signature"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/keystore/session"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/secretstore/services"
)

// KeyStore holds the methods of the interfaces BaseKeyStore
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
func (s *KeyStore) SignTx(ctx context.Context, netChain *big.Int, a ethcommon.Address, tx *ethtypes.Transaction) ([]byte, *ethcommon.Hash, error) {
	// Creates a new signing session
	sess := session.NewSigningSession(s.SecretStore)
	err := sess.SetAccount(ctx, &a)
	if err != nil {
		return []byte{}, nil, errors.FromError(err).ExtendComponent(component)
	}

	err = sess.SetChain(netChain)
	if err != nil {
		return []byte{}, nil, errors.FromError(err).ExtendComponent(component)
	}

	// Run signing session
	Raw, Hash, err := sess.ExecuteForTx(tx)
	if err != nil {
		return []byte{}, nil, errors.FromError(err).ExtendComponent(component)
	}
	return Raw, Hash, nil
}

// SignPrivateEEATx signs a private transaction
func (s *KeyStore) SignPrivateEEATx(ctx context.Context, netChain *big.Int, a ethcommon.Address, tx *ethtypes.Transaction, privateArgs *types.PrivateArgs) ([]byte, *ethcommon.Hash, error) {
	// Creates a new signing session
	sess := session.NewSigningSession(s.SecretStore)
	err := sess.SetAccount(ctx, &a)
	if err != nil {
		return []byte{}, nil, errors.FromError(err).ExtendComponent(component)
	}
	err = sess.SetChain(netChain)
	if err != nil {
		return []byte{}, nil, errors.FromError(err).ExtendComponent(component)
	}

	// Run signing session
	Raw, Hash, err := sess.ExecuteForEEATx(tx, privateArgs)
	if err != nil {
		return []byte{}, nil, err
	}
	return Raw, Hash, nil
}

// SignPrivateTesseraTx signs a private transaction using Tessera
func (s *KeyStore) SignPrivateTesseraTx(ctx context.Context, netChain *big.Int, a ethcommon.Address, tx *ethtypes.Transaction) ([]byte, *ethcommon.Hash, error) {
	// Creates a new signing session
	sess := session.NewSigningSession(s.SecretStore)
	err := sess.SetAccount(ctx, &a)
	if err != nil {
		return []byte{}, nil, err
	}

	Raw, Hash, err := sess.ExecuteForTesseraTx(tx)
	if err != nil {
		return []byte{}, nil, err
	}
	return Raw, Hash, nil
}

// SignMsg returns a signed message and its hash
func (s *KeyStore) SignMsg(ctx context.Context, a ethcommon.Address, msg string) ([]byte, *ethcommon.Hash, error) {
	// Creates a new signing session
	sess := session.NewSigningSession(s.SecretStore)
	err := sess.SetAccount(ctx, &a)
	if err != nil {
		return []byte{}, nil, errors.FromError(err).ExtendComponent(component)
	}

	// Run signing session
	// TODO: Add a possibility to use another ecdsa than ethereum's
	//
	Raw, Hash, err := sess.ExecuteForMsg([]byte(msg), signature.EthECDSA)
	if err != nil {
		return []byte{}, nil, errors.FromError(err).ExtendComponent(component)
	}

	return Raw, Hash, nil
}

// SignRawHash returns a signed raw hash
func (s *KeyStore) SignRawHash(
	a ethcommon.Address,
	hash []byte,
) (rsv []byte, err error) {

	return []byte{}, errors.FeatureNotSupportedError("SignRawHash not implemented yet").SetComponent(component)
}

// GenerateAccount create and stores a new account in the vault
func (s *KeyStore) GenerateAccount(ctx context.Context) (*ethcommon.Address, error) {
	w := account.NewAccount(s.SecretStore)
	err := w.Generate()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	err = w.Store(ctx)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return w.Address(), nil
}

// ImportPrivateKey adds a private key in the vault
// TODO: this is Unsafe and should be removed soon
func (s *KeyStore) ImportPrivateKey(ctx context.Context, priv string) error {

	w := account.NewAccount(s.SecretStore)
	err := w.FromPrivateKey(priv)
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	err = w.Store(ctx)
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	return nil
}
