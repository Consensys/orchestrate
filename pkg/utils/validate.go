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
		hash, err := hexutil.Decode(fl.Field().String())
		if err != nil || len(hash) != ethcommon.HashLength {
			return false
		}
	}

	return true
}

func isValidMethodSig(fl validator.FieldLevel) bool {
	return IsValidSignature(fl.Field().String())
}

func isDuration(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		_, err := time.ParseDuration(fl.Field().String())
		if err != nil {
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
