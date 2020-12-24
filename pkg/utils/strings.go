package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// ShortString makes hashes short for a limited column size
func ShortString(s string, tailLength int) string {
	runes := []rune(s)
	if len(runes)/2 > tailLength {
		return string(runes[:tailLength]) + "..." + string(runes[len(runes)-tailLength:])
	}
	return s
}

func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[seededRand.Intn(len(letter))]
	}
	return string(b)
}

func RandHexString(n int) string {
	var letterRunes = []rune("abcdef0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[seededRand.Intn(len(letterRunes))]
	}
	return string(b)
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
