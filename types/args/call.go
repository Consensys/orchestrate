package args

import (
	errors "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/utils"
)

var component = "types.args"

// Short returns a string representation of the method
func (c *Call) Short() string {
	return c.GetMethod().GetName()
}

// IsConstructor indicate whether this method for contract deployment
func (c *Call) IsConstructor() bool {
	return c.GetMethod().IsConstructor()
}

// SignatureToCall returns a Call object from a short String
func SignatureToCall(sig string) (*Call, error) {
	if !utils.IsValidSignature(sig) {
		return nil, errors.InvalidSignatureError(sig).SetComponent(component)
	}

	return &Call{
		Method: &abi.Method{
			Signature: sig,
		},
	}, nil
}
