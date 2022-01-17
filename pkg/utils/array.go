package utils

import (
	"encoding/json"
)

func CastInterfaceToObject(data, result interface{}) error {
	dataB, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(dataB, result)
	if err != nil {
		return err
	}

	return nil
}

func ArrayIndexOf(iarr, elem interface{}) int {
	if arr, ok := iarr.([]string); ok {
		for idx, v := range arr {
			if v == elem.(string) {
				return idx
			}
		}
	}

	return -1
}

func ArrayIntersection(iarr, jarr interface{}) interface{} {
	intersect := []string{}
	arrOne, okOne := iarr.([]string)
	arrTwo, okTwo := jarr.([]string)
	if okOne && okTwo {
		for _, v1 := range arrOne {
			for _, v2 := range arrTwo {
				if v1 == v2 {
					intersect = append(intersect, v1)
				}
			}
		}
	}

	return intersect
}
