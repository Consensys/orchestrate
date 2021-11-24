package utils

import "strings"

const (
	AuthSeparator = ":"
	BearerPrefix  = "Bearer "
)

func ParseBearerToken(auth string) (string, bool) {
	if len(auth) < len(BearerPrefix) || !strings.EqualFold(auth[:len(BearerPrefix)], BearerPrefix) {
		return "", false
	}

	return auth[len(BearerPrefix):], true
}
