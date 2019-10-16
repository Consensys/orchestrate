package parser

import (
	"regexp"

	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/errors"
)

var (
	aliasRegexp  = `^(?P<key>[0-9a-zA-Z\.]+):(?P<value>[0-9a-zA-Z]+)`
	aliasPattern = regexp.MustCompile(aliasRegexp)
)

func FromAlias(alias string) (key, value string, err error) {
	parts := aliasPattern.FindStringSubmatch(alias)
	if len(parts) != 3 {
		return "", "", errors.InvalidFormatError("invalid alias %q (expected format %q)", alias, aliasRegexp)
	}

	return parts[1], parts[2], nil
}
