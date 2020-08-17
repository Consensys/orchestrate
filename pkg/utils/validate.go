package utils

import (
	"math/big"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-playground/validator/v10"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

var (
	validate *validator.Validate
)

func isHex(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		_, err := hexutil.Decode(fl.Field().String())
		if err != nil {
			return false
		}
	}

	return true
}

func isBig(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		_, ok := new(big.Int).SetString(fl.Field().String(), 10)
		return ok
	}

	return true
}

func isHash(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		return IsHash(fl.Field().String())
	}

	return true
}

func IsHash(input string) bool {
	hash, err := hexutil.Decode(input)
	if err != nil || len(hash) != ethcommon.HashLength {
		return false
	}

	return true
}

func isValidMethodSig(fl validator.FieldLevel) bool {
	return IsValidSignature(fl.Field().String())
}

func isDuration(fl validator.FieldLevel) bool {
	if fl.Field().Type().String() == "*string" {
		val := fl.Field().Interface().(*string)
		if val != nil {
			_, err := time.ParseDuration(*val)
			if err != nil {
				return false
			}
		}
	} else if fl.Field().String() != "" {
		_, err := time.ParseDuration(fl.Field().String())
		if err != nil {
			return false
		}
	}

	return true
}

func isPrivateTxManagerType(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		switch fl.Field().String() {
		case TesseraChainType, OrionChainType:
			return true
		default:
			return false
		}
	}

	return true
}

func isPriority(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		switch fl.Field().String() {
		case PriorityVeryLow, PriorityLow, PriorityMedium, PriorityHigh, PriorityVeryHigh:
			return true
		default:
			return false
		}
	}

	return true
}

func isJobType(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		switch fl.Field().String() {
		case
			EthereumTransaction,
			EthereumRawTransaction,
			OrionEEATransaction,
			OrionMarkingTransaction,
			TesseraPrivateTransaction,
			TesseraPublicTransaction:
			return true
		default:
			return false
		}
	}

	return true
}

func isJobStatus(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		switch fl.Field().String() {
		case StatusCreated, StatusStarted, StatusPending, StatusRecovering, StatusWarning, StatusMined, StatusFailed:
			return true
		default:
			return false
		}
	}

	return true
}

func isGasIncrementLevel(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		switch fl.Field().String() {
		case GasIncrementVeryLow, GasIncrementLow, GasIncrementMedium, GasIncrementHigh, GasIncrementVeryHigh:
			return true
		default:
			return false
		}
	}

	return true
}

func init() {
	if validate != nil {
		return
	}

	validate = validator.New()
	_ = validate.RegisterValidation("isHex", isHex)
	_ = validate.RegisterValidation("isBig", isBig)
	_ = validate.RegisterValidation("isHash", isHash)
	_ = validate.RegisterValidation("isDuration", isDuration)
	_ = validate.RegisterValidation("isValidMethodSig", isValidMethodSig)
	_ = validate.RegisterValidation("isPrivateTxManagerType", isPrivateTxManagerType)
	_ = validate.RegisterValidation("isPriority", isPriority)
	_ = validate.RegisterValidation("isJobType", isJobType)
	_ = validate.RegisterValidation("isJobStatus", isJobStatus)
	_ = validate.RegisterValidation("isGasIncrementLevel", isGasIncrementLevel)
}

func GetValidator() *validator.Validate {
	return validate
}

func HandleValidatorError(validatorErrors validator.ValidationErrors) []error {
	errs := make([]error, 0)
	for _, validatorError := range validatorErrors {
		err := errors.DataError("invalid %s got %s", validatorError.Field(), validatorError.Value())
		errs = append(errs, err)
	}

	return errs
}
