package args

import (
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/utils"
)

// Short returns a string representation of the method
func (c *Call) Short() string {
	return c.GetMethod().GetName()
}

// IsConstructor indicate whether this method for contract deployment
func (c *Call) IsConstructor() bool {
	return c.GetMethod().IsConstructor()
}

// SignatureToCall returns a Call object from a short String
func SignatureToCall(s string) (*Call, error) {
	if !utils.IsValidSignature(s) {
		return nil, fmt.Errorf("invalid signature format, expecting func(type1,type2) got %v", s)
	}

	return &Call{
		Method: &abi.Method{
			Signature: s,
		},
	}, nil
}
