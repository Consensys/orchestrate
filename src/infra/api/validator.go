package api

import (
	"math/big"
	"reflect"
	"time"

	"github.com/consensys/orchestrate/pkg/utils"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-playground/validator/v10"
)

var (
	validate      *validator.Validate
	StringPtrType = reflect.TypeOf(new(string))
	StringType    = reflect.TypeOf("")
)

func isHex(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		return utils.IsHexString(fl.Field().String())
	}

	return true
}

func isHexAddress(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		return ethcommon.IsHexAddress(fl.Field().String())
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

func init() {
	if validate != nil {
		return
	}

	validate = validator.New()
	_ = validate.RegisterValidation("isHex", isHex)
	_ = validate.RegisterValidation("isHexAddress", isHexAddress)
	_ = validate.RegisterValidation("isBig", isBig)
	_ = validate.RegisterValidation("isHash", isHash)
	_ = validate.RegisterValidation("isDuration", isDuration)
	_ = validate.RegisterValidation("minDuration", minDuration)
}

func GetValidator() *validator.Validate {
	return validate
}
