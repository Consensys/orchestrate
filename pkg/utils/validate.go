package utils

import (
	"math/big"
	"reflect"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-playground/validator/v10"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

var (
	validate      *validator.Validate
	StringPtrType = reflect.TypeOf(new(string))
	StringType    = reflect.TypeOf("")
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
	_, err := convDuration(fl)
	return err == nil
}

func minDuration(fl validator.FieldLevel) bool {
	min, err := time.ParseDuration(fl.Param())
	if err != nil {
		return false
	}

	v, err := convDuration(fl)
	if err != nil {
		return false
	}

	if v != 0 && v.Milliseconds() < min.Milliseconds() {
		return false
	}

	return true
}

func convDuration(fl validator.FieldLevel) (time.Duration, error) {
	switch fl.Field().Type() {
	case StringPtrType:
		val := fl.Field().Interface().(*string)
		if val != nil {
			return time.ParseDuration(*val)
		}
		return time.Duration(0), nil
	case StringType:
		if fl.Field().String() != "" {
			return time.ParseDuration(fl.Field().String())
		}
		return time.Duration(0), nil
	default:
		return time.Duration(0), nil
	}
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
			TesseraMarkingTransaction:
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
		case StatusCreated, StatusStarted, StatusPending, StatusRecovering, StatusWarning, StatusMined, StatusFailed, StatusStored, StatusResending:
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

func isKeyType(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		switch fl.Field().String() {
		case Secp256k1:
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
	_ = validate.RegisterValidation("minDuration", minDuration)
	_ = validate.RegisterValidation("isValidMethodSig", isValidMethodSig)
	_ = validate.RegisterValidation("isPrivateTxManagerType", isPrivateTxManagerType)
	_ = validate.RegisterValidation("isPriority", isPriority)
	_ = validate.RegisterValidation("isJobType", isJobType)
	_ = validate.RegisterValidation("isJobStatus", isJobStatus)
	_ = validate.RegisterValidation("isGasIncrementLevel", isGasIncrementLevel)
	_ = validate.RegisterValidation("isKeyType", isKeyType)
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
