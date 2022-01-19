package utils

import (
	"fmt"
	"math/big"
	"reflect"
	"strconv"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func StringerToString(v fmt.Stringer) string {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		// use of IsNil method
		if reflect.ValueOf(v).IsNil() {
			return ""
		}
	}

	return v.String()
}

func IsHexString(s string) bool {
	_, err := hexutil.Decode(s)
	return err == nil
}

func StringToHexBytes(v string) hexutil.Bytes {
	if v == "" {
		return nil
	}

	if vb, err := hexutil.Decode(v); err == nil {
		return vb
	}

	return nil
}

func BigIntStringToHex(v string) *hexutil.Big {
	if v == "" {
		return nil
	}

	if bv, ok := new(big.Int).SetString(v, 10); ok {
		return (*hexutil.Big)(bv)
	}

	return nil
}

func HexToBigIntString(v *hexutil.Big) string {
	if v == nil {
		return ""
	}

	return v.ToInt().String()
}

func StringToUint64(v string) *uint64 {
	if v == "" {
		return nil
	}

	if vi, err := strconv.ParseUint(v, 10, 64); err == nil {
		return &vi
	}

	return nil
}
