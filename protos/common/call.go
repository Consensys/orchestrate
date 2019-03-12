package common

import (
	"fmt"
	"regexp"
)

// Short returns a string representation of the method
func (method *Method) Short() string {
	if method.GetContract() == "" {
		return ""
	}

	if method.GetTag() == "" {
		return fmt.Sprintf("%v@%v", method.GetName(), method.GetContract())
	}

	return fmt.Sprintf("%v@%v[%v]", method.GetName(), method.GetContract(), method.GetTag())
}

// IsDeploy indicate wether this method for contract deployment
func (method *Method) IsDeploy() bool {
	return method.Name == "constructor"
}

var methodRegexp = `(?P<name>[a-zA-Z]+)@(?P<contract>[a-zA-Z0-9]+)(\[(?P<tag>[0-9a-zA-Z-\.]+)\])?`
var methodPattern = regexp.MustCompile(methodRegexp)

// FromShortMethod returns a Method object from a short String
func FromShortMethod(s string) (*Method, error) {
	parts := methodPattern.FindStringSubmatch(s)

	if len(parts) < 3 {
		return nil, fmt.Errorf("%v is invalid short method (expected format %q)", s, methodRegexp)
	}

	name, contract, tag := parts[1], parts[2], parts[4]

	return &Method{
		Name:     name,
		Contract: contract,
		Tag:      tag,
	}, nil
}
