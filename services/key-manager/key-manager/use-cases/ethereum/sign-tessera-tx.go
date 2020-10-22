package ethereum

import (
	"context"

	"github.com/consensys/quorum/common/hexutil"
	quorumtypes "github.com/consensys/quorum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/store"
)

const signTesseraTransactionComponent = "use-cases.ethereum.sign-tessera-transaction"

// signTesseraTxUseCase is a use case to sign a tessera transaction using an existing account
type signTesseraTxUseCase struct {
	vault store.Vault
}

// NewSignTransactionUseCase creates a new signTesseraTxUseCase
func NewSignTesseraTransactionUseCase(vault store.Vault) SignTesseraTransactionUseCase {
	return &signTesseraTxUseCase{
		vault: vault,
	}
}

// Execute signs a Tessera private transaction
func (uc *signTesseraTxUseCase) Execute(ctx context.Context, address, namespace string, tx *quorumtypes.Transaction) (string, error) {
	logger := log.WithContext(ctx).WithField("namespace", namespace).WithField("address", address)
	logger.Debug("signing tessera private transaction")

	retrievedPrivKey, err := uc.vault.Ethereum().FindOne(ctx, address, namespace)
	if err != nil {
		return "", errors.FromError(err).ExtendComponent(signTesseraTransactionComponent)
	}

	privKey, err := NewECDSAFromPrivKey(retrievedPrivKey)
	if err != nil {
		return "", errors.FromError(err).ExtendComponent(signTesseraTransactionComponent)
	}

	signer := quorumtypes.QuorumPrivateTxSigner{}
	tx.SetPrivate()
	h := signer.Hash(tx)
	signature, err := crypto.Sign(h[:], privKey)
	if err != nil {
		errMessage := "failed to sign tessera private transaction"
		log.WithContext(ctx).WithError(err).Error(errMessage)
		return "", errors.CryptoOperationError(errMessage).ExtendComponent(signTesseraTransactionComponent)
	}

	logger.Info("tessera private transaction signed successfully")
	return hexutil.Encode(signature), nil
}
