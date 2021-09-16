package utils

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// ShortString makes hashes short for a limited column size
func ShortString(s string, tailLength int) string {
	runes := []rune(s)
	if len(runes)/2 > tailLength {
		return string(runes[:tailLength]) + "..." + string(runes[len(runes)-tailLength:])
	}
	return s
}

func IsHexString(s string) bool {
	_, err := hexutil.Decode(s)
	return err == nil
}

func MustEncodeBigInt(s string) *big.Int {
	b, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic("failed to convert to big.Int")
	}

	return b
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
