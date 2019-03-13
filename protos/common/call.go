package common

import (
	"fmt"
	"regexp"

	abi "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/abi"
)

// Short returns a string representation of the method
func (c *Call) Short() string {
	if c.GetMethod().GetName() == "" {
		return ""
	}

	return fmt.Sprintf("%v@%v", c.GetMethod().GetName(), c.GetContract().Short())
}

// IsDeploy indicate wether this method for contract deployment
func (c *Call) IsDeploy() bool {
	return c.GetMethod().IsDeploy()
}

var callRegexp = `^(?P<method>[a-zA-Z]+)@(?P<contract>.+)$`
var callPattern = regexp.MustCompile(callRegexp)

// StringToCall returns a Call object from a short String
func StringToCall(s string) (*Call, error) {
	parts := callPattern.FindStringSubmatch(s)

	if len(parts) != 3 {
		return nil, fmt.Errorf("%v is invalid short method (expected format %q)", s, callRegexp)
	}

	contract, err := abi.StringToContract(parts[2])
	if err != nil {
		return nil, err
	}

	return &Call{
		Contract: contract,
		Method: &abi.Method{
			Name: parts[1],
		},
	}, nil
}
