package utils

import (
	"fmt"
	"reflect"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

// ShortString makes hashes short for a limited column size
func ShortString(s string, tailLength int) string {
	runes := []rune(s)
	if len(runes)/2 > tailLength {
		return string(runes[:tailLength]) + "..." + string(runes[len(runes)-tailLength:])
	}
	return s
}

func ValueToString(v interface{}) string {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		// use of IsNil method
		if reflect.ValueOf(v).IsNil() {
			return ""
		}

		return fmt.Sprintf("%v", reflect.ValueOf(v).Elem())
	}

	return fmt.Sprintf("%v", v)
}

func StringToEthHash(s string) *ethcommon.Hash {
	if s == "" {
		return nil
	}

	hash := ethcommon.HexToHash(s)
	return &hash
}

func ToEthAddr(s string) *ethcommon.Address {
	if s == "" {
		return nil
	}

	add := ethcommon.HexToAddress(s)
	return &add
}
