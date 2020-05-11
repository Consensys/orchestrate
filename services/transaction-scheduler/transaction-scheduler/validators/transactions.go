package validators

import (
	"context"
	"crypto/md5"
	"encoding/hex"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/interfaces"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/utils"
)

//go:generate mockgen -source=transactions.go -destination=mocks/transactions.go -package=mocks

const txValidatorComponent = "transaction-validator"

type TransactionValidator interface {
	ValidateRequestHash(ctx context.Context, params interface{}, idempotencyKey string) (string, error)
	ValidateChainExists(ctx context.Context, chainUUID string) error
}

// transactionValidator is a validator for transaction requests (business logic)
type transactionValidator struct {
	db                  interfaces.DB
	chainRegistryClient client.ChainRegistryClient
}

// NewTransactionValidator creates a new TransactionValidator
func NewTransactionValidator(db interfaces.DB, chainRegistryClient client.ChainRegistryClient) TransactionValidator {
	return &transactionValidator{db: db, chainRegistryClient: chainRegistryClient}
}

func (txValidator *transactionValidator) ValidateRequestHash(ctx context.Context, params interface{}, idempotencyKey string) (string, error) {
	log.WithContext(ctx).WithField("idempotency_key", idempotencyKey).Debug("validating idempotency key")

	jsonParams, err := utils.ObjectToJSON(params)
	if err != nil {
		return "", errors.FromError(err).ExtendComponent(txValidatorComponent)
	}

	hash := md5.Sum([]byte(jsonParams))
	requestHash := hex.EncodeToString(hash[:])

	txRequestToCompare, err := txValidator.db.TransactionRequest().FindOneByIdempotencyKey(ctx, idempotencyKey)
	if err != nil && !errors.IsNotFoundError(err) {
		return "", errors.FromError(err).ExtendComponent(txValidatorComponent)
	}

	if txRequestToCompare != nil && txRequestToCompare.RequestHash != requestHash {
		errMessage := "a transaction request with the same idempotency key and different params already exists"
		log.WithError(err).WithField("idempotency_key", idempotencyKey).Error(errMessage)
		return "", errors.AlreadyExistsError(errMessage)
	}

	return requestHash, nil
}

func (txValidator *transactionValidator) ValidateChainExists(ctx context.Context, chainUUID string) error {
	// Validate that the chainUUID exists
	_, err := txValidator.chainRegistryClient.GetChainByUUID(ctx, chainUUID)
	if err != nil {
		errMessage := "failed to get chain"
		log.WithError(err).WithField("chain_uuid", chainUUID).Error(errMessage)
		return errors.InvalidParameterError(errMessage)
	}

	return nil
}
