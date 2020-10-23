package ethereum

import (
	"context"
	"encoding/base64"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/rlp"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"

	"github.com/ethereum/go-ethereum/crypto"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/store"
)

const signEEATransactionComponent = "use-cases.ethereum.sign-eea-transaction"

// signEEATxUseCase is a use case to sign a Quorum private transaction using an existing account
type signEEATxUseCase struct {
	vault store.Vault
}

// NewSignEEATransactionUseCase creates a new signEEATxUseCase
func NewSignEEATransactionUseCase(vault store.Vault) SignEEATransactionUseCase {
	return &signEEATxUseCase{
		vault: vault,
	}
}

// Execute signs a Quorum private transaction
func (uc *signEEATxUseCase) Execute(
	ctx context.Context,
	address, namespace string,
	chainID *big.Int,
	tx *ethtypes.Transaction,
	privateArgs *entities.PrivateETHTransactionParams,
) (string, error) {
	logger := log.WithContext(ctx).WithField("namespace", namespace).WithField("address", address)
	logger.Debug("signing eea private transaction")

	retrievedPrivKey, err := uc.vault.Ethereum().FindOne(ctx, address, namespace)
	if err != nil {
		return "", errors.FromError(err).ExtendComponent(signEEATransactionComponent)
	}

	privKey, err := NewECDSAFromPrivKey(retrievedPrivKey)
	if err != nil {
		return "", errors.FromError(err).ExtendComponent(signEEATransactionComponent)
	}

	h, err := computeHash(tx, privateArgs, chainID)
	if err != nil {
		return "", errors.FromError(err).ExtendComponent(signEEATransactionComponent)
	}

	signature, err := crypto.Sign(h, privKey)
	if err != nil {
		errMessage := "failed to sign eea private transaction"
		log.WithContext(ctx).WithError(err).Error(errMessage)
		return "", errors.CryptoOperationError(errMessage).ExtendComponent(signEEATransactionComponent)
	}

	logger.Info("eea private transaction signed successfully")
	return hexutil.Encode(signature), nil
}

func computeHash(tx *ethtypes.Transaction, privateArgs *entities.PrivateETHTransactionParams, chain *big.Int) ([]byte, error) {
	privateFromEncoded, err := base64.StdEncoding.DecodeString(privateArgs.PrivateFrom)
	if err != nil {
		errMessage := "invalid base64 privateFrom"
		log.WithError(err).WithField("private_from", privateArgs.PrivateFrom).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	var privateRecipientEncoded interface{}
	if privateArgs.PrivacyGroupID != "" {
		privateRecipientEncoded, err = base64.StdEncoding.DecodeString(privateArgs.PrivacyGroupID)
		if err != nil {
			errMessage := "invalid base64 privacyGroupId"
			log.WithError(err).WithField("privacy_group_id", privateArgs.PrivacyGroupID).Error(errMessage)
			return nil, errors.InvalidParameterError(errMessage)
		}
	} else {
		var privateForByteSlice [][]byte
		for _, v := range privateArgs.PrivateFor {
			b, der := base64.StdEncoding.DecodeString(v)
			if der != nil {
				errMessage := "invalid base64 privateFor"
				log.WithError(der).WithField("Private_for", v).Error(errMessage)
				return nil, errors.InvalidParameterError(errMessage)
			}
			privateForByteSlice = append(privateForByteSlice, b)
		}
		privateRecipientEncoded = privateForByteSlice
	}

	hash, err := rlp.Hash([]interface{}{
		tx.Nonce(),
		tx.GasPrice(),
		tx.Gas(),
		tx.To(),
		tx.Value(),
		tx.Data(),
		chain,
		uint(0),
		uint(0),
		privateFromEncoded,
		privateRecipientEncoded,
		privateArgs.PrivateTxType,
	})
	if err != nil {
		errMessage := "failed to hash eea transaction for signature"
		log.WithError(err).Error(errMessage)
		return nil, errors.CryptoOperationError(errMessage)
	}

	return hash[:], nil
}
