package utils

import "strings"

func ParseAuth(prefix, auth string) (string, bool) {
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return "", false
	}

	return auth[len(prefix):], true
}
