package utils

import (
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	gherkin "github.com/cucumber/messages-go/v10"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

// TODO: improve Regex to capture sub values instead of doing 2
var regexSlice = regexp.MustCompile(`\[([\d+])+]`)

func GetField(fieldPath string, val reflect.Value) (reflect.Value, error) {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if strings.Contains(fieldPath, ".") {
		keyValue := strings.Split(fieldPath, ".")

		var err error
		val, err = getField(keyValue[0], val)
		if err != nil {
			return reflect.Value{}, err
		}
		return GetField(strings.Join(keyValue[1:], "."), val)
	}

	return getField(fieldPath, val)
}

func getField(fieldPath string, val reflect.Value) (reflect.Value, error) {
	if val.Kind() == reflect.Invalid {
		return reflect.Value{}, errors.DataError("Could not get '%s' the field does not exist", fieldPath)
	}

	var key string
	if val.Kind() == reflect.Interface {
		val = val.Elem()
	}

	// if field is an array or slice
	sliceMatch := regexSlice.FindAllStringSubmatch(fieldPath, -1)
	if len(sliceMatch) > 0 {
		key = fieldPath
		for _, v := range sliceMatch {
			key = strings.Replace(key, v[0], "", 1)
		}
	} else {
		key = fieldPath
	}

	switch val.Kind() {
	case reflect.Map:
		val = val.MapIndex(reflect.ValueOf(key))
	default:
		val = val.FieldByName(key)
	}

	if val.Kind() == reflect.Interface {
		val = val.Elem()
	}

	// if field is an array or slice
	if len(sliceMatch) > 0 && val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		if len(sliceMatch) == 0 {
			return reflect.Value{}, fmt.Errorf("%s is an array - expected '[' ']'", fieldPath)
		}
		for _, v := range sliceMatch {
			i, _ := strconv.ParseInt(v[1], 10, 64)
			if int(i) >= val.Len() {
				return reflect.Value{}, fmt.Errorf("%s length is only %d could not reach %d", fieldPath, val.Len(), i)
			}
			val = val.Index(int(i))
		}
	}

	if val.Kind() == reflect.Invalid {
		return reflect.Value{}, errors.DataError("Could not get '%s' the field does not exist", fieldPath)
	}

	return val, nil
}

func RemoveCell(s []*gherkin.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell, index int) []*gherkin.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell {
	return append(s[:index], s[index+1:]...)
}

func ExtractTable(srcTable *gherkin.PickleStepArgument_PickleTable, fields []string) *gherkin.PickleStepArgument_PickleTable {
	isField := make(map[string]bool)
	for _, field := range fields {
		isField[field] = true
	}
	newTable := &gherkin.PickleStepArgument_PickleTable{
		Rows: []*gherkin.PickleStepArgument_PickleTable_PickleTableRow{{}},
	}

	// Find indexed of fields to be extracted in the headers and remove the fields in the original srcTable
	headers := srcTable.Rows[0]
	fieldIndex := make(map[int]bool)
	for i, h := range headers.Cells {
		if isField[h.Value] {
			fieldIndex[i] = true
			newTable.Rows[0].Cells = append(newTable.Rows[0].Cells, &gherkin.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
				Value: h.Value,
			})
			headers.Cells = RemoveCell(headers.Cells, i)
		}
	}

	if len(fieldIndex) == 0 {
		return nil
	}

	// Copy values in the new srcTable and remove the cell in the original srcTable
	for _, r := range srcTable.Rows[1:] {
		newRow := &gherkin.PickleStepArgument_PickleTable_PickleTableRow{}
		newTable.Rows = append(newTable.Rows, newRow)
		for i := range fieldIndex {
			newRow.Cells = append(newRow.Cells, &gherkin.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
				Value: r.Cells[i].Value,
			})
			r.Cells = RemoveCell(r.Cells, i)
		}
	}

	return newTable
}

func CmpField(field reflect.Value, value string) error {
	switch value {
	case "~":
		if isEmpty(field) {
			return fmt.Errorf("did not expected to be empty")
		}
	case "-":
		if !isEmpty(field) {
			return fmt.Errorf("expected to be empty but got %v", field)
		}
	default:
		if !isEqual(value, field) {
			return fmt.Errorf("expected '%value' but got '%v'", value, field)
		}
	}
	return nil
}

func isEqual(s string, val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(s, 10, 64)
		if val.Int() != n || err != nil {
			return false
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(s, 10, 64)
		if val.Uint() != n || err != nil {
			return false
		}
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(s, 64)
		if val.Float() != n || err != nil {
			return false
		}
	case reflect.String:
		if val.String() != s && fmt.Sprintf("%v", val) != s {
			return false
		}
	default:
		switch val.Type() {
		case reflect.TypeOf(&big.Int{}):
			if val.Interface().(*big.Int).String() != s {
				return false
			}
		case reflect.TypeOf(common.Address{}):
			if val.Interface().(common.Address).Hex() != s {
				return false
			}
		case reflect.TypeOf(common.Hash{}):
			if val.Interface().(common.Hash).Hex() != s {
				return false
			}
		}
	}
	return true
}

func isEmpty(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Ptr:
		return val.IsNil()
	case reflect.String:
		return val.String() == ""
	default:
		switch val.Type() {
		case reflect.TypeOf(common.Address{}):
			if val.Interface().(common.Address) == (common.Address{}) {
				return true
			}
		case reflect.TypeOf(common.Hash{}):
			if val.Interface().(common.Hash) != (common.Hash{}) {
				return true
			}
		}

	}
	return false
}
