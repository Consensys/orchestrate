package validators

import (
	"context"
	"crypto/md5"
	"encoding/hex"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/utils"
)

//go:generate mockgen -source=transactions.go -destination=mocks/transactions.go -package=mocks

const txValidatorComponent = "transaction-validator"

type TransactionValidator interface {
	ValidateRequestHash(ctx context.Context, params interface{}, idempotencyKey string) (string, error)
}

// Transaction is a validator for transaction requests (business logic)
type transactionValidator struct {
	txRequestDA store.TransactionRequestAgent
}

// NewTransactionValidator creates a new TransactionValidator
func NewTransactionValidator(txRequestDA store.TransactionRequestAgent) TransactionValidator {
	return &transactionValidator{txRequestDA: txRequestDA}
}

func (txValidator *transactionValidator) ValidateRequestHash(ctx context.Context, params interface{}, idempotencyKey string) (string, error) {
	jsonParams, err := utils.ObjectToJSON(params)
	if err != nil {
		return "", errors.FromError(err).ExtendComponent(txValidatorComponent)
	}

	hash := md5.Sum([]byte(jsonParams))
	requestHash := hex.EncodeToString(hash[:])

	txRequestToCompare, err := txValidator.txRequestDA.FindOneByIdempotencyKey(ctx, idempotencyKey)
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
