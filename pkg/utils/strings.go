package utils

import (
	"encoding/json"
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

func BytesToString(b []byte) string {
	if b == nil {
		return ""
	}

	return string(b)
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

func ParseIArrayToStringArray(ints []interface{}) ([]string, error) {
	strings := make([]string, len(ints))
	for idx, val := range ints {
		switch reflect.TypeOf(val).Kind() {
		case reflect.Slice:
			rVal := reflect.ValueOf(val)
			ret := make([]interface{}, rVal.Len())
			for jdx := 0; jdx < rVal.Len(); jdx++ {
				ret[jdx] = rVal.Index(jdx).Interface()
			}

			sv, err := ParseIArrayToStringArray(ret)
			if err != nil {
				return []string{}, err
			}

			b, err := json.Marshal(sv)
			if err != nil {
				return []string{}, err
			}
			strings[idx] = string(b)
		default:
			strings[idx] = fmt.Sprint(val)
		}
	}

	return strings, nil
}
