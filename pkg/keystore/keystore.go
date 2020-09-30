package keystore

import (
	"context"
	"math/big"

	quorumtypes "github.com/consensys/quorum/core/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore/crypto/signature"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore/session"
)

// KeyStore holds the methods of the interfaces BaseKeyStore
type keystore struct {
	sessMng session.AccountManager
}

// NewKeyStore construct a BaseKeyStore from a client
func NewKeyStore(mng session.AccountManager) KeyStore {
	return &keystore{
		sessMng: mng,
	}
}

// SignTx returns a signed transaction. It is perfectly equivalent to SignTx
func (ks *keystore) SignTx(ctx context.Context, netChain *big.Int, addr ethcommon.Address, tx *ethtypes.Transaction) ([]byte, *ethcommon.Hash, error) {
	// Creates a new signing session
	sess, err := ks.sessMng.SigningSession(ctx, addr)
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
func (ks *keystore) SignPrivateEEATx(ctx context.Context, netChain *big.Int, addr ethcommon.Address, tx *ethtypes.Transaction, privateArgs *types.PrivateArgs) ([]byte, *ethcommon.Hash, error) {
	// Creates a new signing session
	sess, err := ks.sessMng.SigningSession(ctx, addr)
	if err != nil {
		return []byte{}, nil, errors.FromError(err).ExtendComponent(component)
	}

	err = sess.SetChain(netChain)
	if err != nil {
		return []byte{}, nil, errors.FromError(err).ExtendComponent(component)
	}

	// Run signing session
	PrivRaw, EnclaveKey, err := sess.ExecuteForEEATx(tx, privateArgs)
	if err != nil {
		return []byte{}, nil, err
	}
	return PrivRaw, EnclaveKey, nil
}

// SignPrivateTesseraTx signs a private transaction using Tessera
func (ks *keystore) SignPrivateTesseraTx(ctx context.Context, netChain *big.Int, addr ethcommon.Address, tx *quorumtypes.Transaction) ([]byte, *ethcommon.Hash, error) {
	// Creates a new signing session
	sess, err := ks.sessMng.SigningSession(ctx, addr)
	if err != nil {
		return []byte{}, nil, err
	}

	// Set tx.data.V = 37 or 38 for signature
	if !tx.IsPrivate() {
		tx.SetPrivate()
	}

	Raw, Hash, err := sess.ExecuteForTesseraTx(tx)
	if err != nil {
		return []byte{}, nil, err
	}
	return Raw, Hash, nil
}

// SignMsg returns a signed message and its hash
func (ks *keystore) SignMsg(ctx context.Context, addr ethcommon.Address, msg string) ([]byte, *ethcommon.Hash, error) {
	// Creates a new signing session
	sess, err := ks.sessMng.SigningSession(ctx, addr)
	if err != nil {
		return []byte{}, nil, errors.FromError(err).ExtendComponent(component)
	}

	// Run signing session
	// TODO: Add a possibility to use another ecdsa than ethereum'ks
	//
	Raw, Hash, err := sess.ExecuteForMsg([]byte(msg), signature.EthECDSA)
	if err != nil {
		return []byte{}, nil, errors.FromError(err).ExtendComponent(component)
	}

	return Raw, Hash, nil
}

// SignRawHash returns a signed raw hash
func (ks *keystore) SignRawHash(
	a ethcommon.Address,
	hash []byte,
) (rsv []byte, err error) {
	return []byte{}, errors.FeatureNotSupportedError("SignRawHash not implemented yet").SetComponent(component)
}

// GenerateAccount create and stores a new account in the vault
func (ks *keystore) GenerateAccount(ctx context.Context) (ethcommon.Address, error) {
	return ks.sessMng.GenerateAccount(ctx)
}

// ImportPrivateKey adds a private key in the vault
func (ks *keystore) ImportPrivateKey(ctx context.Context, priv string) error {
	return ks.sessMng.ImportPrivateKey(ctx, priv)
}
