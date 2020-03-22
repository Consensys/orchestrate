package utils

import (
	"sort"
)

func ContainsString(sl []string, v string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func UniqueString(sl []string) []string {
	var uniq []string

	set := make(map[string]struct{})
	for _, v := range sl {
		if _, exist := set[v]; !exist {
			set[v] = struct{}{}
			uniq = append(uniq, v)
		}
	}

	sort.Strings(uniq)

	return uniq
}
