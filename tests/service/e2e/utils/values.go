package utils

import (
	"encoding/json"
	"math/big"
	"reflect"
	"strings"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/types/tx"
	gherkin "github.com/cucumber/messages-go/v10"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mitchellh/mapstructure"
)

func underlyingType(structType reflect.Type) reflect.Type {
	if structType.Kind() != reflect.Ptr {
		return structType
	}
	return underlyingType(structType.Elem())
}

// TODO: unmarshall any slice types
func toValue(structType reflect.Type, strValue string) (reflect.Value, error) {
	var v reflect.Value
	var err error
	switch structType {
	case reflect.TypeOf("string"), reflect.TypeOf(entities.PrivateTxManagerType("")):
		v = reflect.ValueOf(strValue)
	case reflect.TypeOf(common.Address{}):
		v = reflect.ValueOf(common.HexToAddress(strValue))
	case reflect.TypeOf(common.Hash{}):
		v = reflect.ValueOf(common.HexToHash(strValue))
	case reflect.TypeOf(big.Int{}):
		b, ok := new(big.Int).SetString(strValue, 10)
		if !ok {
			err = errors.DataError("%s is an invalid %s", strValue, structType.Name())
		}
		v = reflect.ValueOf(b)
	// Should cover ints, uints, []ints, []string....
	default:
		inter := reflect.New(structType).Interface()
		err = json.Unmarshal([]byte(strValue), inter)
		v = reflect.ValueOf(inter).Elem()
	}
	if err != nil {
		return reflect.Value{}, errors.DataError("invalid %s - got %v", structType.Name(), err)
	}
	return v, nil
}

func toMapStruct(structType reflect.Type, mapStruct map[string]interface{}, keys, values []string) error {
	structType = underlyingType(structType)

	for i, key := range keys {

		k := strings.Split(key, ".")

		// assign value if not a nested struct
		if len(k) == 1 {
			var f reflect.Type
			if structType.Kind() == reflect.Map {
				f = structType.Elem()
			} else {
				fi, ok := structType.FieldByName(strings.Title(key))
				if !ok {
					return errors.DataError("%s field does not exist in %s", key, structType.Name())
				}
				f = fi.Type
			}
			typ := underlyingType(f)
			if values[i] == "" {
				continue
			}
			v, err := toValue(typ, values[i])
			if err != nil {
				return err
			}
			reflect.ValueOf(mapStruct).SetMapIndex(reflect.ValueOf(key), v)
			continue
		}

		// check if the nested field exist
		field, ok := structType.FieldByName(strings.Title(k[0]))
		if !ok {
			return errors.DataError("%s field does not exist in %s", k[0], structType.Name())
		}

		// check if the nested field has already been assigned
		nestedMap, ok := mapStruct[k[0]]
		if !ok {
			m := make(map[string]interface{})
			mapStruct[k[0]] = m
			err := toMapStruct(field.Type, m, []string{strings.Join(k[1:], ".")}, []string{values[i]})
			if err != nil {
				return errors.FromError(err)
			}
			continue
		}

		// check if the nested field is a map to interface
		m, ok := nestedMap.(map[string]interface{})
		if !ok {
			return errors.DataError("invalid payload")
		}
		err := toMapStruct(field.Type, m, []string{strings.Join(k[1:], ".")}, []string{values[i]})
		if err != nil {
			return errors.FromError(err)
		}
	}

	return nil
}

func ParseTable(i interface{}, table *gherkin.PickleStepArgument_PickleTable) ([]interface{}, error) {
	var interfaces []interface{}

	header := table.Rows[0]
	rows := table.Rows[1:]

	var keys []string
	for _, h := range header.Cells {
		keys = append(keys, h.Value)
	}

	for _, row := range rows {
		var values []string
		for _, r := range row.Cells {
			values = append(values, r.Value)
		}
		mapStruct := make(map[string]interface{})
		err := toMapStruct(reflect.TypeOf(i), mapStruct, keys, values)

		if err != nil {
			return nil, errors.FromError(err)
		}

		newStruct := reflect.New(reflect.TypeOf(i)).Interface()
		err = mapstructure.Decode(mapStruct, newStruct)
		if err != nil {
			return nil, errors.FromError(err)
		}
		interfaces = append(interfaces, newStruct)
	}

	return interfaces, nil
}

func ParseEnvelope(table *gherkin.PickleStepArgument_PickleTable) ([]*tx.Envelope, error) {
	var envelopes []*tx.Envelope

	// TODO: Parse Enum for METHOD or JOBTYPE

	interfaces, err := ParseTable(tx.Envelope{}, table)
	if err != nil {
		return nil, err
	}

	for _, e := range interfaces {
		envelopes = append(envelopes, e.(*tx.Envelope).SafeEnvelope())
	}

	return envelopes, nil
}
