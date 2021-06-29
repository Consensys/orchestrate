package utils

func ArrayIndexOf(iarr, elem interface{}) int {
	if arr, ok := iarr.([]interface{}); ok {
		for idx, v := range arr {
			if v == elem {
				return idx
			}
		}
	}

	return -1
}

func ArrayIntersection(iarr, jarr interface{}) interface{} {
	intersect := []interface{}{}
	arrOne, ok := iarr.([]string)
	if !ok {
		return intersect
	}
	arrTwo, ok := jarr.([]string)
	if !ok {
		return intersect
	}

	for _, v1 := range arrOne {
		for _, v2 := range arrTwo {
			if v1 == v2 {
				intersect = append(intersect, v1)
			}
		}
	}

	return intersect
}
