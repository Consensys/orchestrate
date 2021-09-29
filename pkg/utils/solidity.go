package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/consensys/orchestrate/pkg/errors"
)

const nameRegex = "[a-zA-Z_][a-zA-Z0-9_]*"
const subTypeRegex = `\[[0-9]*\]`

var dynamicTypes = []string{"bytes", "string"}
var staticTypesAliases = []string{"uint", "int", "byte"}
var staticTypes = []string{"address", "bool"}
var typeRegex string

var signatureRegexp string
var signaturePattern *regexp.Regexp

func init() {
	for i := 8; i < 257; i += 8 {
		staticTypes = append(staticTypes,
			fmt.Sprintf("uint%d", i),
			fmt.Sprintf("int%d", i))
	}
	for i := 1; i < 33; i++ {
		staticTypes = append(staticTypes, fmt.Sprintf("bytes%d", i))
	}

	typeRegex = strings.Join(append(staticTypes, append(dynamicTypes, staticTypesAliases...)...), "|")
	signatureRegexp = fmt.Sprintf(`^(%[1]v)\(((?:%[2]v)(?:%[3]v)*(?:,(?:%[2]v)(?:%[3]v)*)*)?\)$`, nameRegex, typeRegex, subTypeRegex)
	signaturePattern = regexp.MustCompile(signatureRegexp)
}

// ParseSignature returns name and arguments string of a signature
func ParseSignature(sig string) (name, args string, err error) {
	parts := signaturePattern.FindStringSubmatch(sig)

	if len(parts) != 3 {
		return "", "", errors.InvalidSignatureError(sig).SetComponent(component)
	}
	return parts[1], parts[2], nil
}

// VerifyMessageSignature makes hashes short for a limited column size
func IsValidSignature(s string) bool {
	return signaturePattern.MatchString(s)
}
